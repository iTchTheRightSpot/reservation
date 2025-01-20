package models

type ProfileEntity struct {
	ProfileId uint64  `json:"profile_id"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Locked    bool    `json:"locked"`
	ImageKey  *string `json:"image_key"`
}

type ProfilePayload struct {
	Firstname string `json:"firstname" validate:"required,min=1,max=50"`
	Lastname  string `json:"lastname" validate:"required,min=1,max=50"`
	Email     string `json:"email" validate:"required,min=1,max=320"`
	Password  string `json:"password" validate:"required,min=8,max=15"`
}

type RoleEnum string

const (
	STAFF     RoleEnum = "STAFF"
	DEVELOPER RoleEnum = "DEVELOPER"
	USER      RoleEnum = "USER"
)

type RoleEntity struct {
	RoleId    uint64   `json:"role_id"`
	Role      RoleEnum `json:"role"`
	ProfileId uint64   `json:"profile_id"`
}

type PermissionEnum string

const (
	READ   PermissionEnum = "READ"
	WRITE  PermissionEnum = "WRITE"
	DELETE PermissionEnum = "DELETE"
)

type PermissionEntity struct {
	PermissionId uint64         `json:"permission_id"`
	Permission   PermissionEnum `json:"permission"`
	RoleId       uint64         `json:"role_id"`
}

type RolePermissionEnum struct {
	Role        RoleEnum         `json:"role"`
	Permissions []PermissionEnum `json:"permissions"`
}

type RolePermissionEntity struct {
	Role        RoleEntity         `json:"role"`
	Permissions []PermissionEntity `json:"permissions"`
}

type ProfileRolePermissionEntity struct {
	Profile        ProfileEntity          `json:"profile"`
	RolePermission []RolePermissionEntity `json:"role_perm"`
}
