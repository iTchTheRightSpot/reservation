package shift

import (
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"slices"
	"time"
)

type Shift struct {
	ShiftId       uint64    `json:"shift_id"`
	StaffId       uint64    `json:"staff_id"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"shift_end"`
	IsEnabled     bool      `json:"is_enabled"`
	IsReoccurring bool      `json:"is_reoccurring"`
}

type ShiftDto struct {
	StaffUUID string          `json:"staff_uuid" validate:"required"`
	Times     *[]ShiftSegment `json:"shift_segments" validate:"required,dive,required"`
}

type ShiftSegment struct {
	IsVisible     bool   `json:"is_visible" validate:"required"`
	IsReoccurring bool   `json:"is_reoccurring" validate:"required"`
	Start         string `json:"start" validate:"required"` // in ISO 8601 standard
	Duration      int    `json:"duration" validate:"required"`
}

type ScheduledPeriod struct {
	IsVisible     bool
	IsReoccurring bool
	Start         time.Time
	End           time.Time
}

func (dto *ShiftDto) CheckForOverlappingSegments(timezone *time.Location) ([]ScheduledPeriod, error) {
	if dto.Times == nil {
		return nil, fmt.Errorf("time_slots cannot be nil")
	}

	arr := make([]ScheduledPeriod, len(*dto.Times))

	for idx, slot := range *dto.Times {
		parse, err := time.Parse(utils.TimeFormat, slot.Start)
		if err != nil {
			return nil, err
		}

		start := parse.In(timezone)
		end := start.Add(time.Duration(slot.Duration) * time.Second)

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
