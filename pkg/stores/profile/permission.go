package profile

import (
	"context"
	"github.com/iTchTheRightSpot/reservation/pkg"
	"github.com/iTchTheRightSpot/reservation/pkg/models"
	log "github.com/iTchTheRightSpot/utility/utils"
)

type IPermissionStore interface {
	Save(ctx context.Context, p *models.PermissionEntity) error
	Delete(ctx context.Context, permissionId uint64) (int64, error)
}

type permissionStore struct {
	logger log.ILogger
	db     pkg.Db
}

func NewPermissionStore(l log.ILogger, db pkg.Db) IPermissionStore {
	return &permissionStore{logger: l, db: db}
}

func (dep *permissionStore) Save(ctx context.Context, p *models.PermissionEntity) error {
	if p == nil {
		return &log.ServerError{Message: "permission object is nil"}
	}

	q := `
    	INSERT INTO permission (permission, role_id)
        VALUES ($1, $2)
        RETURNING permission_id, permission, role_id
	`

	row := dep.db.QueryRowContext(ctx, q, p.Permission, p.RoleId)
	if err := row.Scan(&p.PermissionId, &p.Permission, &p.RoleId); err != nil {
		dep.logger.Error(ctx, err.Error())
		return &log.InsertionError{Message: "error saving to permission"}
	}

	return nil
}

func (dep *permissionStore) Delete(ctx context.Context, permissionId uint64) (int64, error) {
	res, err := dep.db.ExecContext(ctx, "DELETE FROM permission WHERE permission_id = $1", permissionId)
	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return 0, &log.InsertionError{Message: "error deleting permission"}
	}
	i, err := res.RowsAffected()
	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return 0, &log.InsertionError{Message: "error rows affected"}
	}
	return i, nil
}