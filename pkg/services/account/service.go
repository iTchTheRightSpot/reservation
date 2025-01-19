package account

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IAccountService interface {
	Register(ctx context.Context, obj *profile.ProfilePayload) error
}

type accountService struct {
	logger utils.ILogger
	adp    *stores.Adapters
	ps     auth.IPasswordService
}

func NewAccountService(l utils.ILogger, a *stores.Adapters, ps auth.IPasswordService) IAccountService {
	return &accountService{logger: l, adp: a, ps: ps}
}

func (dep *accountService) Register(ctx context.Context, obj *profile.ProfilePayload) error {
	return nil
}
