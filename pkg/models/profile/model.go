package profile

type Profile struct {
	ProfileId uint64  `json:"profile_id"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Email     string  `json:"email"`
	ImageKey  *string `json:"image_key"`
}
