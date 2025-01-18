package staff

import "github.com/google/uuid"

type Staff struct {
	StaffId   uint64    `json:"staff_id"`
	UUID      uuid.UUID `json:"uuid"`
	Bio       *string   `json:"bio"`
	ProfileId *uint64   `json:"profile_id"`
}

type StaffServiceEntity struct {
	JunctionId uint64 `json:"junction_id"`
	StaffId    uint64 `json:"staff_id"`
	ServiceId  uint64 `json:"service_id"`
}

type StaffStoreFrontDb struct {
	Name     string    `json:"name"`
	UUID     uuid.UUID `json:"uuid"`
	Bio      *string   `json:"bio"`
	ImageKey *string   `json:"image_key"`
}
