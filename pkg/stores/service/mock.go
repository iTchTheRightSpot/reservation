package service

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
)

type MockServiceStore struct {
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

func (dep *MockServiceStore) Save(context.Context, *service.ServiceEntity) (*service.ServiceEntity, error) {
	dep.ServiceSaveCalled = true
	return dep.ServiceSave, dep.ServiceSaveError
}

func (dep *MockServiceStore) ServiceByName(context.Context, string) (*service.ServiceEntity, error) {
	dep.ServiceByNameCalled = true
	return dep.ServiceByNameReturn, dep.ServiceByNameError
}

func (dep *MockServiceStore) ServicesByStaffId(context.Context, uint64) ([]*service.ServiceEntity, error) {
	dep.ServicesByStaffIdCalled = true
	return dep.ServicesByStaffIdReturn, dep.ServicesByStaffIdError
}
