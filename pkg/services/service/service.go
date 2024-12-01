package service

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"strings"
)

type IService interface {
	Create(ctx context.Context, p *service.ServicePayload) error
}

type serviceImpl struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewServiceImpl(l utils.ILogger, a *stores.Adapters) IService {
	return &serviceImpl{logger: l, adapters: a}
}

func (dep *serviceImpl) Create(ctx context.Context, p *service.ServicePayload) error {
	s := service.ServiceEntity{
		Name:          strings.TrimSpace(p.Name),
		Price:         p.Price,
		IsVisible:     p.IsVisible,
		IsReoccurring: p.IsReoccurring,
		Duration:      p.Duration,
		CleanUpTime:   p.CleanUpTime,
	}

	if _, err := dep.adapters.ServiceStore.Save(ctx, &s); err != nil {
		dep.logger.Error(err.Error())
		return &utils.InsertionError{Message: err.Error()}
	}
	return nil
}
