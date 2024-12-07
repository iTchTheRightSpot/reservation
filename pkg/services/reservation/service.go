package reservation

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/mail"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IReservationService interface {
	Create(ctx context.Context, p *reservation.ReservationPayload) error
	AvailableDates(ctx context.Context, o *reservation.AvailableTimesPayload) (*[]reservation.ReservationTimeSlots, error)
}

type reservationService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
	cache    pkg.ICache[string, []reservation.ReservationTimeSlots]
	mail     mail.IMailService
}

func NewReservationService(
	l utils.ILogger,
	a *stores.Adapters,
	c pkg.ICache[string, []reservation.ReservationTimeSlots],
	m mail.IMailService,
) IReservationService {
	return &reservationService{logger: l, adapters: a, cache: c, mail: m}
}

func (dep *reservationService) createKey(o *reservation.AvailableTimesPayload) string {
	return fmt.Sprintf("%s_%v_%v_%v_%v_%s", strings.Join(o.Services, "_"), o.StaffId, o.Day, o.Month, o.Year, o.StartDateTime.Location().String())
}

func (dep *reservationService) generateChunks(schedules []*schedule.Schedule, duration int) []*reservation.Chunks {
	var arr []*reservation.Chunks
	var wg sync.WaitGroup
	var mu sync.Mutex // to ensure safe concurrent access to 'arr'

	for _, sch := range schedules {
		wg.Add(1)
		go func(sch *schedule.Schedule) {
			defer wg.Done()

			var times []time.Time
			tempStart := sch.Start

			for tempStart.Before(sch.End) {
				times = append(times, tempStart)
				tempStart = tempStart.Add(time.Duration(duration) * time.Second)
			}

			mu.Lock()
			arr = append(arr, &reservation.Chunks{
				Start: sch.Start,
				Times: times,
			})
			mu.Unlock()
		}(sch)
	}

	wg.Wait()
	return arr
}

func (dep *reservationService) filterChunks(ctx context.Context, staffId uint64, duration int, chunks []*reservation.Chunks, location *time.Location) (*[]reservation.ReservationTimeSlots, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []reservation.ReservationTimeSlots
	errChan := make(chan error, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)

		go func(chunk *reservation.Chunks) {
			defer wg.Done()

			var validTimes []string

			for _, timeSlot := range chunk.Times {
				end := timeSlot.Add(time.Duration(duration) * time.Second)

				count, err := dep.adapters.ReservationStore.CountReservationsInRange(ctx, staffId, timeSlot, end, reservation.CONFIRMED)
				if err != nil {
					errChan <- &utils.BadRequestError{Message: err.Error()}
					return
				}

				if count == 0 {
					validTimes = append(validTimes, fmt.Sprintf("%v", timeSlot.In(location).UnixMilli()))
				}
			}

			if len(validTimes) > 0 {
				mu.Lock()
				results = append(results, reservation.ReservationTimeSlots{
					Date:  fmt.Sprintf("%v", chunk.Start.In(location).UnixMilli()),
					Times: validTimes,
				})
				mu.Unlock()
			}
		}(chunk)
	}

	wg.Wait()
	close(errChan)

	return &results, <-errChan
}

func (dep *reservationService) AvailableDates(ctx context.Context, o *reservation.AvailableTimesPayload) (*[]reservation.ReservationTimeSlots, error) {
	key := dep.createKey(o)
	val := dep.cache.Get(key)
	if val != nil {
		return val, nil
	}

	staf, err := dep.adapters.StaffStore.StaffByUUID(ctx, o.StaffId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "invalid staff Id"}
	}

	services, err := dep.matchStaffServices(ctx, o.Services, staf)
	if err != nil {
		return nil, err
	}

	duration := dep.sumUpServiceDuration(services)

	schedules, err := dep.adapters.ScheduleStore.SchedulesInRangeAndVisibilityAndDifference(
		ctx, staf.StaffId, o.StartDateTime, o.EndDateTime, true, duration)

	chunks := dep.generateChunks(schedules, duration)

	filter, err := dep.filterChunks(ctx, staf.StaffId, duration, chunks, o.StartDateTime.Location())
	if err != nil {
		return nil, err
	}

	dep.cache.Put(key, *filter)

	return filter, nil
}

