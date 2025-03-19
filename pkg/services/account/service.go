package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"slices"
)

type IAccountService interface {
	Register(ctx context.Context, obj *models.ProfilePayload) error
	Login(ctx context.Context, obj *models.Login) (*models.JwtResponse, error)
	ActiveUser(ctx context.Context, obj *models.JwtObj) (*models.ActiveUser, error)
	AddRoleAndPermission(ctx context.Context, o *models.RoleAndPermissionPayload) error
	DeleteRole(ctx context.Context, staffUUID, role string) error
	DeletePermission(ctx context.Context, staffUUID, role, permission string) error
}

type accountService struct {
	logger     utils.ILogger
	adp        *stores.Adapters
	ps         auth.IPasswordService
	jwt        auth.IJwtService
	staffCache pkg.ICache[string, []*staff.AllStaffsEntity]
}

func NewAccountService(l utils.ILogger, a *stores.Adapters, j auth.IJwtService, ps auth.IPasswordService, c pkg.ICache[string, []*staff.AllStaffsEntity]) IAccountService {
	return &accountService{logger: l, adp: a, jwt: j, ps: ps, staffCache: c}
}

func (dep *accountService) ActiveUser(ctx context.Context, obj *models.JwtObj) (*models.ActiveUser, error) {
	s, err := dep.adp.ProfileStore.ProfileByStaffUUID(ctx, obj.UserId)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, err
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
		return nil, err
	}

	staf, err := dep.adp.StaffStore.StaffByProfileId(ctx, prp.Profile.ProfileId)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, err
	}

	acs := dep.accessControls(prp)

	return dep.jwt.Encode(&models.JwtObj{UserId: staf.UUID.String(), AccessControls: *acs}, utils.TwoDaysInSeconds)
}

func (dep *accountService) Register(ctx context.Context, obj *models.ProfilePayload) error {
	if _, err := dep.adp.ProfileStore.ProfileByEmail(ctx, obj.Email); err == nil {
		return err
	}

	pass, err := dep.ps.Encode(obj.Password)
	if err != nil {
		dep.logger.Error(err.Error())
		return err
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
			return err
		}

		r := models.RoleEntity{Role: models.STAFF, ProfileId: p.ProfileId}
		if err = adps.RoleStore.Save(ctx, &r); err != nil {
			dep.logger.Error(err.Error())
			return err
		}

		if err = adps.PermissionStore.Save(ctx, &models.PermissionEntity{Permission: models.READ, RoleId: r.RoleId}); err != nil {
			dep.logger.Error(err.Error())
			return err
		}

		bio := "Ready to put a smile on your face 🌞"
		err = adps.StaffStore.Save(ctx, &staff.StaffEntity{ProfileId: &p.ProfileId, UUID: uuid.New(), Bio: &bio})
		if err != nil {
			dep.logger.Error(err.Error())
			return err
		}

		dep.staffCache.Clear()
		return err
	})
}

func (dep *accountService) AddRoleAndPermission(ctx context.Context, o *models.RoleAndPermissionPayload) error {
	prp, err := dep.adp.ProfileStore.ProfileRolesAndPermissionByStaffUUID(ctx, o.UserId)
	if err != nil {
		dep.logger.Error(err.Error())
		return err
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
							return err
						}
					}
				}
			} else {
				// save normally
				r := models.RoleEntity{ProfileId: prp.Profile.ProfileId, Role: enum.Role}
				err = adps.RoleStore.Save(ctx, &r)

				if err != nil {
					dep.logger.Error(err.Error())
					return err
				}

				for _, p := range enum.Permissions {
					err = adps.PermissionStore.Save(ctx, &models.PermissionEntity{RoleId: r.RoleId, Permission: p})

					if err != nil {
						dep.logger.Error(err.Error())
						return err
					}
				}
			}
		}

		dep.staffCache.Clear()
		return nil
	})
}

func (dep *accountService) DeleteRole(ctx context.Context, staffUUID, role string) error {
	prp, err := dep.adp.ProfileStore.ProfileRolesAndPermissionByStaffUUID(ctx, staffUUID)
	if err != nil {
		dep.logger.Error(err.Error())
		return err
	}

	idx := slices.IndexFunc(prp.RolePermission, func(e models.RolePermissionEntity) bool {
		return string(e.Role.Role) == role
	})

	if idx == -1 {
		return err
	}

	return dep.adp.Transaction.RunInTransaction(func(adps *stores.Adapters) error {
		c, err := adps.RoleStore.Delete(ctx, prp.RolePermission[idx].Role.RoleId)
		if err != nil {
			dep.logger.Error(err.Error())
			return err
		}
		dep.logger.Log("number of roles affected after role deletion", c)
		dep.staffCache.Clear()
		return nil
	})
}

func (dep *accountService) DeletePermission(ctx context.Context, staffUUID, role, permission string) error {
	prp, err := dep.adp.ProfileStore.ProfileRolesAndPermissionByStaffUUID(ctx, staffUUID)
	if err != nil {
		dep.logger.Error(err.Error())
		return err
	}

	idx := slices.IndexFunc(prp.RolePermission, func(e models.RolePermissionEntity) bool {
		return string(e.Role.Role) == role
	})

	if idx == -1 {
		return &utils.NotFoundError{Message: "role does not exist"}
	}

	idx1 := slices.IndexFunc(prp.RolePermission[idx].Permissions, func(e models.PermissionEntity) bool {
		return string(e.Permission) == permission
	})

	if idx1 == -1 {
		return &utils.NotFoundError{Message: "permission does not exist"}
	}

	c, err := dep.adp.PermissionStore.Delete(ctx, prp.RolePermission[idx].Permissions[idx1].PermissionId)
	if err != nil {
		dep.logger.Error(err.Error())
		return err
	}

	dep.staffCache.Clear()
	dep.logger.Log("number of roles affected after permission deletion", c)
	return nil
}
