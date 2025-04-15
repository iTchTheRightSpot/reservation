package profile

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/utility/utils"
)

type IRoleStore interface {
	Save(ctx context.Context, r *models.RoleEntity) error
	Delete(ctx context.Context, roleId uint64) (int64, error)
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
		return &utils.ServerError{Message: "role object is nil"}
	}

	q := `
    	INSERT INTO role (role, profile_id)
        VALUES ($1, $2)
        RETURNING role_id, role, profile_id
	`

	row := dep.db.QueryRowContext(ctx, q, r.Role, r.ProfileId)
	if err := row.Scan(&r.RoleId, &r.Role, &r.ProfileId); err != nil {
		dep.logger.Error(ctx, err.Error())
		return &utils.InsertionError{}
	}

	return nil
}

func (dep *roleStore) Delete(ctx context.Context, roleId uint64) (int64, error) {
	res, err := dep.db.ExecContext(ctx, "DELETE FROM role WHERE role_id = $1", roleId)
	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return 0, &utils.InsertionError{Message: "error deleting role"}
	}
	i, err := res.RowsAffected()
	if err != nil {
		return 0, &utils.InsertionError{Message: "error deleting role no."}
	}
	return i, nil
}