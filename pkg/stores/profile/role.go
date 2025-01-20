package profile

import (
	"context"
	"errors"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IRoleStore interface {
	Save(ctx context.Context, r *models.RoleEntity) error
}

type roleStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewRoleStore(l utils.ILogger, db pkg.Db) IRoleStore {
	return &roleStore{logger: l, db: db}
}

func (dep *roleStore) Save(ctx context.Context, r *models.RoleEntity) error {
	if r == nil {
		return errors.New("role object is nil")
	}

	q := `
    	INSERT INTO role (role, profile_id)
        VALUES ($1, $2)
        RETURNING role_id, role, profile_id
	`

	row := dep.db.QueryRowContext(ctx, q, r.Role, r.ProfileId)

	if err := row.Scan(&r.RoleId, &r.Role, &r.ProfileId); err != nil {
		dep.logger.Error(err.Error())
		return errors.New("exception saving to role table")
	}

	return nil
}
