package account

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IAccountService interface {
	Register(ctx context.Context, obj *models.ProfilePayload) error
	Login(ctx context.Context, obj *models.Login) (*models.JwtResponse, error)
}

type accountService struct {
	logger utils.ILogger
	adp    *stores.Adapters
	ps     auth.IPasswordService
	jwt    auth.IJwtService
}

func NewAccountService(l utils.ILogger, a *stores.Adapters, j auth.IJwtService, ps auth.IPasswordService) IAccountService {
	return &accountService{logger: l, adp: a, jwt: j, ps: ps}
}

func (dep *accountService) accessControls(o *models.ProfileRolePermissionEntity) *[]models.RolePermissionEnum {
	arr := make([]models.RolePermissionEnum, len(o.RolePermission))

	for i, rp := range o.RolePermission {
		rpp := make([]models.PermissionEnum, len(rp.Permissions))

		for j, perm := range rp.Permissions {
			rpp[j] = perm.Permission
		}
		arr[i] = models.RolePermissionEnum{
			Role:        rp.Role.Role,
			Permissions: rpp,
		}
	}

	return &arr
}

// Login TODO implement logic to temporary lock account on brute force
// TODO validate account is not locked
func (dep *accountService) Login(ctx context.Context, obj *models.Login) (*models.JwtResponse, error) {
	prp, err := dep.adp.ProfileStore.ProfileRolesAndPermissionByEmail(ctx, obj.Email)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "invalid email or password"}
	}

	if err = dep.ps.Validate([]byte(prp.Profile.Password), []byte(obj.Password)); err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "invalid email or password"}
	}

	staf, err := dep.adp.StaffStore.StaffByProfileId(ctx, prp.Profile.ProfileId)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("invalid email or password")
	}

	acs := dep.accessControls(prp)

	return dep.jwt.Encode(&models.JwtObj{UserId: staf.UUID.String(), AccessControls: *acs}, utils.TwoDaysInSeconds)
}

func (dep *accountService) Register(ctx context.Context, obj *models.ProfilePayload) error {
	if _, err := dep.adp.ProfileStore.ProfileByEmail(ctx, obj.Email); err == nil {
		return errors.New("email already used")
	}

	pass, err := dep.ps.Encode(obj.Password)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.BadRequestError{Message: "error encoding password. please try a different password"}
	}

	return dep.adp.Transaction.RunInTransaction(func(adps *stores.Adapters) error {
		p := models.ProfileEntity{
			Firstname: obj.Firstname,
			Lastname:  obj.Lastname,
			Email:     obj.Email,
			Password:  string(pass),
		}
		if err = adps.ProfileStore.Save(ctx, &p); err != nil {
			dep.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error creating account. profile"}
		}

		r := models.RoleEntity{Role: models.STAFF, ProfileId: p.ProfileId}
		if err = adps.RoleStore.Save(ctx, &r); err != nil {
			dep.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error creating account. role"}
		}

		if err = adps.PermissionStore.Save(ctx, &models.PermissionEntity{Permission: models.READ, RoleId: r.RoleId}); err != nil {
			dep.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error creating account. permission"}
		}
		bio := "Ready to put a smile on your face 🌞"
		return adps.StaffStore.Save(ctx, &staff.Staff{ProfileId: &p.ProfileId, UUID: uuid.New(), Bio: &bio})
	})
}
