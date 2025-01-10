package reservation

import (
	"context"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"strings"
	"time"
)

type IReservationStore interface {
	CountReservationsInRange(ctx context.Context, staffId uint64, start, end time.Time, statuses ...reservation.ReservationEnum) (int, error)
	Save(ctx context.Context, r *reservation.Reservation) error
	ReservationById(ctx context.Context, reservationId uint64) (*reservation.Reservation, error)
	UpdateReservationStatus(ctx context.Context, reservationId uint64, status reservation.ReservationEnum) error
}

type reservationStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewReservationStore(l utils.ILogger, db pkg.Db) IReservationStore {
	return &reservationStore{logger: l, db: db}
}

func (dep *reservationStore) UpdateReservationStatus(ctx context.Context, reservationId uint64, status reservation.ReservationEnum) error {
	q := "UPDATE reservation SET status = $2 WHERE reservation_id = $1"

	_, err := dep.db.ExecContext(ctx, q, reservationId, status)
	if err != nil {
		dep.logger.Error(err.Error())
		return errors.New("error updating reservation status")
	}

	return nil
}

func (dep *reservationStore) ReservationById(ctx context.Context, reservationId uint64) (*reservation.Reservation, error) {
	var r reservation.Reservation

	q := "SELECT * FROM reservation WHERE reservation_id = $1"
	row := dep.db.QueryRowContext(ctx, q, reservationId)

	err := row.Scan(
		&r.ReservationId, &r.Name, &r.Email, &r.Description, &r.Phone, &r.Price, &r.Status, &r.CreatedAt, &r.ScheduledFor, &r.ExpireAt, &r.StaffId)

	if err != nil {
		dep.logger.Error(err.Error())
		return nil, fmt.Errorf("error retrieving reservation with id %v", reservationId)
	}

	return &r, nil
}

func (dep *reservationStore) CountReservationsInRange(ctx context.Context, staffId uint64, start, end time.Time, statuses ...reservation.ReservationEnum) (int, error) {
	arr := make([]string, len(statuses))
	for i, status := range statuses {
		arr[i] = fmt.Sprintf("'%s'", string(status))
	}

	q := fmt.Sprintf(`
		SELECT COUNT(*) FROM reservation
		WHERE staff_id = $1 AND (
		    (scheduled_for BETWEEN $2 AND $3) OR
		    (expire_at BETWEEN $2 AND $3)
		)
		AND status IN (%s)
	`, fmt.Sprintf("%s", strings.Join(arr, ", ")))

	var count int

	row := dep.db.QueryRowContext(ctx, q, staffId, start, end)
	if err := row.Scan(&count); err != nil {
		dep.logger.Error(err.Error())
		return 0, fmt.Errorf("error counting reservations in range")
	}

	return count, nil
}

func (dep *reservationStore) Save(ctx context.Context, r *reservation.Reservation) error {
	if r == nil {
		return errors.New("reservation cannot be nil")
	}

	q := `
        INSERT INTO reservation (staff_id, name, email, description, phone, price, status, created_at, scheduled_for, expire_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING reservation_id, staff_id, name, email, description, phone, price, status, created_at, scheduled_for, expire_at;
    `

	row := dep.db.QueryRowContext(
		ctx, q, r.StaffId, r.Name, r.Email, r.Description, r.Phone, r.Price, r.Status, r.CreatedAt, r.ScheduledFor, r.ExpireAt)

	err := row.Scan(
		&r.ReservationId, &r.StaffId, &r.Name, &r.Email, &r.Description, &r.Phone, &r.Price, &r.Status, &r.CreatedAt, &r.ScheduledFor, &r.ExpireAt)

	if err != nil {
		dep.logger.Error(err.Error())
		return errors.New("error saving reservation")
	}

	dep.logger.Error("reservation saved")
	return nil
}
