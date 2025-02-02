package schedule

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"time"
)

type IScheduleService interface {
	Schedules(ctx context.Context, payload *schedule.AllSchedulesPayload) ([]*schedule.ScheduleResponse, error)
	Create(ctx context.Context, dto *schedule.SchedulePayload) error
	Update(ctx context.Context, p *schedule.UpdateSchedulePayload) error
	Delete(ctx context.Context, scheduleId uint64) error
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

	from := time.Date(p.Year, time.Month(p.Month), 1, 0, 0, 0, 0, dep.logger.Timezone())
	to := time.Date(from.Year(), from.Month()+1, 0, 23, 59, 59, 999999999, from.Location())
	schedules, err := dep.adapters.ScheduleStore.SchedulesInRange(ctx, s.StaffId, from, to)
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

	return dep.adapters.Transaction.RunInTransaction(func(adapters *stores.Adapters) error {
		for _, segment := range segments {
			err = adapters.ScheduleStore.Save(ctx, &schedule.ScheduleEntity{
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

func (dep *scheduleService) Update(ctx context.Context, p *schedule.UpdateSchedulePayload) error {
	o, err := dep.adapters.ScheduleStore.ScheduleByScheduleId(ctx, p.ScheduleId)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.NotFoundError{Message: "invalid schedule_id"}
	}

	if o.IsVisible == p.IsVisible && o.IsReoccurring == p.IsReoccurring {
		return nil
	}

	o.IsVisible = p.IsVisible
	o.IsReoccurring = p.IsReoccurring

	_, err = dep.adapters.ScheduleStore.Update(ctx, o)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.InsertionError{Message: "error updating schedule"}
	}

	return nil
}

func (dep *scheduleService) Delete(ctx context.Context, scheduleId uint64) error {
	o, err := dep.adapters.ScheduleStore.ScheduleByScheduleId(ctx, scheduleId)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.NotFoundError{Message: "invalid schedule_id"}
	}

	count, err := dep.adapters.ReservationStore.CountReservationsInRange(ctx, o.StaffId, o.Start, o.End, reservation.CONFIRMED, reservation.CANCELLED)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.NotFoundError{Message: "validation checks failed before attempting deletion"}
	}

	if count > 0 {
		m := "deletion failed. schedule references reservations"
		dep.logger.Error(m)
		return &utils.InsertionError{Message: m}
	}

	_, err = dep.adapters.ScheduleStore.Delete(ctx, o.ScheduleId)
	if err != nil {
		dep.logger.Error(err.Error())
		return &utils.InsertionError{Message: "error deleting schedule"}
	}

	return nil
}
