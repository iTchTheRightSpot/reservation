package profile

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
)

type MockProfileStore struct {
	MockProfileSave       *profile.Profile
	MockProfileSaveError  error
	MockProfileSaveCalled bool
}

func (dep *MockProfileStore) Save(context.Context, *profile.Profile) error {
	dep.MockProfileSaveCalled = true
	return dep.MockProfileSaveError
}

type MockPermissionStore struct {
	MockPermissionSave       *models.Permission
	MockPermissionSaveError  error
	MockPermissionSaveCalled bool
}

func (dep *MockPermissionStore) Save(context.Context, *models.Permission) error {
	dep.MockPermissionSaveCalled = true
	return dep.MockPermissionSaveError
}

type MockRoleStore struct {
	MockRoleSave       *models.Role
	MockRoleSaveError  error
	MockRoleSaveCalled bool
}

func (dep *MockRoleStore) Save(context.Context, *models.Role) error {
	dep.MockRoleSaveCalled = true
	return dep.MockRoleSaveError
}
