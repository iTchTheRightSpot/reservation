package models

type RoleEnum string

const (
	STAFF     RoleEnum = "STAFF"
	DEVELOPER RoleEnum = "DEVELOPER"
	USER      RoleEnum = "USER"
)

type Role struct {
	RoleId    uint64   `json:"role_id"`
	Role      RoleEnum `json:"role"`
	ProfileId uint64   `json:"profile_id"`
}
