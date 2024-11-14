package staff

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IStaffStore interface {
	Save(ctx context.Context, r *staff.Staff) (*staff.Staff, error)
}

type staffStore struct {
	logger utils.ILogger
	db     utils.Db
}

func NewStaffStore(l utils.ILogger, db utils.Db) IStaffStore {
	return &staffStore{logger: l, db: db}
}

func (dep *staffStore) Save(ctx context.Context, r *staff.Staff) (*staff.Staff, error) {
	if r == nil {
		return nil, fmt.Errorf("staff object is nil")
	}

	if r.ProfileId == nil {
		return nil, fmt.Errorf("profile_id cannot be nil to save a staff")
	}

	q := `
    	INSERT INTO staff (staff_uuid, bio, profile_id)
        VALUES ($1, $2, $3)
        RETURNING staff_id, staff_uuid, bio, profile_id
	`

	row := dep.db.QueryRowContext(ctx, q, r.StaffUUID, r.Bio, r.ProfileId)

	err := row.Scan(&r.StaffId, &r.StaffUUID, &r.Bio, &r.ProfileId)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to staff table")
	}

	return r, nil
}
