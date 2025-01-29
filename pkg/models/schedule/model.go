package schedule

import (
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"slices"
	"time"
)

type ScheduleEntity struct {
	ScheduleId    uint64    `json:"schedule_id"`
	StaffId       uint64    `json:"staff_id"`
	Start         time.Time `json:"schedule_start"`
	End           time.Time `json:"schedule_end"`
	IsVisible     bool      `json:"is_visible"`
	IsReoccurring bool      `json:"is_reoccurring"`
}

type UpdateSchedulePayload struct {
	ScheduleId    uint64 `json:"schedule_id" validate:"required"`
	IsVisible     bool   `json:"is_visible"`
	IsReoccurring bool   `json:"is_reoccurring"`
}

type ScheduleResponse struct {
	ScheduleId    uint64 `json:"schedule_id"`
	IsVisible     bool   `json:"is_visible"`
	IsReoccurring bool   `json:"is_reoccurring"`
	Start         string `json:"start"`
	End           string `json:"end"`
}

type AllSchedulesPayload struct {
	StaffUUID string         `json:"staff_uuid" validate:"required,min=36,max=37"`
	Month     int            `validate:"required,min=1,max=12"`
	Year      int            `validate:"required"`
	Timezone  *time.Location `validate:"required"`
}

type SchedulePayload struct {
	StaffId string                    `json:"staff_id" validate:"required"`
	Times   *[]ScheduleSegmentPayload `json:"schedule_segments" validate:"required,min=1,dive,required"`
}

type ScheduleSegmentPayload struct {
	IsVisible     bool   `json:"is_visible"`
	IsReoccurring bool   `json:"is_reoccurring"`
	Start         string `json:"start" validate:"required"` // in ISO 8601 standard
	Duration      int    `json:"duration" validate:"required"`
}

type ScheduledPeriod struct {
	IsVisible     bool
	IsReoccurring bool
	Start         time.Time
	End           time.Time
}

func (dto *SchedulePayload) CheckForOverlappingSegments(now time.Time, timezone *time.Location) ([]ScheduledPeriod, error) {
	if dto.Times == nil {
		return nil, errors.New("time_slots cannot be nil")
	}

	arr := make([]ScheduledPeriod, len(*dto.Times))

	for idx, slot := range *dto.Times {
		parse, err := time.Parse(utils.TimeFormat, slot.Start)
		if err != nil {
			return nil, err
		}

		start := parse.In(timezone)

		// validate no date in the past
		if now.After(start) {
			times := *dto.Times
			return nil, fmt.Errorf("%s is in the past ", times[idx].Start)
		}

		end := start.Add(time.Duration(slot.Duration) * time.Second)

		// validate start and end are within the same dat
		if start.Day() != end.Day() || start.Month() != end.Month() || start.Year() != end.Year() {
			t := *dto.Times
			return nil, fmt.Errorf("%s plus duration cannot include the next day", t[idx].Start)
		}

		// validate no conflicts
		conflict := slices.ContainsFunc(arr, func(obj ScheduledPeriod) bool {
			return start.Before(obj.End) && end.After(obj.Start)
		})

		if conflict {
			return nil, fmt.Errorf("conflicting date %s & duration %v", slot.Start, slot.Duration)
		}

		arr[idx] = ScheduledPeriod{
			IsReoccurring: slot.IsReoccurring,
			IsVisible:     slot.IsVisible,
			Start:         start,
			End:           end,
		}
	}

	return arr, nil
}
