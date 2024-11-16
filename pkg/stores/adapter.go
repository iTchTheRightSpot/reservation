package stores

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/shift"
	staffStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type Adapters struct {
	ProfileStore    profile.IProfileStore
	RoleStore       profile.IRoleStore
	PermissionStore profile.IPermissionStore
	ShiftStore      shift.IShiftStore
	StaffStore      staffStore.IStaffStore
	Transaction     ITransactionProvider
}

func NewAdapters(l utils.ILogger, db pkg.Db, p ITransactionProvider) *Adapters {
	return &Adapters{
		ProfileStore:    profile.NewProfileStore(l, db),
		RoleStore:       profile.NewRoleStore(l, db),
		PermissionStore: profile.NewPermissionStore(l, db),
		ShiftStore:      shift.NewShiftStore(l, db),
		StaffStore:      staffStore.NewStaffStore(l, db),
		Transaction:     p,
	}
}
