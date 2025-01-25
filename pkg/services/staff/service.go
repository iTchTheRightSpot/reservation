package staff

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IStaffService interface {
	LinkServiceToStaff(ctx context.Context, staffUUID, serviceName string) error
	AllStaffs(ctx context.Context) ([]*staff.AllStaffsEntity, error)
}

type staffService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
	cache pkg.ICache[string, []*staff.AllStaffsEntity]
	key string
}

func NewStaffService(l utils.ILogger, a *stores.Adapters, c pkg.ICache[string, []*staff.AllStaffsEntity]) IStaffService {
	return &staffService{logger: l, adapters: a, cache: c, key: "key"}
}

func (dep *staffService) AllStaffs(ctx context.Context) ([]*staff.AllStaffsEntity, error) {
	val := dep.cache.Get(dep.key)
	if val != nil {
		return *val, nil
	}

	arr, err := dep.adapters.StaffStore.AllStaffs(ctx)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "error retrieving staffs"}
	}

	dep.cache.Put(dep.key, arr)
	return arr, nil
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
