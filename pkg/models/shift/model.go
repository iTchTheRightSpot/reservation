package shift

import (
	"time"
)

type Shift struct {
	ShiftId   uint64    `json:"shift_id"`
	StaffId   uint64    `json:"staff_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"shift_end"`
	IsEnabled bool      `json:"is_enabled"`
}

type TimeSlot struct {
	IsVisible bool   `json:"is_visible" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	StartTime string `json:"start_time" validate:"required"`
	Duration  int    `json:"duration" validate:"required"`
}

type ShiftDto struct {
	StaffUUID string      `json:"staff_uuid" validate:"required"`
	TimeSlots *[]TimeSlot `json:"time_slots" validate:"required,dive,required"`
}
