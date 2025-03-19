package schedule

import (
	"context"
	"database/sql"
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
		return nil, &utils.NotFoundError{Message: "error retrieve schedules in range"}
	}

	defer func(rs *sql.Rows) {
		if err = rs.Close(); err != nil {
			dep.logger.Error(err.Error())
			err = &utils.ServerError{Message: "error closing stream after schedules in range"}
		}
	}(rows)

	var arr []*schedule.ScheduleEntity
	for rows.Next() {
		var s schedule.ScheduleEntity

		err = rows.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)
		if err != nil {
			dep.logger.Error(err.Error())
			return nil, &utils.ServerError{Message: "error scanning schedules in range"}
		}

		arr = append(arr, &s)
	}

	return arr, err
}

func (dep *scheduleStore) Save(ctx context.Context, s *schedule.ScheduleEntity) error {
	if s == nil {
		return &utils.ServerError{Message: "schedule object is nil"}
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
		return &utils.InsertionError{Message: "error saving schedule"}
	}

	return nil
}

func (dep *scheduleStore) Update(ctx context.Context, s *schedule.ScheduleEntity) (int64, error) {
	if s == nil {
		return 0, &utils.ServerError{Message: "schedule object is nil"}
	}

	q := "UPDATE schedule SET is_visible = $2, is_reoccurring = $3 WHERE schedule_id = $1"

	res, err := dep.db.ExecContext(ctx, q, s.ScheduleId, s.IsVisible, s.IsReoccurring)
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, &utils.InsertionError{Message: "error updating schedule"}
	}

	i, err := res.RowsAffected()
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, &utils.InsertionError{Message: "error updating schedule no."}
	}

	return i, nil
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
		dep.logger.Error(err.Error())
		return 0, &utils.NotFoundError{Message: "invalid staff id"}
	} else if count == -1 {
		return 0, &utils.NotFoundError{Message: "error counting existing schedules for staff"}
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
		return 0, &utils.NotFoundError{Message: "error counting schedules in range & visibility"}
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
		return nil, &utils.NotFoundError{Message: "error retrieving available schedules in range"}
	}

	defer func(rs *sql.Rows) {
		if err = rs.Close(); err != nil {
			dep.logger.Error(err.Error())
			err = &utils.ServerError{Message: "error closing stream after schedules in range"}
		}
	}(rows)

	var arr []*schedule.ScheduleEntity
	for rows.Next() {
		var s schedule.ScheduleEntity

		err = rows.Scan(&s.ScheduleId, &s.Start, &s.End, &s.IsVisible, &s.IsReoccurring, &s.StaffId)
		if err != nil {
			dep.logger.Error(err.Error())
			return nil, &utils.ServerError{Message: "error scanning schedules in range"}
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
		return nil, &utils.NotFoundError{Message: "error retrieving schedules by id"}
	}

	return &s, nil
}

func (dep *scheduleStore) Delete(ctx context.Context, scheduleId uint64) (int64, error) {
	result, err := dep.db.ExecContext(ctx, "DELETE FROM schedule WHERE schedule_id = $1", scheduleId)
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, &utils.InsertionError{Message: "error deleting schedule"}
	}

	i, err := result.RowsAffected()
	if err != nil {
		dep.logger.Error(err.Error())
		return 0, &utils.InsertionError{Message: "error deleting schedule. row"}
	}
	return i, nil
}
