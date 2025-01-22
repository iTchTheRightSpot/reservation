package service_type

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"strings"
)

type IServiceType interface {
	Create(ctx context.Context, p *service.ServiceTypePayload) error
	ServiceTypes(ctx context.Context) (interface{}, error)
	StaffsByServiceTypes(ctx context.Context, services *[]string) (interface{}, error)
}

type serviceTypeImpl struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewServiceImpl(l utils.ILogger, a *stores.Adapters) IServiceType {
	return &serviceTypeImpl{logger: l, adapters: a}
}

func (dep *serviceTypeImpl) ServiceTypes(ctx context.Context) (interface{}, error) {
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

func (dep *serviceTypeImpl) StaffsByServiceTypes(ctx context.Context, services *[]string) (interface{}, error) {
	stafs, err := dep.adapters.StaffStore.StaffsByServices(ctx, services)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "no staffs found for said service(s)"}
	}

	type ui struct {
		StaffId  string  `json:"staff_id"`
		Name     string  `json:"name"`
		ImageKey *string `json:"image_key"`
		Bio      *string `json:"bio"`
	}

	arr := make([]*ui, len(stafs))

	var str strings.Builder
	for i, staf := range stafs {
		str.Reset()

		if staf.Bio == nil {
			str.WriteString("Ready to put a smile on your face🌞")
		} else {
			str.WriteString(*staf.Bio)
		}
		ss := str.String()
		arr[i] = &ui{StaffId: staf.UUID.String(), Name: staf.Name, ImageKey: staf.ImageKey, Bio: &ss}
	}

	return arr, nil
}

func (dep *serviceTypeImpl) Create(ctx context.Context, p *service.ServiceTypePayload) error {
	s := service.ServiceTypeEntity{
		Name:        strings.TrimSpace(p.Name),
		Price:       p.Price,
		IsVisible:   p.IsVisible,
		Duration:    p.Duration,
		CleanUpTime: p.CleanUpTime,
	}

	if err := dep.adapters.ServiceStore.Save(ctx, &s); err != nil {
		dep.logger.Error(err.Error())
		return &utils.InsertionError{Message: err.Error()}
	}
	return nil
}
