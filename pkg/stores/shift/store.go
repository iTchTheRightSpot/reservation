package shift

import (
	"context"
	"fmt"
	utils2 "github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/shift"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IShiftStore interface {
	Save(ctx context.Context, s *shift.Shift) (*shift.Shift, error)
}

type profileStore struct {
	logger utils.ILogger
	db     utils2.Db
}

func NewShiftStore(l utils.ILogger, db utils2.Db) IShiftStore {
	return &profileStore{logger: l, db: db}
}

func (dep *profileStore) Save(ctx context.Context, s *shift.Shift) (*shift.Shift, error) {
	if s == nil {
		return nil, fmt.Errorf("shift object is nil")
	}

	q := `
		INSERT INTO shift (start, shift_end, is_enabled, staff_id)
		VALUES ($1, $2, $3, $4)
		RETURNING shift_id, start, shift_end, is_enabled, staff_id
	`

	row := dep.db.QueryRowContext(ctx, q, s.Start, s.End, s.IsEnabled, s.StaffId)
	err := row.Scan(&s.ShiftId, &s.Start, &s.End, &s.IsEnabled, &s.StaffId)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to shift table")
	}

	return s, nil
}
