package profile

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IProfileStore interface {
	Save(ctx context.Context, p *models.ProfileEntity) error
	ProfileByEmail(ctx context.Context, email string) (*models.ProfileEntity, error)
	ProfileRolesAndPermissionByEmail(ctx context.Context, email string) (*models.ProfileRolePermissionEntity, error)
}

type profileStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewProfileStore(l utils.ILogger, db pkg.Db) IProfileStore {
	return &profileStore{logger: l, db: db}
}

func (dep *profileStore) ProfileRolesAndPermissionByEmail(ctx context.Context, email string) (*models.ProfileRolePermissionEntity, error) {
	q := `
        SELECT
          row_to_json(p.*) AS profile,
          json_agg(
              json_build_object(
                  'role', row_to_json(r),
                  'permissions', permissions.permissions
              )
          ) AS role_perm
      FROM profile p
      INNER JOIN role r ON r.profile_id = p.profile_id
      INNER JOIN (
          SELECT
              r.role_id,
              json_agg(perm.*) AS permissions
          FROM permission perm
          INNER JOIN role r ON r.role_id = perm.role_id
          GROUP BY r.role_id
      ) permissions ON permissions.role_id = r.role_id
      WHERE p.email = $1
      GROUP BY p.profile_id;
    `

	rows, err := dep.db.QueryContext(ctx, q, email)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("error retrieving profile, roles, and permissions")
	}

	defer func(rows *sql.Rows) { err = rows.Close() }(rows)

	var result *models.ProfileRolePermissionEntity

	for rows.Next() {
		var profileData json.RawMessage
		var rolePermData json.RawMessage

		if err = rows.Scan(&profileData, &rolePermData); err != nil {
			dep.logger.Error(err.Error())
			return nil, errors.New("error scanning database rows")
		}

		var pro models.ProfileEntity
		if err = json.Unmarshal(profileData, &pro); err != nil {
			dep.logger.Error(err.Error())
			return nil, errors.New("error unmarshalling profile data")
		}

		var rolePerms []models.RolePermissionEntity
		if err = json.Unmarshal(rolePermData, &rolePerms); err != nil {
			dep.logger.Error(err.Error())
			return nil, errors.New("error unmarshalling role permissions data")
		}

		result = &models.ProfileRolePermissionEntity{
			Profile:        pro,
			RolePermission: rolePerms,
		}
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("error iterating through rows")
	}

	if result == nil {
		return nil, errors.New("profile not found")
	}

	return result, nil
}

func (dep *profileStore) ProfileByEmail(ctx context.Context, email string) (*models.ProfileEntity, error) {
	var p models.ProfileEntity

	row := dep.db.QueryRowContext(ctx, "SELECT * FROM profile WHERE email = $1 LIMIT 1", email)
	err := row.Scan(&p.ProfileId, &p.Firstname, &p.Lastname, &p.Email, &p.Password, &p.ImageKey)

	if err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("error retrieving profile by email")
	}

	return &p, nil
}

func (dep *profileStore) Save(ctx context.Context, p *models.ProfileEntity) error {
	if p == nil {
		return errors.New("profile object is nil")
	}

	q := `
		INSERT INTO profile (firstname, lastname, email, password, image_key)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING profile_id, firstname, lastname, email, password, image_key
	`

	row := dep.db.QueryRowContext(ctx, q, p.Firstname, p.Lastname, p.Email, p.Password, p.ImageKey)
	err := row.Scan(&p.ProfileId, &p.Firstname, &p.Lastname, &p.Email, &p.Password, &p.ImageKey)

	if err != nil {
		dep.logger.Error(err.Error())
		return errors.New("error saving to profile table")
	}

	return nil
}
