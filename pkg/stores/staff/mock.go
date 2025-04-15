package staff

import (
	"context"
	"github.com/iTchTheRightSpot/reservation/pkg/models/staff"
)

type MockStaffStore struct {
	StaffSave         *staff.StaffEntity
	StaffSaveError    error
	StaffSaveCalled   bool
	StaffByUUIDReturn *staff.StaffEntity
	StaffByUUIDError  error
	StaffByUUIDCalled bool
}

func (dep *MockStaffStore) StaffsByServices(context.Context, *[]string) ([]*staff.StaffStoreFrontDb, error) {
	return nil, nil
}

func (dep *MockStaffStore) AllStaffs(context.Context) ([]*staff.AllStaffsEntity, error) {
	return nil, nil
}

func (dep *MockStaffStore) Save(context.Context, *staff.StaffEntity) error {
	dep.StaffSaveCalled = true
	return dep.StaffSaveError
}

func (dep *MockStaffStore) StaffByProfileId(context.Context, uint64) (*staff.StaffEntity, error) {
	panic("implement me")
}

func (dep *MockStaffStore) StaffByUUID(context.Context, string) (*staff.StaffEntity, error) {
	dep.StaffByUUIDCalled = true
	return dep.StaffByUUIDReturn, dep.StaffByUUIDError
}

type MockStaffServiceStore struct {
	StaffServiceEntitySave           *staff.StaffServiceEntity
	StaffServiceEntitySaveError      error
	StaffServiceEntitySaveCalled     bool
	CountByStaffIdAndServiceIdResult int
	CountByStaffIdAndServiceIdError  error
	CountByStaffIdAndServiceIdCalled bool
}

func (dep *MockStaffServiceStore) Save(context.Context, *staff.StaffServiceEntity) error {
	dep.StaffServiceEntitySaveCalled = true
	return dep.StaffServiceEntitySaveError
}

func (dep *MockStaffServiceStore) CountByStaffIdAndServiceId(context.Context, uint64, uint64) (int, error) {
	dep.CountByStaffIdAndServiceIdCalled = true
	return dep.CountByStaffIdAndServiceIdResult, dep.CountByStaffIdAndServiceIdError
}