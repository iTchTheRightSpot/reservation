package schedule

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"sync"
)

type IScheduleService interface {
	Create(ctx context.Context, dto *schedule.SchedulePayload) error
}

type scheduleService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewScheduleService(l utils.ILogger, a *stores.Adapters) IScheduleService {
	return &scheduleService{logger: l, adapters: a}
}

func (dep *scheduleService) validateSegments(ctx context.Context, staffID uint64, segments []schedule.ScheduledPeriod) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(segments))

	for _, seg := range segments {
		wg.Add(1)
		go func(segment schedule.ScheduledPeriod) {
			defer wg.Done()
			count, err := dep.adapters.ShiftStore.CountExistingShiftsForStaff(ctx, staffID, segment.Start, segment.End)
			if err != nil {
				errChan <- fmt.Errorf("error checking existing shifts: %w", err)
				return
			}
			if count > 0 {
				errChan <- fmt.Errorf("duplicate shift detected from %v to %v", segment.Start, segment.End)
			}
		}(seg)
	}

	wg.Wait()
	close(errChan)

	return <-errChan
}

func (dep *scheduleService) Create(ctx context.Context, dto *schedule.SchedulePayload) error {
	segments, err := dto.CheckForOverlappingSegments(dep.logger.Date(), dep.logger.Timezone())
	if err != nil {
		dep.logger.Error(err)
		return err
	}

	staff, err := dep.adapters.StaffStore.StaffByUUID(ctx, dto.StaffId)
	if err != nil {
		return err
	}

	if err = dep.validateSegments(ctx, staff.StaffId, segments); err != nil {
		return err
	}

	return dep.adapters.Transaction.RunInTransaction(func(adapters *stores.Adapters) error {
		for _, segment := range segments {
			_, err = adapters.ShiftStore.Save(ctx, &schedule.Schedule{
				StaffId:   staff.StaffId,
				Start:     segment.Start,
				End:       segment.End,
				IsVisible: segment.IsVisible,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
