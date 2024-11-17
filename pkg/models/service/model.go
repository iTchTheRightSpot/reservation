package service

type Service struct {
	ServiceId   uint64  `json:"service_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsVisible   bool    `json:"is_visible"`
	Duration    uint8   `json:"duration"`
	CleanUpTime uint8   `json:"clean_up_time"`
}

type ServicePayload struct {
	Name        string  `json:"name" validate:"required,min=50"`
	Price       float64 `json:"price" validate:"required"`
	IsVisible   bool    `json:"is_visible"`
	Duration    uint8   `json:"duration" validate:"required"`
	CleanUpTime uint8   `json:"clean_up_time" validate:"required"`
}
