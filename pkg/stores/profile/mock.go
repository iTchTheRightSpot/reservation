package profile

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
)

type MockProfileStore struct {
	MockProfileSave       *models.ProfileEntity
	MockProfileSaveError  error
	MockProfileSaveCalled bool
}

func (dep *MockProfileStore) ProfileByEmail(context.Context, string) (*models.ProfileEntity, error) {
	return nil, nil
}

func (dep *MockProfileStore) ProfileRolesAndPermissionByEmail(context.Context, string) (*models.ProfileRolePermissionEntity, error) {
	return nil, nil
}

func (dep *MockProfileStore) ProfileByStaffUUID(context.Context, string) (*models.ProfileEntity, error) {
	return nil, nil
}

func (dep *MockProfileStore) Save(context.Context, *models.ProfileEntity) error {
	dep.MockProfileSaveCalled = true
	return dep.MockProfileSaveError
}

type MockPermissionStore struct {
	MockPermissionSave       *models.PermissionEntity
	MockPermissionSaveError  error
	MockPermissionSaveCalled bool
}

func (dep *MockPermissionStore) Save(context.Context, *models.PermissionEntity) error {
	dep.MockPermissionSaveCalled = true
	return dep.MockPermissionSaveError
}

type MockRoleStore struct {
	MockRoleSave       *models.RoleEntity
	MockRoleSaveError  error
	MockRoleSaveCalled bool
}

func (dep *MockRoleStore) Save(context.Context, *models.RoleEntity) error {
	dep.MockRoleSaveCalled = true
	return dep.MockRoleSaveError
}