func (dep *reservationService) matchStaffServices(ctx context.Context, requestedServices []string, staffObj *staff.Staff) ([]*service.ServiceEntity, error) {
	serviceEntities, err := dep.adapters.ServiceStore.ServicesByStaffId(ctx, staffObj.StaffId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "invalid service for staff"}
	}

	arr := make([]*service.ServiceEntity, 0)

	for _, entity := range serviceEntities {
		match := slices.ContainsFunc(requestedServices, func(s string) bool {
			lc := strings.ToUpper(strings.TrimSpace(s))
			up := strings.ToUpper(strings.TrimSpace(entity.Name))
			return lc == up
		})

		if match {
			arr = append(arr, entity)
		}
	}

	if len(arr) != len(requestedServices) {
		mess := "1 or more services were not found for selected staff"
		dep.logger.Error(mess)
		return nil, &utils.BadRequestError{Message: mess}
	}

	return arr, nil
}

func (dep *reservationService) sumUpServiceDuration(s []*service.ServiceEntity) int {
	count := 0
	for _, entity := range s {
		count += entity.Duration + entity.CleanUpTime
	}
	return count
}

func (dep *reservationService) sumUpServicePrice(s []*service.ServiceEntity) float64 {
	var count float64
	for _, entity := range s {
		count += entity.Price
	}
	return count
}

func (dep *reservationService) createReservation(ctx context.Context, p *reservation.ReservationPayload, matchedServices []*service.ServiceEntity, s *staff.Staff, start time.Time, end time.Time) error {
	return dep.adapters.Transaction.RunInTransaction(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable}, func(adapters *stores.Adapters) error {
		priceSum := dep.sumUpServicePrice(matchedServices)

		reserv := &reservation.Reservation{
			StaffId:      s.StaffId,
			Name:         strings.TrimSpace(p.Name),
			Email:        strings.TrimSpace(p.Email),
			Description:  &p.Description,
			Address:      &p.Address,
			Phone:        &p.Phone,
			ImageKey:     nil,
			Price:        priceSum,
			Status:       reservation.CONFIRMED,
			CreatedAt:    dep.logger.Date(),
			ScheduledFor: start,
			ExpireAt:     end,
		}

		if err := adapters.ReservationStore.Save(ctx, reserv); err != nil {
			dep.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error creating reservation"}
		}

		for _, entity := range matchedServices {
			err := adapters.ReservationServiceStore.Save(ctx, &reservation.ReservationServiceEntity{
				ReservationId: reserv.ReservationId,
				ServiceId:     entity.ServiceId,
			})
			if err != nil {
				dep.logger.Error(err)
				return &utils.InsertionError{Message: "error creating reservation"}
			}
		}
		return nil
	})
}

func (dep *reservationService) Create(ctx context.Context, p *reservation.ReservationPayload) error {
	staffObj, err := dep.adapters.StaffStore.StaffByUUID(ctx, strings.TrimSpace(p.StaffId))
	if err != nil {
		return &utils.NotFoundError{Message: "invalid staff id"}
	}

	services, err := dep.matchStaffServices(ctx, p.Services, staffObj)
	if err != nil {
		return err
	}

	parseInt, err := strconv.ParseInt(p.Time, 10, 64)
	if err != nil {
		return &utils.BadRequestError{Message: err.Error()}
	}

	start := time.UnixMilli(parseInt)
	if start.Before(dep.logger.Date()) {
		mess := "cannot make a reservation for a past day"
		dep.logger.Error(mess)
		return &utils.BadRequestError{Message: mess}
	}

	end := start.Add(time.Second * time.Duration(dep.sumUpServiceDuration(services)))

	count, err := dep.adapters.ScheduleStore.CountSchedulesInRangeAndVisibility(ctx, staffObj.StaffId, start, end, true)
	if err != nil {
		return err
	}
	if count < 1 {
		mess := "invalid reservation time"
		dep.logger.Error(mess)
		return &utils.BadRequestError{Message: mess}
	}

	count, err = dep.adapters.ReservationStore.CountReservationsInRange(
		ctx, staffObj.StaffId, start, end, reservation.CONFIRMED)

	if err != nil {
		return err
	}
	if count > 0 {
		mess := "reservation time is not available failed"
		dep.logger.Error(mess)
		return &utils.BadRequestError{Message: mess}
	}

	if err = dep.createReservation(ctx, p, services, staffObj, start, end); err != nil {
		return err
	}

	dep.cache.Clear()
	return dep.mail.SendReservationConfirmation()
}
