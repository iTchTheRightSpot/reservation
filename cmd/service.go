package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type serviceRegistry struct {
	JwtService         auth.IJwtService
	ScheduleService    schedule.IScheduleService
	ServiceImpl        service.IService
	StaffService       staff.IStaffService
	ReservationService reservation.IReservationService
}

func newServiceRegistry(s *sql.DB, l utils.ILogger, e *config.SecretVariables) *serviceRegistry {
	a := stores.NewAdapters(l, s, stores.NewTransactionProvider(l, s))
	return &serviceRegistry{
		JwtService:         auth.NewJwtService(l, e),
		ScheduleService:    schedule.NewScheduleService(l, a),
		ServiceImpl:        service.NewServiceImpl(l, a),
		StaffService:       staff.NewStaffService(l, a),
		ReservationService: reservation.NewReservationService(l, a),
	}
}
