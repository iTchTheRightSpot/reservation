package models

import "time"

type JwtObj struct {
	UserId         string               `json:"user_id"`
	AccessControls []RolePermissionEnum `json:"access_controls"`
	ExpireAt       *time.Time           `json:"expire_at"`
}

type JwtResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

type ActiveUser struct {
	UserId         string               `json:"user_id"`
	Firstname      string               `json:"firstname"`
	ImageKey       *string              `json:"image_key"`
	AccessControls []RolePermissionEnum `json:"access_controls"`
}

type Login struct {
	Email    string `json:"email" validate:"required,min=1,max=320"`
	Password string `json:"password" validate:"required,min=8,max=15"`
}

type AddRoleAndPermissionPayload struct {
	UserId         string               `json:"user_id" validate:"required"`
	RolePermission []RolePermissionEnum `json:"role_permission" validate:"required"`
}
