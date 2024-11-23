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

type ReservationService struct {
	JunctionId    uint64 `json:"junction_id"`
	ReservationId uint64 `json:"reservation_id"`
	ServiceId     uint64 `json:"service_id"`
}
