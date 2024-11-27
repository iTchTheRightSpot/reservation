package stores

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type Adapters struct {
	ProfileStore      profile.IProfileStore
	RoleStore         profile.IRoleStore
	PermissionStore   profile.IPermissionStore
	ScheduleStore     schedule.IScheduleStore
	ServiceStore      service.IServiceStore
	StaffStore        staff.IStaffStore
	StaffServiceStore staff.IStaffServiceStore
	ReservationStore  reservation.IReservationStore
	Transaction       ITransactionProvider
}

func NewAdapters(l utils.ILogger, db pkg.Db, tx ITransactionProvider) *Adapters {
	return &Adapters{
		ProfileStore:      profile.NewProfileStore(l, db),
		RoleStore:         profile.NewRoleStore(l, db),
		PermissionStore:   profile.NewPermissionStore(l, db),
		ScheduleStore:     schedule.NewScheduleStore(l, db),
		ServiceStore:      service.NewServiceStore(l, db),
		StaffStore:        staff.NewStaffStore(l, db),
		StaffServiceStore: staff.NewStaffServiceStore(l, db),
		ReservationStore:  reservation.NewReservationStore(l, db),
		Transaction:       tx,
	}
}
