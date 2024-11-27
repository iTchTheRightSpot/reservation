package reservation

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
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
}

func NewReservationService(l utils.ILogger, a *stores.Adapters) IReservationService {
	return &reservationService{logger: l, adapters: a}
}

func (dep *reservationService) matchStaffServices(requestedServices []*string, availableServices []*service.ServiceEntity) ([]*service.ServiceEntity, error) {
	arr := make([]*service.ServiceEntity, 0)

	for _, entity := range availableServices {
		match := slices.ContainsFunc(requestedServices, func(s *string) bool {
			lc := strings.ToUpper(strings.TrimSpace(*s))
			up := strings.ToUpper(strings.TrimSpace(entity.Name))
			return lc == up
		})

		if match {
			arr = append(arr, entity)
		}
	}

	if len(arr) != len(requestedServices) {
		return nil, fmt.Errorf("1 or more services were not found for selected staff")
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

func (dep *reservationService) Create(ctx context.Context, p *reservation.ReservationPayload) error {
	s, err := dep.adapters.StaffStore.StaffByUUID(ctx, strings.TrimSpace(p.StaffId))
	if err != nil {
		return err
	}

	ser, err := dep.adapters.ServiceStore.ServicesByStaffId(ctx, s.StaffId)
	if err != nil {
		return err
	}

	matchedServices, err := dep.matchStaffServices(p.Services, ser)
	if err != nil {
		dep.logger.Error(err.Error())
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

	end := start.Add(time.Second * time.Duration(dep.sumUpServiceDuration(matchedServices)))

	count, err := dep.adapters.ScheduleStore.CountSchedulesInRangeAndVisibility(ctx, s.StaffId, start, end, true)
	if count <= 0 {
		return fmt.Errorf("invalid reservation time")
	}

	return fmt.Errorf("yet to implement reservation service")
}
