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
	"slices"
)

type IAccountService interface {
	Register(ctx context.Context, obj *models.ProfilePayload) error
	Login(ctx context.Context, obj *models.Login) (*models.JwtResponse, error)
	ActiveUser(ctx context.Context, obj *models.JwtObj) (*models.ActiveUser, error)
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

func (dep *accountService) ActiveUser(ctx context.Context, obj *models.JwtObj) (*models.ActiveUser, error) {
	s, err := dep.adp.ProfileStore.ProfileByStaffUUID(ctx, obj.UserId)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "invalid user id"}
	}
	return &models.ActiveUser{
		UserId:         obj.UserId,
		Firstname:      s.Firstname,
		ImageKey:       s.ImageKey,
		AccessControls: obj.AccessControls,
	}, nil
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

func (dep *accountService) Login(ctx context.Context, obj *models.Login) (*models.JwtResponse, error) {
	prp, err := dep.adp.ProfileStore.ProfileRolesAndPermissionByEmail(ctx, obj.Email)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "invalid email or password"}
	}

	if prp.Profile.Locked {
		return nil, &utils.AuthenticationError{Message: "account locked. Please reset your password"}
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
		return adps.StaffStore.Save(ctx, &staff.StaffEntity{ProfileId: &p.ProfileId, UUID: uuid.New(), Bio: &bio})
	})
}

func (dep *accountService) AddRoleAndPermission(ctx context.Context, o *models.AddRoleAndPermissionPayload) error {
	prp, err := dep.adp.ProfileStore.ProfileRolesAndPermissionByStaffUUID(ctx, o.UserId)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.NotFoundError{Message: "invalid staff id"}
	}

	return dep.adp.Transaction.RunInTransaction(func(adps *stores.Adapters) error {
		for _, enum := range o.RolePermission {
			idx := slices.IndexFunc(prp.RolePermission, func(e models.RolePermissionEntity) bool {
				return e.Role.Role == enum.Role
			})

			// validate if role does not exist
			if idx != -1 {
				rp := prp.RolePermission[idx]
				for _, p := range enum.Permissions {
					e := slices.ContainsFunc(rp.Permissions, func(pe models.PermissionEntity) bool {
						return pe.Permission == p
					})

					if !e {
						err = adps.PermissionStore.Save(ctx, &models.PermissionEntity{RoleId: rp.Role.RoleId, Permission: p})
						if err != nil {
							dep.logger.Error(err.Error())
							return &utils.InsertionError{Message: "error saving permission " + string(p)}
						}
					}
				}
			} else {
				// save normally
				r := models.RoleEntity{ProfileId: prp.Profile.ProfileId, Role: enum.Role}
				err = adps.RoleStore.Save(ctx, &r)

				if err != nil {
					dep.logger.Error(err.Error())
					return &utils.InsertionError{Message: "error saving role " + string(enum.Role)}
				}

				for _, p := range enum.Permissions {
					err = adps.PermissionStore.Save(ctx, &models.PermissionEntity{RoleId: r.RoleId, Permission: p})

					if err != nil {
						dep.logger.Error(err.Error())
						return &utils.InsertionError{Message: "error saving permission " + string(p)}
					}
				}
			}
		}
		return nil
	})
}
