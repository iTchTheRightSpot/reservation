package staff

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IStaffService interface {
	AllUsers(ctx context.Context) ([]*staff.AllStaffsEntity, error)
}

type staffService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
	cache    pkg.ICache[string, []*staff.AllStaffsEntity]
	key      string
}

func NewStaffService(l utils.ILogger, a *stores.Adapters, c pkg.ICache[string, []*staff.AllStaffsEntity]) IStaffService {
	return &staffService{logger: l, adapters: a, cache: c, key: "key"}
}

func (dep *staffService) AllUsers(ctx context.Context) ([]*staff.AllStaffsEntity, error) {
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
