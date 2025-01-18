package service_type

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
)

type MockServiceTypeStore struct {
	ServiceSave             *service.ServiceTypeEntity
	ServiceSaveError        error
	ServiceSaveCalled       bool
	ServiceByNameReturn     *service.ServiceTypeEntity
	ServiceByNameError      error
	ServiceByNameCalled     bool
	ServicesByStaffIdReturn []*service.ServiceTypeEntity
	ServicesByStaffIdError  error
	ServicesByStaffIdCalled bool
}

func (dep *MockServiceTypeStore) ServicesByStatus(context.Context, bool) ([]*service.ServiceTypeEntity, error) {
	return nil, nil
}

func (dep *MockServiceTypeStore) Save(context.Context, *service.ServiceTypeEntity) error {
	dep.ServiceSaveCalled = true
	return dep.ServiceSaveError
}

func (dep *MockServiceTypeStore) ServiceByName(context.Context, string) (*service.ServiceTypeEntity, error) {
	dep.ServiceByNameCalled = true
	return dep.ServiceByNameReturn, dep.ServiceByNameError
}

func (dep *MockServiceTypeStore) ServicesByStaffId(context.Context, uint64, bool) ([]*service.ServiceTypeEntity, error) {
	dep.ServicesByStaffIdCalled = true
	return dep.ServicesByStaffIdReturn, dep.ServicesByStaffIdError
}
