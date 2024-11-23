package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type serviceRegistry struct {
	JwtService   auth.IJwtService
	ShiftService schedule.IScheduleService
	ServiceImpl  service.IService
}

func newServiceRegistry(s *sql.DB, l utils.ILogger, e *config.SecretVariables) *serviceRegistry {
	a := stores.NewAdapters(l, s, stores.NewTransactionProvider(l, s))
	security := auth.NewJwtService(l, e)

	return &serviceRegistry{
		JwtService:   security,
		ShiftService: schedule.NewScheduleService(l, a),
		ServiceImpl:  service.NewServiceImpl(l, a),
	}
}
