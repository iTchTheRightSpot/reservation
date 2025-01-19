package profile

type Profile struct {
	ProfileId uint64  `json:"profile_id"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Email     string  `json:"email"`
	Password  []byte  `json:"password"`
	ImageKey  *string `json:"image_key"`
}

type ProfilePayload struct {
	Firstname string `json:"firstname" validate:"required,min=1,max=50"`
	Lastname  string `json:"lastname" validate:"required,min=1,max=50"`
	Email     string `json:"email" validate:"required,min=1,max=320"`
	Password  string `json:"password" validate:"required,min=8,max=15"`
}
