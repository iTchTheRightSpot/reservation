package service_type

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"strings"
)

type IService interface {
	Create(ctx context.Context, p *service.ServiceTypePayload) error
	Services(ctx context.Context) (interface{}, error)
}

type serviceImpl struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewServiceImpl(l utils.ILogger, a *stores.Adapters) IService {
	return &serviceImpl{logger: l, adapters: a}
}

func (dep *serviceImpl) Services(ctx context.Context) (interface{}, error) {
	db, err := dep.adapters.ServiceStore.ServicesByStatus(ctx, true)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "an error occurred retrieving service types"}
	}

	type ui struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Duration int     `json:"duration"`
	}

	arr := make([]*ui, len(db))

	for i, e := range db {
		arr[i] = &ui{Name: e.Name, Price: e.Price, Duration: e.Duration}
	}

	return arr, nil
}

func (dep *serviceImpl) Create(ctx context.Context, p *service.ServiceTypePayload) error {
	s := service.ServiceTypeEntity{
		Name:          strings.TrimSpace(p.Name),
		Price:         p.Price,
		IsVisible:     p.IsVisible,
		IsReoccurring: p.IsReoccurring,
		Duration:      p.Duration,
		CleanUpTime:   p.CleanUpTime,
	}

	if err := dep.adapters.ServiceStore.Save(ctx, &s); err != nil {
		dep.logger.Error(err.Error())
		return &utils.InsertionError{Message: err.Error()}
	}
	return nil
}
