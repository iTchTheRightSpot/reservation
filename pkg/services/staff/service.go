package staff

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IStaffService interface {
	LinkServiceToStaff(ctx context.Context, staffUUID, serviceName string) error
}

type staffService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewStaffService(l utils.ILogger, a *stores.Adapters) IStaffService {
	return &staffService{logger: l, adapters: a}
}

func (dep *staffService) LinkServiceToStaff(ctx context.Context, staffUUID, serviceName string) error {
	s, err := dep.adapters.StaffStore.StaffByUUID(ctx, staffUUID)
	if err != nil {
		return &utils.NotFoundError{Message: "invalid staff id"}
	}

	service, err := dep.adapters.ServiceStore.ServiceByName(ctx, serviceName)
	if err != nil {
		return &utils.NotFoundError{Message: "invalid service name"}
	}

	count, err := dep.adapters.StaffServiceStore.CountByStaffIdAndServiceId(ctx, s.StaffId, service.ServiceId)
	if err != nil {
		return &utils.InsertionError{Message: err.Error()}
	}

	if count > 0 {
		return &utils.InsertionError{Message: "service already linked to staff"}
	}

	err = dep.adapters.StaffServiceStore.Save(ctx, &staff.StaffServiceEntity{
		StaffId: s.StaffId, ServiceId: service.ServiceId,
	})

	if err != nil {
		return &utils.InsertionError{Message: err.Error()}
	}

	return nil
}
