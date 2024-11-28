package stores

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type mockLiveTransactionProvider struct {
	logger utils.ILogger
	tx     *sql.Tx
}

func MockLiveTransactionProvider(l utils.ILogger, tx *sql.Tx) ITransactionProvider {
	return &mockLiveTransactionProvider{logger: l, tx: tx}
}

func (p *mockLiveTransactionProvider) RunInTransaction(txFunc func(adapters *Adapters) error) error {
	return txFunc(NewAdapters(p.logger, p.tx, nil))
}

type MockUnitTransactionProvider struct {
	ProfileStore      profile.IProfileStore
	RoleStore         profile.IRoleStore
	PermissionStore   profile.IPermissionStore
	ScheduleStore     schedule.IScheduleStore
	ServiceStore      service.IServiceStore
	StaffStore        staff.IStaffStore
	StaffServiceStore staff.IStaffServiceStore
	ReservationStore  reservation.IReservationStore
}

func (m *MockUnitTransactionProvider) RunInTransaction(txFunc func(adapters *Adapters) error) error {
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
