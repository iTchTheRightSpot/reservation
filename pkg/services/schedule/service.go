package schedule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"sync"
	"time"
)

type IScheduleService interface {
	Schedules(ctx context.Context, payload *schedule.AllSchedulesPayload) ([]*schedule.ScheduleResponse, error)
	Create(ctx context.Context, dto *schedule.SchedulePayload) error
}

type scheduleService struct {
	logger   utils.ILogger
	adapters *stores.Adapters
}

func NewScheduleService(l utils.ILogger, a *stores.Adapters) IScheduleService {
	return &scheduleService{logger: l, adapters: a}
}

func (dep *scheduleService) Schedules(ctx context.Context, p *schedule.AllSchedulesPayload) ([]*schedule.ScheduleResponse, error) {
	s, err := dep.adapters.StaffStore.StaffByUUID(ctx, p.StaffUUID)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "invalid staff id"}
	}

	start := time.Date(p.Year, time.Month(p.Month), 1, 0, 0, 0, 0, dep.logger.Timezone())
	schedules, err := dep.adapters.ScheduleStore.SchedulesInRange(ctx, s.StaffId, start, start.AddDate(0, 1, -1))
	if err != nil {
		return nil, &utils.BadRequestError{Message: err.Error()}
	}

	res := make([]*schedule.ScheduleResponse, len(schedules))
	for i, ele := range schedules {
		res[i] = &schedule.ScheduleResponse{
			ScheduleId:    ele.ScheduleId,
			IsVisible:     ele.IsVisible,
			IsReoccurring: ele.IsReoccurring,
			Start:         fmt.Sprintf("%v", ele.Start.In(p.Timezone).UnixMilli()),
			End:           fmt.Sprintf("%v", ele.End.In(p.Timezone).UnixMilli()),
		}
	}

	return res, nil
}

func (dep *scheduleService) validateSegments(ctx context.Context, staffID uint64, segments []schedule.ScheduledPeriod) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(segments))

	for _, seg := range segments {
		wg.Add(1)
		go func(segment schedule.ScheduledPeriod) {
			defer wg.Done()
			count, err := dep.adapters.ScheduleStore.CountExistingSchedulesForStaff(ctx, staffID, segment.Start, segment.End)
			if err != nil {
				errChan <- errors.New(err.Error())
				return
			}
			if count > 0 {
				errChan <- fmt.Errorf("duplicate schedule detected from %v to %v", segment.Start, segment.End)
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
		dep.logger.Error(err.Error())
		return &utils.BadRequestError{Message: err.Error()}
	}

	staff, err := dep.adapters.StaffStore.StaffByUUID(ctx, dto.StaffId)
	if err != nil {
		return &utils.BadRequestError{Message: "invalid staff id"}
	}

	if err = dep.validateSegments(ctx, staff.StaffId, segments); err != nil {
		return &utils.BadRequestError{Message: err.Error()}
	}

	// TODO look into exclusion constraint
	return dep.adapters.Transaction.RunInTransaction(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable}, func(adapters *stores.Adapters) error {
		for _, segment := range segments {
			_, err = adapters.ScheduleStore.Save(ctx, &schedule.Schedule{
				StaffId:   staff.StaffId,
				Start:     segment.Start,
				End:       segment.End,
				IsVisible: segment.IsVisible,
			})
			if err != nil {
				return &utils.InsertionError{Message: "duplicate schedule"}
			}
		}
		return nil
	})
}
