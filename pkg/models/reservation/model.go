package reservation

import "time"

type ReservationEnum string

const (
	CONFIRMED ReservationEnum = "CONFIRMED"
	CANCELLED ReservationEnum = "CANCELLED"
)

type Reservation struct {
	ReservationId uint64          `json:"reservation_id"`
	StaffId       uint64          `json:"staff_id"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Description   *string         `json:"description"`
	Address       *string         `json:"address"`
	Phone         *string         `json:"phone"`
	ImageKey      *string         `json:"image_key"`
	Price         float64         `json:"price"`
	Status        ReservationEnum `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	ScheduledFor  time.Time       `json:"scheduled_for"`
	ExpireAt      time.Time       `json:"expire_at"`
}

type ReservationServiceEntity struct {
	JunctionId    uint64 `json:"junction_id"`
	ReservationId uint64 `json:"reservation_id"`
	ServiceId     uint64 `json:"service_id"`
}

type ReservationPayload struct {
	StaffId     string    `json:"staff_id" validate:"required,min=36,max=37"`
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email" validate:"required,max=320"`
	Description string    `json:"description" validate:"max=255"`
	Address     string    `json:"address" validate:"required,max=255"`
	Phone       string    `json:"phone" validate:"required,max=20"`
	Services    []*string `json:"services" validate:"required,min=1,dive,required"`
	Timezone    string    `json:"timezone"`
	Time        string    `json:"time" validate:"required"`
}

type ReservationTimeSlots struct {
	Date  string   `json:"date"`
	Times []string `json:"times"`
}
