package schedule

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"time"
)

type MockScheduleStore struct {
	ScheduleSave                                     *schedule.ScheduleEntity
	ScheduleSaveError                                error
	ScheduleSaveCalled                               bool
	SchedulesInRangeReturn                           []*schedule.ScheduleEntity
	SchedulesInRangeError                            error
	SchedulesInRangeCalled                           bool
	CountExistingSchedulesForStaffReturn             int
	CountExistingSchedulesForStaffError              error
	CountExistingSchedulesForStaffCalled             bool
	CountSchedulesInRangeAndVisibilityReturn         int
	CountSchedulesInRangeAndVisibilityError          error
	CountSchedulesInRangeAndVisibilityCalled         bool
	SchedulesInRangeAndVisibilityAndDifferenceReturn []*schedule.ScheduleEntity
	SchedulesInRangeAndVisibilityAndDifferenceError  error
	SchedulesInRangeAndVisibilityAndDifferenceCalled bool
}

func (dep *MockScheduleStore) Save(context.Context, *schedule.ScheduleEntity) error {
	dep.ScheduleSaveCalled = true
	return dep.ScheduleSaveError
}

func (dep *MockScheduleStore) Update(context.Context, *schedule.ScheduleEntity) (int64, error) {
	return 0, nil
}

func (dep *MockScheduleStore) Delete(context.Context, uint64) (int64, error) {
	return 0, nil
}

func (dep *MockScheduleStore) SchedulesInRange(context.Context, uint64, time.Time, time.Time) ([]*schedule.ScheduleEntity, error) {
	dep.SchedulesInRangeCalled = true
	return dep.SchedulesInRangeReturn, dep.SchedulesInRangeError
}

func (dep *MockScheduleStore) CountExistingSchedulesForStaff(context.Context, uint64, time.Time, time.Time) (int, error) {
	dep.CountExistingSchedulesForStaffCalled = true
	return dep.CountExistingSchedulesForStaffReturn, dep.CountExistingSchedulesForStaffError
}

func (dep *MockScheduleStore) CountSchedulesInRangeAndVisibility(context.Context, uint64, time.Time, time.Time, bool) (int, error) {
	dep.CountSchedulesInRangeAndVisibilityCalled = true
	return dep.CountSchedulesInRangeAndVisibilityReturn, dep.CountSchedulesInRangeAndVisibilityError
}

func (dep *MockScheduleStore) SchedulesWithinTimeframe(context.Context, uint64, time.Time, time.Time, bool, int) ([]*schedule.ScheduleEntity, error) {
	dep.SchedulesInRangeAndVisibilityAndDifferenceCalled = true
	return dep.SchedulesInRangeAndVisibilityAndDifferenceReturn, dep.SchedulesInRangeAndVisibilityAndDifferenceError
}

func (dep *MockScheduleStore) ScheduleByScheduleId(context.Context, uint64) (*schedule.ScheduleEntity, error) {
	return nil, nil
}
