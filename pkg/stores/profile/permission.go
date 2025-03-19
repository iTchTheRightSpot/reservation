package profile

import (
	"context"
	"errors"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IPermissionStore interface {
	Save(ctx context.Context, p *models.PermissionEntity) error
	Delete(ctx context.Context, permissionId uint64) (int64, error)
}

type permissionStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewPermissionStore(l utils.ILogger, db pkg.Db) IPermissionStore {
	return &permissionStore{logger: l, db: db}
}

func (dep *permissionStore) Save(ctx context.Context, p *models.PermissionEntity) error {
	if p == nil {
		return errors.New("permission object is nil")
	}

	q := `
    	INSERT INTO permission (permission, role_id)
        VALUES ($1, $2)
        RETURNING permission_id, permission, role_id
	`

	row := dep.db.QueryRowContext(ctx, q, p.Permission, p.RoleId)
	if err := row.Scan(&p.PermissionId, &p.Permission, &p.RoleId); err != nil {
		dep.logger.Error(err.Error())
		return &utils.InsertionError{Message: "exception saving to permission table"}
	}

	return nil
}

func (dep *permissionStore) Delete(ctx context.Context, permissionId uint64) (int64, error) {
	res, err := dep.db.ExecContext(ctx, "DELETE FROM permission WHERE permission_id = $1", permissionId)
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, &utils.InsertionError{Message: "error deleting permission"}
	}
	i, err := res.RowsAffected()
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, &utils.InsertionError{Message: "error deleting permission x2"}
	}
	return i, nil
}
