package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/mail"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/service_type"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type serviceRegistry struct {
	JwtService         auth.IJwtService
	ScheduleService    schedule.IScheduleService
	ServiceImpl        service_type.IServiceType
	StaffService       staff.IStaffService
	ReservationService reservation.IReservationService
	MailService        mail.IMailService
}

func newServiceRegistry(s *sql.DB, l utils.ILogger, e *config.SecretVariables) *serviceRegistry {
	a := stores.NewAdapters(l, s, stores.NewTransactionProvider(l, s))
	m := mail.NewMailService(l, e)

	return &serviceRegistry{
		JwtService:      auth.NewJwtService(l, e),
		ScheduleService: schedule.NewScheduleService(l, a),
		ServiceImpl:     service_type.NewServiceImpl(l, a),
		StaffService:    staff.NewStaffService(l, a),
		ReservationService: reservation.NewReservationService(
			l,
			a,
			pkg.NewInMemoryCache[string, []*model.ReservationTimeSlots](l, 30, 30),
			m,
		),
		MailService: m,
	}
}
