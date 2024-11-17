package models

import "time"

type JwtObj struct {
	UserId         string           `json:"user_id"`
	AccessControls []RolePermission `json:"access_controls"`
	ExpireAt       *time.Time       `json:"expire_at"`
}

type JwtResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}
