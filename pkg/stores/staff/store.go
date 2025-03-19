package staff

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"strings"
)

type IStaffStore interface {
	Save(ctx context.Context, r *staff.StaffEntity) error
	StaffByUUID(ctx context.Context, staffUUID string) (*staff.StaffEntity, error)
	StaffsByServices(ctx context.Context, s *[]string) ([]*staff.StaffStoreFrontDb, error)
	StaffByProfileId(ctx context.Context, id uint64) (*staff.StaffEntity, error)
	AllStaffs(ctx context.Context) ([]*staff.AllStaffsEntity, error)
}

type staffStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewStaffStore(l utils.ILogger, db pkg.Db) IStaffStore {
	return &staffStore{logger: l, db: db}
}

func (dep *staffStore) AllStaffs(ctx context.Context) ([]*staff.AllStaffsEntity, error) {
	q := `
        SELECT
          p.firstname,
          p.lastname,
          p.email,
          p.locked,
          p.image_key,
          st.uuid as user_id,
          st.bio,
          json_agg(
              json_build_object(
                  'role', r.role,
                  'permissions', permissions.permissions
              )
          ) AS access_controls
        FROM profile p
        INNER JOIN staff st ON st.profile_id = p.profile_id
		INNER JOIN role r ON r.profile_id = p.profile_id
		INNER JOIN (
			SELECT
				r.role_id,
				json_agg(perm.permission) AS permissions
			FROM permission perm
			INNER JOIN role r ON r.role_id = perm.role_id
			GROUP BY r.role_id
		) permissions ON permissions.role_id = r.role_id
		GROUP BY p.firstname, p.lastname, p.email, p.locked, p.image_key, st.uuid, st.bio;
    `

	rows, err := dep.db.QueryContext(ctx, q)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "error retrieving all staffs"}
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			dep.logger.Error(err.Error())
			err = &utils.ServerError{Message: "error closing stream after retrieving all staffs"}
		}
	}(rows)

	var result = make([]*staff.AllStaffsEntity, 0)

	for rows.Next() {
		var pd staff.AllStaffsEntity
		var rolePermData json.RawMessage

		if err = rows.Scan(
			&pd.Firstname,
			&pd.Lastname,
			&pd.Email,
			&pd.Locked,
			&pd.ImageKey,
			&pd.UUID,
			&pd.Bio,
			&rolePermData,
		); err != nil {
			dep.logger.Error(err.Error())
			return nil, &utils.ServerError{Message: "error scanning database rows after retrieving all staffs"}
		}

		if err = json.Unmarshal(rolePermData, &pd.AccessControls); err != nil {
			dep.logger.Error(err.Error())
			return nil, &utils.ServerError{Message: "error unmarshalling all staffs"}
		}

		result = append(result, &pd)
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.ServerError{Message: "error iterating through all staffs"}
	}

	return result, err
}

func (dep *staffStore) StaffByProfileId(ctx context.Context, profileId uint64) (*staff.StaffEntity, error) {
	var r staff.StaffEntity
	row := dep.db.QueryRowContext(ctx, "SELECT * FROM staff WHERE profile_id = $1", profileId)
	if err := row.Scan(&r.StaffId, &r.UUID, &r.Bio, &r.ProfileId); err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "invalid staff with id"}
	}
	return &r, nil
}

func (dep *staffStore) StaffsByServices(ctx context.Context, s *[]string) ([]*staff.StaffStoreFrontDb, error) {
	if s == nil {
		return nil, errors.New("services arr is nil")
	}

	// dynamically construct the placeholder part of the query (e.g., $1, $2, $3 for each service)
	placeholders := make([]string, len(*s))
	for i := range *s {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// join placeholders into a string (e.g., $1, $2, $3)
	placeholderStr := strings.Join(placeholders, ", ")

	q := fmt.Sprintf(`
    	SELECT
			st.uuid, p.firstname, p.image_key, st.bio
		FROM staff st
    	INNER JOIN profile p ON p.profile_id = st.profile_id
    	INNER JOIN staff_service sts ON sts.staff_id = st.staff_id
		INNER JOIN service_type s ON s.service_id = sts.service_id
    	WHERE s.name IN (%s)
    	GROUP BY st.staff_id, st.uuid, p.firstname, p.image_key, st.bio
    	HAVING COUNT(DISTINCT s.name) = %d
	`, placeholderStr, len(*s))

	// flatten the slice to pass it as individual arguments
	args := make([]interface{}, len(*s))
	for i, service := range *s {
		args[i] = service
	}

	rows, err := dep.db.QueryContext(ctx, q, args...)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.NotFoundError{Message: "error retrieving staffs by services"}
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			dep.logger.Error(err.Error())
			err = &utils.ServerError{Message: "error closing stream after retrieving staffs by services"}
		}
	}(rows)

	var arr []*staff.StaffStoreFrontDb

	for rows.Next() {
		var sb staff.StaffStoreFrontDb

		if err = rows.Scan(&sb.UUID, &sb.Name, &sb.ImageKey, &sb.Bio); err != nil {
			dep.logger.Error(err)
			return nil, &utils.ServerError{Message: "exception scanning staffs by services"}
		}

		arr = append(arr, &sb)
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(err.Error())
		return nil, &utils.ServerError{Message: "error iterating staffs by services"}
	}

	return arr, err
}

func (dep *staffStore) Save(ctx context.Context, r *staff.StaffEntity) error {
	if r == nil {
		return errors.New("staff object is nil")
	}

	q := `
    	INSERT INTO staff (uuid, bio, profile_id)
        VALUES ($1, $2, $3)
        RETURNING staff_id, uuid, bio, profile_id
	`

	row := dep.db.QueryRowContext(ctx, q, r.UUID, r.Bio, r.ProfileId)

	err := row.Scan(&r.StaffId, &r.UUID, &r.Bio, &r.ProfileId)

	if err != nil {
		dep.logger.Error(err.Error())
		return errors.New("exception saving to staff table")
	}

	return nil
}

func (dep *staffStore) StaffByUUID(ctx context.Context, staffUUID string) (*staff.StaffEntity, error) {
	var r staff.StaffEntity
	q := "SELECT * FROM staff WHERE uuid = $1"

	row := dep.db.QueryRowContext(ctx, q, staffUUID)
	err := row.Scan(&r.StaffId, &r.UUID, &r.Bio, &r.ProfileId)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, fmt.Errorf("exception retrieving staff with uuid %s", staffUUID)
	}

	return &r, nil
}
