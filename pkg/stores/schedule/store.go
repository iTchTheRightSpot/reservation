package schedule

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"time"
)

type IScheduleStore interface {
	Save(ctx context.Context, s *schedule.Schedule) (*schedule.Schedule, error)
	CountExistingSchedulesForStaff(ctx context.Context, staffId uint64, start, end time.Time) (int, error)
	SchedulesInRange(ctx context.Context, staffId uint64, start time.Time, end time.Time) ([]*schedule.Schedule, error)
	CountSchedulesInRangeAndVisibility(ctx context.Context, staffId uint64, start time.Time, end time.Time, isVisible bool) (int, error)
}

type scheduleStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewScheduleStore(l utils.ILogger, db pkg.Db) IScheduleStore {
	return &scheduleStore{logger: l, db: db}
}

func (dep *scheduleStore) SchedulesInRange(ctx context.Context, staffId uint64, start time.Time, end time.Time) ([]*schedule.Schedule, error) {
	var arr []*schedule.Schedule

	q := `
		SELECT * FROM schedule s
		WHERE s.staff_id = $1
		AND (
			(s.schedule_start BETWEEN $2 AND $3) OR
			(s.schedule_end BETWEEN $2 AND $3)
		)
	`

	rows, err := dep.db.QueryContext(ctx, q, staffId, start, end)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, fmt.Errorf("error retrieve schedules in range")
	}

	defer func(rs *sql.Rows) { err = rs.Close() }(rows)

	for rows.Next() {
		var s schedule.Schedule

		err = rows.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)
		if err != nil {
			dep.logger.Error(err.Error())
			return nil, fmt.Errorf("error scanning schedules")
		}

		arr = append(arr, &s)
	}

	return arr, err
}

func (dep *scheduleStore) Save(ctx context.Context, s *schedule.Schedule) (*schedule.Schedule, error) {
	if s == nil {
		return nil, fmt.Errorf("schedule object is nil")
	}

	q := `
		INSERT INTO schedule (schedule_start, schedule_end, is_visible, is_reoccurring, staff_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING schedule_id, schedule_start, schedule_end, is_visible, is_reoccurring, staff_id
	`

	row := dep.db.QueryRowContext(ctx, q, s.Start, s.End, s.IsVisible, s.IsReoccurring, s.StaffId)
	err := row.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to schedule table")
	}

	return s, nil
}

func (dep *scheduleStore) CountExistingSchedulesForStaff(ctx context.Context, staffId uint64, start, end time.Time) (int, error) {
	q := `
		SELECT COUNT(s.schedule_id) FROM schedule s
		WHERE s.staff_id = $1
		AND (
			(s.schedule_start BETWEEN $2 AND $3) OR
			(s.schedule_end BETWEEN $2 AND $3)
		)
	`

	count := -1
	row := dep.db.QueryRowContext(ctx, q, staffId, start, end)

	if err := row.Scan(&count); err != nil {
		dep.logger.Error(err)
		return 0, err
	} else if count == -1 {
		dep.logger.Error(fmt.Sprintf("exception counting existing schedules for staff id %v", staffId))
		return 0, fmt.Errorf("exception counting existing schedules for staff")
	}

	return count, nil
}

func (dep *scheduleStore) CountSchedulesInRangeAndVisibility(ctx context.Context, staffId uint64, start time.Time, end time.Time, isVisible bool) (int, error) {
	var count int

	q := `
      SELECT COUNT(*) FROM schedule s
      WHERE s.staff_id = $1
      AND (
          ($2 BETWEEN s.schedule_start AND s.schedule_end) AND
          ($3 BETWEEN s.schedule_start AND s.schedule_end)
      )
      AND is_visible = $4
    `

	row := dep.db.QueryRowContext(ctx, q, staffId, start, end, isVisible)
	if err := row.Scan(&count); err != nil {
		dep.logger.Error(err.Error())
		return 0, fmt.Errorf("exception counting schedules in range & visibility")
	}

	return count, nil
}
