package staff

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
)

type MockStaffStore struct {
	StaffSave         *staff.Staff
	StaffSaveError    error
	StaffSaveCalled   bool
	StaffByUUIDReturn *staff.Staff
	StaffByUUIDError  error
	StaffByUUIDCalled bool
}

func (dep *MockStaffStore) Save(context.Context, *staff.Staff) error {
	dep.StaffSaveCalled = true
	return dep.StaffSaveError
}

func (dep *MockStaffStore) StaffByUUID(context.Context, string) (*staff.Staff, error) {
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
