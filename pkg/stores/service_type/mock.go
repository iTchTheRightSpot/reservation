package service_type

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
)

type MockServiceTypeStore struct {
	ServiceSave             *service.ServiceEntity
	ServiceSaveError        error
	ServiceSaveCalled       bool
	ServiceByNameReturn     *service.ServiceEntity
	ServiceByNameError      error
	ServiceByNameCalled     bool
	ServicesByStaffIdReturn []*service.ServiceEntity
	ServicesByStaffIdError  error
	ServicesByStaffIdCalled bool
}

func (dep *MockServiceTypeStore) Save(context.Context, *service.ServiceEntity) (*service.ServiceEntity, error) {
	dep.ServiceSaveCalled = true
	return dep.ServiceSave, dep.ServiceSaveError
}

func (dep *MockServiceTypeStore) ServiceByName(context.Context, string) (*service.ServiceEntity, error) {
	dep.ServiceByNameCalled = true
	return dep.ServiceByNameReturn, dep.ServiceByNameError
}

func (dep *MockServiceTypeStore) ServicesByStaffId(context.Context, uint64) ([]*service.ServiceEntity, error) {
	dep.ServicesByStaffIdCalled = true
	return dep.ServicesByStaffIdReturn, dep.ServicesByStaffIdError
}
