package service_type

import (
	"context"
	"github.com/iTchTheRightSpot/reservation/pkg/models/service_type"
	"github.com/iTchTheRightSpot/reservation/pkg/models/staff"
	"github.com/iTchTheRightSpot/reservation/pkg/stores"
	"github.com/iTchTheRightSpot/utility/utils"
	"strings"
)

type IServiceType interface {
	Create(ctx context.Context, p *service_type.ServiceTypePayload) error
	ServiceTypes(ctx context.Context) (interface{}, error)
	StaffsByServiceTypes(ctx context.Context, services *[]string) (interface{}, error)
	CRMServiceTypes(ctx context.Context) (interface{}, error)
	LinkServiceToStaff(ctx context.Context, obj *service_type.LinkServiceTypeToStaffPayload) error
	Update(ctx context.Context, dto *service_type.ServiceTypePayload) error
	ServicesByStaffUUID(ctx context.Context, staffUUID string) (interface{}, error)
}

type serviceTypeImpl struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewServiceImpl(l utils.ILogger, a *stores.Adapters) IServiceType {
	return &serviceTypeImpl{logger: l, adapters: a}
}

func (dep *serviceTypeImpl) Update(ctx context.Context, dto *service_type.ServiceTypePayload) error {
	s, err := dep.adapters.ServiceStore.ServiceTypeByName(ctx, strings.TrimSpace(dto.Name))
	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return err
	}
	s.Name = strings.TrimSpace(dto.Name)
	s.Price = dto.Price
	s.IsVisible = dto.IsVisible
	s.Duration = dto.Duration
	s.CleanUpTime = dto.CleanUpTime
	return dep.adapters.ServiceStore.Update(ctx, s)
}

func (dep *serviceTypeImpl) CRMServiceTypes(ctx context.Context) (interface{}, error) {
	db, err := dep.adapters.ServiceStore.ServiceTypes(ctx)
	if err != nil {
		return nil, err
	}

	type ui struct {
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		IsVisible bool    `json:"is_visible"`
		Duration  int     `json:"duration"`
		CleanUp   int     `json:"clean_up_time"`
	}

	arr := make([]*ui, len(db))

	for i, e := range db {
		arr[i] = &ui{Name: e.Name, Price: e.Price, IsVisible: e.IsVisible, Duration: e.Duration, CleanUp: e.CleanUpTime}
	}

	if len(arr) == 0 {
		return []interface{}{}, nil
	}

	return arr, nil
}

func (dep *serviceTypeImpl) ServiceTypes(ctx context.Context) (interface{}, error) {
	db, err := dep.adapters.ServiceStore.ServiceTypes(ctx)
	if err != nil {
		return nil, err
	}

	type ui struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Duration int     `json:"duration"`
	}

	var arr []*ui

	for _, e := range db {
		if e.IsVisible {
			arr = append(arr, &ui{Name: e.Name, Price: e.Price, Duration: e.Duration})
		}
	}

	if len(arr) == 0 {
		return []interface{}{}, nil
	}

	return arr, nil
}

func (dep *serviceTypeImpl) StaffsByServiceTypes(ctx context.Context, services *[]string) (interface{}, error) {
	stafs, err := dep.adapters.StaffStore.StaffsByServices(ctx, services)
	if err != nil {
		return nil, err
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

func (dep *serviceTypeImpl) Create(ctx context.Context, p *service_type.ServiceTypePayload) error {
	s := service_type.ServiceTypeEntity{
		Name:        strings.TrimSpace(p.Name),
		Price:       p.Price,
		IsVisible:   p.IsVisible,
		Duration:    p.Duration,
		CleanUpTime: p.CleanUpTime,
	}

	if err := dep.adapters.ServiceStore.Save(ctx, &s); err != nil {
		dep.logger.Error(ctx, err.Error())
		return err
	}
	return nil
}

func (dep *serviceTypeImpl) LinkServiceToStaff(ctx context.Context, obj *service_type.LinkServiceTypeToStaffPayload) error {
	s, err := dep.adapters.StaffStore.StaffByUUID(ctx, obj.StaffUUID)
	if err != nil {
		return err
	}

	service, err := dep.adapters.ServiceStore.ServiceTypeByName(ctx, obj.Service)
	if err != nil {
		return err
	}

	count, err := dep.adapters.StaffServiceStore.CountByStaffIdAndServiceId(ctx, s.StaffId, service.ServiceId)
	if err != nil {
		return err
	}

	if count > 0 {
		return &utils.InsertionError{Message: "service already linked to staff"}
	}

	err = dep.adapters.StaffServiceStore.Save(ctx, &staff.StaffServiceEntity{
		StaffId: s.StaffId, ServiceId: service.ServiceId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (dep *serviceTypeImpl) ServicesByStaffUUID(ctx context.Context, staffUUID string) (interface{}, error) {
	arr, err := dep.adapters.ServiceStore.ServiceTypesByStaffUUID(ctx, strings.TrimSpace(staffUUID))
	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return []interface{}{}, nil
	}

	return arr, nil
}