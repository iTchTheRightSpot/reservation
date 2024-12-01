package reservation

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/mail"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"slices"
	"strconv"
	"strings"
	"time"
)

type IReservationService interface {
	Create(ctx context.Context, p *reservation.ReservationPayload) error
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

func (dep *reservationService) AvailableDates(o *reservation.AvailableTimesPayload) ([]reservation.ReservationTimeSlots, error) {

	return nil, fmt.Errorf("yet to implement")
}

func (dep *reservationService) services(ctx context.Context, p *reservation.ReservationPayload, err error, staffObj *staff.Staff) ([]*service.ServiceEntity, error) {
	serviceEntities, err := dep.adapters.ServiceStore.ServicesByStaffId(ctx, staffObj.StaffId)
	if err != nil {
		return nil, err
	}

	arr := make([]*service.ServiceEntity, 0)

	for _, entity := range serviceEntities {
		match := slices.ContainsFunc(p.Services, func(s *string) bool {
			lc := strings.ToUpper(strings.TrimSpace(*s))
			up := strings.ToUpper(strings.TrimSpace(entity.Name))
			return lc == up
		})

		if match {
			arr = append(arr, entity)
		}
	}

	if len(arr) != len(p.Services) {
		mess := "1 or more services were not found for selected staff"
		dep.logger.Error(mess)
		return nil, fmt.Errorf(mess)
	}

	return arr, nil
}

func (dep *reservationService) dateInTimezone(p *reservation.ReservationPayload, t *time.Location) (*time.Time, error) {
	num, err := strconv.ParseInt(p.Time, 10, 64)
	if err != nil {
		return nil, err
	}

	if len(p.Timezone) < 1 {
		in := time.UnixMilli(num).In(t)
		return &in, nil
	}

	l, err := time.LoadLocation(p.Timezone)
	if err != nil {
		return nil, err
	}

	in := time.UnixMilli(num).In(l).In(t)
	return &in, nil
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

func (dep *reservationService) createReservation(ctx context.Context, p *reservation.ReservationPayload, matchedServices []*service.ServiceEntity, s *staff.Staff, start *time.Time, end time.Time) error {
	return dep.adapters.Transaction.RunInTransaction(func(adapters *stores.Adapters) error {
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
			ScheduledFor: *start,
			ExpireAt:     end,
		}

		if err := adapters.ReservationStore.SelectForUpdateSave(ctx, reserv); err != nil {
			dep.logger.Error(err.Error())
			return fmt.Errorf("error creating reservation")
		}

		for _, entity := range matchedServices {
			err := adapters.ReservationServiceStore.Save(ctx, &reservation.ReservationServiceEntity{
				ReservationId: reserv.ReservationId,
				ServiceId:     entity.ServiceId,
			})
			if err != nil {
				dep.logger.Error(err)
				return fmt.Errorf("error creating reservation")
			}
		}
		return nil
	})
}

func (dep *reservationService) Create(ctx context.Context, p *reservation.ReservationPayload) error {
	staffObj, err := dep.adapters.StaffStore.StaffByUUID(ctx, strings.TrimSpace(p.StaffId))
	if err != nil {
		return err
	}

	services, err := dep.services(ctx, p, err, staffObj)
	if err != nil {
		return err
	}

	start, err := dep.dateInTimezone(p, dep.logger.Timezone())
	if err != nil {
		dep.logger.Error(err.Error())
		return err
	}

	if start.Before(dep.logger.Date()) {
		dep.logger.Error("cannot make a reservation for a past day")
		return fmt.Errorf("cannot make a reservation for a past day")
	}

	end := start.Add(time.Second * time.Duration(dep.sumUpServiceDuration(services)))

	count, err := dep.adapters.ScheduleStore.CountSchedulesInRangeAndVisibility(ctx, staffObj.StaffId, *start, end, true)
	if err != nil {
		return err
	}
	if count < 1 {
		mess := "invalid reservation time"
		dep.logger.Error(mess)
		return fmt.Errorf(mess)
	}

	count, err = dep.adapters.ReservationStore.CountReservationsInRangeByStaffTimeAndStatuses(
		ctx, staffObj.StaffId, *start, end, reservation.CONFIRMED)

	if err != nil {
		return err
	}
	if count > 0 {
		mess := "reservation creation failed. conflict"
		dep.logger.Error(mess)
		return fmt.Errorf(mess)
	}

	err = dep.createReservation(ctx, p, services, staffObj, start, end)
	if err != nil {
		return err
	}
	dep.cache.Clear()
	return dep.mail.SendReservationConfirmation()
}
