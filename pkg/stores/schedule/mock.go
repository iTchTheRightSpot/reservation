package schedule

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"time"
)

type MockScheduleStore struct {
	ScheduleSave                             *schedule.Schedule
	ScheduleSaveError                        error
	ScheduleSaveCalled                       bool
	SchedulesInRangeReturn                   []*schedule.Schedule
	SchedulesInRangeError                    error
	SchedulesInRangeCalled                   bool
	CountExistingSchedulesForStaffReturn     int
	CountExistingSchedulesForStaffError      error
	CountExistingSchedulesForStaffCalled     bool
	CountSchedulesInRangeAndVisibilityReturn int
	CountSchedulesInRangeAndVisibilityError  error
	CountSchedulesInRangeAndVisibilityCalled bool
}

func (dep *MockScheduleStore) Save(context.Context, *schedule.Schedule) (*schedule.Schedule, error) {
	dep.ScheduleSaveCalled = true
	return dep.ScheduleSave, dep.ScheduleSaveError
}

func (dep *MockScheduleStore) SchedulesInRange(context.Context, uint64, time.Time, time.Time) ([]*schedule.Schedule, error) {
	dep.SchedulesInRangeCalled = true
	return dep.SchedulesInRangeReturn, dep.SchedulesInRangeError
}

func (dep *MockScheduleStore) CountExistingSchedulesForStaff(context.Context, uint64, time.Time, time.Time) (int, error) {
	dep.CountExistingSchedulesForStaffCalled = true
	return dep.CountExistingSchedulesForStaffReturn, dep.CountExistingSchedulesForStaffError
}

func (dep *MockScheduleStore) CountSchedulesInRangeAndVisibility(context.Context, uint64, *time.Time, time.Time, bool) (int, error) {
	dep.CountSchedulesInRangeAndVisibilityCalled = true
	return dep.CountSchedulesInRangeAndVisibilityReturn, dep.CountSchedulesInRangeAndVisibilityError
}
