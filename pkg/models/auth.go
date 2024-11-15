package models

import "time"

type JwtObj struct {
	Roles    []RoleEnum `json:"roles"`
	UserUUID string     `json:"user_uuid"`
	ExpireAt time.Time  `json:"expire_at"`
}

type JwtResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}
