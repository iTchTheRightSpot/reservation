package staff

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"strings"
)

type IStaffStore interface {
	Save(ctx context.Context, r *staff.Staff) error
	StaffByUUID(ctx context.Context, staffUUID string) (*staff.Staff, error)
	StaffsByServices(ctx context.Context, s *[]string) ([]*staff.StaffStoreFrontDb, error)
	StaffByProfileId(ctx context.Context, id uint64) (*staff.Staff, error)
}

type staffStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewStaffStore(l utils.ILogger, db pkg.Db) IStaffStore {
	return &staffStore{logger: l, db: db}
}

func (dep *staffStore) StaffByProfileId(ctx context.Context, profileId uint64) (*staff.Staff, error) {
	var r staff.Staff
	row := dep.db.QueryRowContext(ctx, "SELECT * FROM staff WHERE profile_id = $1", profileId)
	if err := row.Scan(&r.StaffId, &r.UUID, &r.Bio, &r.ProfileId); err != nil {
		dep.logger.Error(err.Error())
		return nil, fmt.Errorf("no staff with id %v", profileId)
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
		return nil, errors.New("exception retrieving staffs by services")
	}

	defer func(rows *sql.Rows) { err = rows.Close() }(rows)

	var arr []*staff.StaffStoreFrontDb

	for rows.Next() {
		var sb staff.StaffStoreFrontDb

		if err = rows.Scan(&sb.UUID, &sb.Name, &sb.ImageKey, &sb.Bio); err != nil {
			dep.logger.Error(err)
			return nil, errors.New("exception scanning staffs by services")
		}

		arr = append(arr, &sb)
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("error iterating staffs by services")
	}

	return arr, err
}

func (dep *staffStore) Save(ctx context.Context, r *staff.Staff) error {
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

func (dep *staffStore) StaffByUUID(ctx context.Context, staffUUID string) (*staff.Staff, error) {
	var r staff.Staff
	q := "SELECT * FROM staff WHERE uuid = $1"

	row := dep.db.QueryRowContext(ctx, q, staffUUID)
	err := row.Scan(&r.StaffId, &r.UUID, &r.Bio, &r.ProfileId)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, fmt.Errorf("exception retrieving staff with uuid %s", staffUUID)
	}

	return &r, nil
}
