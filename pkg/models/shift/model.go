package shift

import "time"

type Shift struct {
	ShiftId   uint64    `json:"shift_id"`
	StaffId   uint64    `json:"staff_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"shift_end"`
	IsEnabled bool      `json:"is_enabled"`
}
