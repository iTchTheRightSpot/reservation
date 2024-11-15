package role

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IRoleStore interface {
	Save(ctx context.Context, r *models.Role) (*models.Role, error)
}

type roleStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewRoleStore(l utils.ILogger, db pkg.Db) IRoleStore {
	return &roleStore{logger: l, db: db}
}

func (dep *roleStore) Save(ctx context.Context, r *models.Role) (*models.Role, error) {
	if r == nil {
		return nil, fmt.Errorf("role object is nil")
	}

	q := `
    	INSERT INTO role (role, profile_id)
        VALUES ($1, $2)
        RETURNING role_id, role, profile_id
	`

	row := dep.db.QueryRowContext(ctx, q, r.Role, r.ProfileId)

	err := row.Scan(&r.RoleId, &r.Role, &r.ProfileId)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to role table")
	}

	return r, nil
}
