package service

type ServiceTypeEntity struct {
	ServiceId     uint64  `json:"service_id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	IsVisible     bool    `json:"is_visible"`
	IsReoccurring bool    `json:"is_reoccurring"`
	Duration      int     `json:"duration"`
	CleanUpTime   int     `json:"clean_up_time"`
}

type ServiceTypePayload struct {
	Name          string  `json:"name" validate:"required,max=50"`
	Price         float64 `json:"price" validate:"required"`
	IsVisible     bool    `json:"is_visible"`
	IsReoccurring bool    `json:"is_reoccurring"`
	Duration      int     `json:"duration" validate:"required"`
	CleanUpTime   int     `json:"clean_up_time" validate:"required"`
}

type ServiceTypeResponse struct {
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	IsVisible     bool    `json:"is_visible"`
	IsReoccurring bool    `json:"is_reoccurring"`
	Duration      int     `json:"duration"`
	CleanUpTime   int     `json:"clean_up_time"`
}
