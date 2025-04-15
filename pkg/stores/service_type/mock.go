package service_type

import (
	"context"
	"github.com/iTchTheRightSpot/reservation/pkg/models/service_type"
)

type MockServiceTypeStore struct {
	ServiceSave             *service_type.ServiceTypeEntity
	ServiceSaveError        error
	ServiceSaveCalled       bool
	ServiceByNameReturn     *service_type.ServiceTypeEntity
	ServiceByNameError      error
	ServiceByNameCalled     bool
	ServicesByStaffIdReturn []*service_type.ServiceTypeEntity
	ServicesByStaffIdError  error
	ServicesByStaffIdCalled bool
}

func (dep *MockServiceTypeStore) Update(context.Context, *service_type.ServiceTypeEntity) error {
	return nil
}

func (dep *MockServiceTypeStore) ServiceTypes(context.Context) ([]*service_type.ServiceTypeEntity, error) {
	return nil, nil
}

func (dep *MockServiceTypeStore) Save(context.Context, *service_type.ServiceTypeEntity) error {
	dep.ServiceSaveCalled = true
	return dep.ServiceSaveError
}

func (dep *MockServiceTypeStore) ServiceTypeByName(context.Context, string) (*service_type.ServiceTypeEntity, error) {
	dep.ServiceByNameCalled = true
	return dep.ServiceByNameReturn, dep.ServiceByNameError
}

func (dep *MockServiceTypeStore) ServiceTypesByStaffId(context.Context, uint64, bool) ([]*service_type.ServiceTypeEntity, error) {
	dep.ServicesByStaffIdCalled = true
	return dep.ServicesByStaffIdReturn, dep.ServicesByStaffIdError
}

func (dep *MockServiceTypeStore) ServiceTypesByStaffUUID(context.Context, string) ([]string, error) {
	dep.ServicesByStaffIdCalled = true
	return nil, dep.ServicesByStaffIdError
}