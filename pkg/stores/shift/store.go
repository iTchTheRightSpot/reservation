package shift

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/shift"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"time"
)

type IShiftStore interface {
	Save(ctx context.Context, s *shift.Shift) (*shift.Shift, error)
	CountExistingShiftsForStaff(ctx context.Context, staffId uint64, start, end time.Time) (int, error)
}

type profileStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewShiftStore(l utils.ILogger, db pkg.Db) IShiftStore {
	return &profileStore{logger: l, db: db}
}

func (dep *profileStore) Save(ctx context.Context, s *shift.Shift) (*shift.Shift, error) {
	if s == nil {
		return nil, fmt.Errorf("shift object is nil")
	}

	q := `
		INSERT INTO shift (shift_start, shift_end, is_enabled, is_reoccurring, staff_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING shift_id, shift_start, shift_end, is_enabled, is_reoccurring, staff_id
	`

	row := dep.db.QueryRowContext(ctx, q, s.Start, s.End, s.IsEnabled, s.IsReoccurring, s.StaffId)
	err := row.Scan(&s.ShiftId, &s.Start, &s.End, &s.IsEnabled, &s.IsReoccurring, &s.StaffId)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to shift table")
	}

	return s, nil
}

func (dep *profileStore) CountExistingShiftsForStaff(ctx context.Context, staffId uint64, start, end time.Time) (int, error) {
	q := `
		SELECT COUNT(s.shift_id) FROM shift s
		WHERE s.staff_id = $1
		AND (
			(s.shift_start BETWEEN $2 AND $3) OR
			(s.shift_end BETWEEN $2 AND $3)
		)
	`

	count := -1
	row := dep.db.QueryRowContext(ctx, q, staffId, start, end)

	if err := row.Scan(&count); err != nil {
		dep.logger.Error(err)
		return 0, err
	} else if count == -1 {
		dep.logger.Error(fmt.Sprintf("exception counting existing shifts for staff id %v", staffId))
		return 0, fmt.Errorf("exception counting existing shifts for staff")
	}

	return count, nil
}
