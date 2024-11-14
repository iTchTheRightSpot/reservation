package models

import "time"

type StaffJwtObj struct {
	Roles     []RoleEnum `json:"roles"`
	StaffUUID string     `json:"staff_uuid"`
}

type JwtResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}
