package account

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
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
	if _, err := dep.adp.ProfileStore.ProfileByEmail(ctx, obj.Email); err == nil {
		return errors.New("email already used")
	}

	pass, err := dep.ps.Encode(obj.Password)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.BadRequestError{Message: "error encoding password. please try a different password"}
	}

	return dep.adp.Transaction.RunInTransaction(func(adps *stores.Adapters) error {
		p := profile.Profile{
			Firstname: obj.Firstname,
			Lastname:  obj.Lastname,
			Email:     obj.Email,
			Password:  pass,
		}
		if err = adps.ProfileStore.Save(ctx, &p); err != nil {
			dep.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error creating account. profile"}
		}

		r := models.Role{Role: models.STAFF, ProfileId: p.ProfileId}
		if err = adps.RoleStore.Save(ctx, &r); err != nil {
			dep.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error creating account. role"}
		}

		if err = adps.PermissionStore.Save(ctx, &models.Permission{Permission: models.READ, RoleId: r.RoleId}); err != nil {
			dep.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error creating account. permission"}
		}
		bio := "Ready to put a smile on your face 🌞"
		return adps.StaffStore.Save(ctx, &staff.Staff{ProfileId: &p.ProfileId, UUID: uuid.New(), Bio: &bio})
	})
}
