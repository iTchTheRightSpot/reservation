package reservation

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"slices"
	"strings"
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

func (dep *reservationService) Create(ctx context.Context, p *reservation.ReservationPayload) error {
	s, err := dep.adapters.StaffStore.StaffByUUID(ctx, strings.TrimSpace(p.StaffId))
	if err != nil {
		return err
	}

	ser, err := dep.adapters.ServiceStore.ServicesByStaffId(ctx, s.StaffId)
	if err != nil {
		return err
	}

	//matchedServices, err := dep.matchStaffServices(p.Services, ser)
	_, err = dep.matchStaffServices(p.Services, ser)
	if err != nil {
		return err
	}

	return fmt.Errorf("yet to implement reservation service")
}
