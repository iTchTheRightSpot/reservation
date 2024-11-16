package models

import "time"

type JwtObj struct {
	UserUUID       string           `json:"user_uuid"`
	AccessControls []RolePermission `json:"access_controls"`
	ExpireAt       *time.Time       `json:"expire_at"`
}

type JwtResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}
