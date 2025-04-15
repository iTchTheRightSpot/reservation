package staff

import (
	"context"
	"github.com/iTchTheRightSpot/reservation/pkg/models/staff"
	"github.com/iTchTheRightSpot/reservation/pkg/stores"
	"github.com/iTchTheRightSpot/utility/cache"
	"github.com/iTchTheRightSpot/utility/utils"
)

type IStaffService interface {
	AllUsers(ctx context.Context) ([]*staff.AllStaffsEntity, error)
}

type staffService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
	cache    cache.ICache[string, []*staff.AllStaffsEntity]
	key      string
}

func NewStaffService(l utils.ILogger, a *stores.Adapters, c cache.ICache[string, []*staff.AllStaffsEntity]) IStaffService {
	return &staffService{logger: l, adapters: a, cache: c, key: "key"}
}

func (dep *staffService) AllUsers(ctx context.Context) ([]*staff.AllStaffsEntity, error) {
	val := dep.cache.Get(dep.key)
	if val != nil {
		return *val, nil
	}

	arr, err := dep.adapters.StaffStore.AllStaffs(ctx)
	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return nil, &utils.NotFoundError{Message: "error retrieving staffs"}
	}

	dep.cache.Put(dep.key, arr)
	return arr, nil
}