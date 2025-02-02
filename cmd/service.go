package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	staffModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/account"
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
	PasswordService    auth.IPasswordService
	AccountService     account.IAccountService
}

func newServiceRegistry(s *sql.DB, l utils.ILogger, e *config.SecretVariables) *serviceRegistry {
	a := stores.NewAdapters(l, s, stores.NewTransactionProvider(l, s))
	m := mail.NewMailService(l, e)
	p := auth.NewPasswordService(l)
	j := auth.NewJwtService(l, e)

	staffCache := pkg.NewInMemoryCache[string, []*staffModel.AllStaffsEntity](l, 10, 10)
	return &serviceRegistry{
		JwtService:         j,
		ScheduleService:    schedule.NewScheduleService(l, a),
		ServiceImpl:        service_type.NewServiceImpl(l, a),
		StaffService:       staff.NewStaffService(l, a, staffCache),
		ReservationService: reservation.NewReservationService(l, a, pkg.NewInMemoryCache[string, []*model.ReservationTimeSlots](l, 30, 30), m),
		AccountService:     account.NewAccountService(l, a, j, p, staffCache),
		PasswordService:    p,
		MailService:        m,
	}
}
