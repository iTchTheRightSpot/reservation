package staff

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IStaffStore interface {
	Save(ctx context.Context, r *staff.Staff) (*staff.Staff, error)
	StaffByUUID(ctx context.Context, staffUUID string) (*staff.Staff, error)
}

type staffStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewStaffStore(l utils.ILogger, db pkg.Db) IStaffStore {
	return &staffStore{logger: l, db: db}
}

func (dep *staffStore) Save(ctx context.Context, r *staff.Staff) (*staff.Staff, error) {
	if r == nil {
		return nil, fmt.Errorf("staff object is nil")
	}

	q := `
    	INSERT INTO staff (uuid, bio, profile_id)
        VALUES ($1, $2, $3)
        RETURNING staff_id, uuid, bio, profile_id
	`

	row := dep.db.QueryRowContext(ctx, q, r.UUID, r.Bio, r.ProfileId)

	err := row.Scan(&r.StaffId, &r.UUID, &r.Bio, &r.ProfileId)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to staff table")
	}

	return r, nil
}

func (dep *staffStore) StaffByUUID(ctx context.Context, staffUUID string) (*staff.Staff, error) {
	var r staff.Staff
	q := "SELECT * FROM staff WHERE uuid = $1"

	row := dep.db.QueryRowContext(ctx, q, staffUUID)
	err := row.Scan(&r.StaffId, &r.UUID, &r.Bio, &r.ProfileId)
	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception retrieving staff with uuid %s", staffUUID)
	}

	return &r, nil
}
