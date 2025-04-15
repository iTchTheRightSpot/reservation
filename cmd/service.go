package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/reservation/config"
	model "github.com/iTchTheRightSpot/reservation/pkg/models/reservation"
	staffModel "github.com/iTchTheRightSpot/reservation/pkg/models/staff"
	"github.com/iTchTheRightSpot/reservation/pkg/services/account"
	"github.com/iTchTheRightSpot/reservation/pkg/services/auth"
	"github.com/iTchTheRightSpot/reservation/pkg/services/mail"
	"github.com/iTchTheRightSpot/reservation/pkg/services/reservation"
	"github.com/iTchTheRightSpot/reservation/pkg/services/schedule"
	"github.com/iTchTheRightSpot/reservation/pkg/services/service_type"
	"github.com/iTchTheRightSpot/reservation/pkg/services/staff"
	"github.com/iTchTheRightSpot/reservation/pkg/stores"
	"github.com/iTchTheRightSpot/utility/cache"
	log "github.com/iTchTheRightSpot/utility/utils"
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

func newServiceRegistry(s *sql.DB, l log.ILogger, e *config.SecretVariables) *serviceRegistry {
	a := stores.NewAdapters(l, s, stores.NewTransactionProvider(l, s))
	m := mail.NewMailService(l, e)
	p := auth.NewPasswordService(l)
	j := auth.NewJwtServiceAsymmetric(l, e)

	staffCache := cache.SyncMapInMemoryCache[string, []*staffModel.AllStaffsEntity](l, 10, 10)
	return &serviceRegistry{
		JwtService:         j,
		ScheduleService:    schedule.NewScheduleService(l, a),
		ServiceImpl:        service_type.NewServiceImpl(l, a),
		StaffService:       staff.NewStaffService(l, a, staffCache),
		ReservationService: reservation.NewReservationService(l, a, cache.SyncMapInMemoryCache[string, []*model.ReservationTimeSlots](l, 30, 30), m),
		AccountService:     account.NewAccountService(l, a, j, p, staffCache),
		PasswordService:    p,
		MailService:        m,
	}
}