package cmd

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
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
	j := auth.NewJwtServiceAsymmetric(l, e)
	//j := auth.NewJwtServiceSymmetric(l, e)

	if e.CookieParam.CookieSecure {
		dummyUser(a, p)
	}

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

func dummyUser(a *stores.Adapters, ps auth.IPasswordService) {
	ctx := context.Background()
	pass, _ := ps.Encode("Password123@#$")
	p := models.ProfileEntity{
		Firstname: "Developer",
		Lastname:  "Lastname",
		Email:     "developer@email.com",
		Password:  string(pass),
	}

	_ = a.ProfileStore.Save(ctx, &p)

	r1 := models.RoleEntity{Role: models.STAFF, ProfileId: p.ProfileId}
	_ = a.RoleStore.Save(ctx, &r1)
	_ = a.PermissionStore.Save(ctx, &models.PermissionEntity{
		Permission: models.READ,
		RoleId:     r1.RoleId,
	})
	_ = a.PermissionStore.Save(ctx, &models.PermissionEntity{
		Permission: models.WRITE,
		RoleId:     r1.RoleId,
	})
	_ = a.PermissionStore.Save(ctx, &models.PermissionEntity{
		Permission: models.DELETE,
		RoleId:     r1.RoleId,
	})

	r2 := models.RoleEntity{Role: models.DEVELOPER, ProfileId: p.ProfileId}
	_ = a.RoleStore.Save(ctx, &r2)
	_ = a.PermissionStore.Save(ctx, &models.PermissionEntity{
		Permission: models.READ,
		RoleId:     r2.RoleId,
	})
	_ = a.PermissionStore.Save(ctx, &models.PermissionEntity{
		Permission: models.WRITE,
		RoleId:     r2.RoleId,
	})
	_ = a.PermissionStore.Save(ctx, &models.PermissionEntity{
		Permission: models.DELETE,
		RoleId:     r2.RoleId,
	})

	lorem := "Lorem ipsum dolor sit amet, consectetur adipisicing elit."
	st := staffModel.StaffEntity{
		UUID:      uuid.New(),
		Bio:       &lorem,
		ProfileId: &p.ProfileId,
	}
	_ = a.StaffStore.Save(ctx, &st)
}
