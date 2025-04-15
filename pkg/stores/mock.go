package stores

import (
	"context"
	"github.com/iTchTheRightSpot/reservation/pkg/stores/profile"
	"github.com/iTchTheRightSpot/reservation/pkg/stores/reservation"
	"github.com/iTchTheRightSpot/reservation/pkg/stores/schedule"
	"github.com/iTchTheRightSpot/reservation/pkg/stores/service_type"
	"github.com/iTchTheRightSpot/reservation/pkg/stores/staff"
)

type MockUnitTransactionProvider struct {
	ProfileStore      profile.IProfileStore
	RoleStore         profile.IRoleStore
	PermissionStore   profile.IPermissionStore
	ScheduleStore     schedule.IScheduleStore
	ServiceStore      service_type.IServiceTypeStore
	StaffStore        staff.IStaffStore
	StaffServiceStore staff.IStaffServiceStore
	ReservationStore  reservation.IReservationStore
}

func (m *MockUnitTransactionProvider) RunInTransaction(_ context.Context, txFunc func(adapters *Adapters) error) error {
	return txFunc(&Adapters{
		ProfileStore:      m.ProfileStore,
		RoleStore:         m.RoleStore,
		PermissionStore:   m.PermissionStore,
		ScheduleStore:     m.ScheduleStore,
		ServiceStore:      m.ServiceStore,
		StaffStore:        m.StaffStore,
		StaffServiceStore: m.StaffServiceStore,
		ReservationStore:  m.ReservationStore,
	})
}