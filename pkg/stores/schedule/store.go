package schedule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"time"
)

type IScheduleStore interface {
	Save(ctx context.Context, s *schedule.ScheduleEntity) error
	Update(ctx context.Context, s *schedule.ScheduleEntity) (int64, error)
	CountExistingSchedulesForStaff(ctx context.Context, staffId uint64, start, end time.Time) (int, error)
	SchedulesInRange(ctx context.Context, staffId uint64, start time.Time, end time.Time) ([]*schedule.ScheduleEntity, error)
	CountSchedulesInRangeAndVisibility(ctx context.Context, staffId uint64, start time.Time, end time.Time, isVisible bool) (int, error)
	SchedulesWithinTimeframe(ctx context.Context, staffId uint64, start time.Time, end time.Time, isVisible bool, duration int) ([]*schedule.ScheduleEntity, error)
	ScheduleByScheduleId(ctx context.Context, scheduleId uint64) (*schedule.ScheduleEntity, error)
	Delete(context.Context, uint64) (int64, error)
}

type scheduleStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewScheduleStore(l utils.ILogger, db pkg.Db) IScheduleStore {
	return &scheduleStore{logger: l, db: db}
}

func (dep *scheduleStore) SchedulesInRange(ctx context.Context, staffId uint64, start time.Time, end time.Time) ([]*schedule.ScheduleEntity, error) {
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
		return nil, errors.New("error retrieve schedules in range")
	}

	defer func(rs *sql.Rows) { err = rs.Close() }(rows)

	var arr []*schedule.ScheduleEntity
	for rows.Next() {
		var s schedule.ScheduleEntity

		err = rows.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)
		if err != nil {
			dep.logger.Error(err.Error())
			return nil, errors.New("error scanning schedules")
		}

		arr = append(arr, &s)
	}

	return arr, err
}

func (dep *scheduleStore) Save(ctx context.Context, s *schedule.ScheduleEntity) error {
	if s == nil {
		return errors.New("schedule object is nil")
	}

	q := `
		INSERT INTO schedule (schedule_start, schedule_end, is_visible, is_reoccurring, staff_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING schedule_id, schedule_start, schedule_end, is_visible, is_reoccurring, staff_id
	`

	row := dep.db.QueryRowContext(ctx, q, s.Start, s.End, s.IsVisible, s.IsReoccurring, s.StaffId)
	err := row.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)

	if err != nil {
		dep.logger.Error(err.Error())
		return errors.New("exception saving to schedule table")
	}

	return nil
}

func (dep *scheduleStore) Update(ctx context.Context, s *schedule.ScheduleEntity) (int64, error) {
	if s == nil {
		return 0, errors.New("schedule object is nil")
	}

	q := "UPDATE schedule SET is_visible = $2, is_reoccurring = $3 WHERE schedule_id = $1"

	res, err := dep.db.ExecContext(ctx, q, s.ScheduleId, s.IsVisible, s.IsReoccurring)
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, errors.New("exception saving to schedule table")
	}

	return res.RowsAffected()
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

func (dep *scheduleStore) SchedulesWithinTimeframe(ctx context.Context, staffId uint64, start time.Time, end time.Time, isVisible bool, duration int) ([]*schedule.ScheduleEntity, error) {
	q := `
	   SELECT * FROM schedule s
	   WHERE s.staff_id = $1
	   AND (
	       (s.schedule_start BETWEEN $2 AND $3) OR
	       (s.schedule_end BETWEEN $2 AND $3)
	   )
	   AND is_visible = $4
	   AND EXTRACT(EPOCH FROM (s.schedule_end - s.schedule_start)) >= $5
	`

	rows, err := dep.db.QueryContext(ctx, q, staffId, start, end, isVisible, duration)
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("error retrieving available schedules in range")
	}

	defer func(rs *sql.Rows) { err = rs.Close() }(rows)

	var arr []*schedule.ScheduleEntity
	for rows.Next() {
		var s schedule.ScheduleEntity

		err = rows.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)
		if err != nil {
			dep.logger.Error(err.Error())
			return nil, errors.New("error scanning schedules")
		}

		arr = append(arr, &s)
	}

	return arr, err
}

func (dep *scheduleStore) ScheduleByScheduleId(ctx context.Context, scheduleId uint64) (*schedule.ScheduleEntity, error) {
	var s schedule.ScheduleEntity

	row := dep.db.QueryRowContext(ctx, "SELECT * FROM schedule s WHERE s.schedule_id = $1 LIMIT 1", scheduleId)
	err := row.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)

	if err != nil {
		dep.logger.Error(err.Error())
		return nil, fmt.Errorf("error retrieving schedule with id %v", scheduleId)
	}

	return &s, nil
}

func (dep *scheduleStore) Delete(ctx context.Context, scheduleId uint64) (int64, error) {
	result, err := dep.db.ExecContext(ctx, "DELETE FROM schedule WHERE schedule_id = $1", scheduleId)
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, fmt.Errorf("error schedule schedule with id %v", scheduleId)
	}

	return result.RowsAffected()
}
