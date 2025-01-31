package account

import (
	"context"
	"errors"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"testing"
)

func TestAccountService(t *testing.T) {
	t.Parallel()

	lg := utils.NewMockLogger()
	ctx := context.Background()

	t.Run("should return error when adding role & permission. Invalid staff id", func(t *testing.T) {
		// given
		ac := accountService{
			logger: lg,
			adp: &stores.Adapters{
				ProfileStore: &profile.MockProfileStore{
					ProfileRolesAndPermissionByStaffUUIDError: errors.New("invalid user"),
				},
			},
		}

		// method to test
		err := ac.AddRoleAndPermission(ctx, &models.AddRoleAndPermissionPayload{UserId: "staff-id"})

		// assert
		if err == nil {
			t.Errorf("expect err, given nil")
		}
		m := "invalid staff id"
		if err.Error() != m {
			t.Errorf("expect '%s', given %s", m, err.Error())
		}
	})

	t.Run("should save role & permission as they do not exist", func(t *testing.T) {
		// given
		payload := models.AddRoleAndPermissionPayload{
			RolePermission: []models.RolePermissionEnum{
				{Role: models.DEVELOPER, Permissions: []models.PermissionEnum{models.READ, models.WRITE}},
			},
		}

		en := models.ProfileRolePermissionEntity{
			Profile: models.ProfileEntity{ProfileId: 1},
			RolePermission: []models.RolePermissionEntity{
				{
					Role: models.RoleEntity{RoleId: 1, Role: models.STAFF},
					Permissions: []models.PermissionEntity{
						{RoleId: 1, Permission: models.READ},
					},
				},
			},
		}

		rStore := profile.MockRoleStore{}
		pStore := profile.MockPermissionStore{}
		ac := accountService{
			logger: lg,
			adp: &stores.Adapters{
				ProfileStore: &profile.MockProfileStore{
					ProfileRolesAndPermissionByStaffUUIDObj:   &en,
					ProfileRolesAndPermissionByStaffUUIDError: nil,
				},
				Transaction: &stores.MockUnitTransactionProvider{
					RoleStore:       &rStore,
					PermissionStore: &pStore,
				},
			},
		}

		// method to test
		err := ac.AddRoleAndPermission(ctx, &payload)

		// assert
		if err != nil {
			t.Errorf("expect nil, given %s", err.Error())
		}

		if !rStore.MockRoleSaveCalled {
			t.Errorf("expect role save to be called, but it was never called")
		}

		if !pStore.MockPermissionSaveCalled {
			t.Errorf("expect permission save to be called, but it was never called")
		}
	})

	t.Run("should only permission as it does not exist", func(t *testing.T) {
		// given
		payload := models.AddRoleAndPermissionPayload{
			RolePermission: []models.RolePermissionEnum{
				{Role: models.STAFF, Permissions: []models.PermissionEnum{models.READ, models.WRITE}},
			},
		}

		en := models.ProfileRolePermissionEntity{
			Profile: models.ProfileEntity{ProfileId: 1},
			RolePermission: []models.RolePermissionEntity{
				{
					Role: models.RoleEntity{RoleId: 1, Role: models.STAFF},
					Permissions: []models.PermissionEntity{
						{RoleId: 1, Permission: models.READ},
					},
				},
			},
		}

		rStore := profile.MockRoleStore{}
		pStore := profile.MockPermissionStore{}
		ac := accountService{
			logger: lg,
			adp: &stores.Adapters{
				ProfileStore: &profile.MockProfileStore{
					ProfileRolesAndPermissionByStaffUUIDObj:   &en,
					ProfileRolesAndPermissionByStaffUUIDError: nil,
				},
				Transaction: &stores.MockUnitTransactionProvider{
					RoleStore:       &rStore,
					PermissionStore: &pStore,
				},
			},
		}

		// method to test
		err := ac.AddRoleAndPermission(ctx, &payload)

		// assert
		if err != nil {
			t.Errorf("expect nil, given %s", err.Error())
		}

		if rStore.MockRoleSaveCalled {
			t.Errorf("expect role save to be called, but it was never called")
		}

		if !pStore.MockPermissionSaveCalled {
			t.Errorf("expect permission save to be called, but it was never called")
		}
	})
}
