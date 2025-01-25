package staff

import (
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
)

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

type AllStaffsEntity struct {
	Firstname      string                       `json:"firstname"`
	Lastname       string                       `json:"lastname"`
	Email          string                       `json:"email"`
	Locked         bool                         `json:"locked"`
	ImageKey       *string                      `json:"image_key"`
	UUID           uuid.UUID                    `json:"user_id"`
	Bio            *string                      `json:"bio"`
	AccessControls []*models.RolePermissionEnum `json:"access_controls"`
}
