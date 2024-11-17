package staff

import "github.com/google/uuid"

type Staff struct {
	StaffId   uint64    `json:"staff_id"`
	UUID      uuid.UUID `json:"uuid"`
	Bio       *string   `json:"bio"`
	ProfileId *uint64   `json:"profile_id"`
}
