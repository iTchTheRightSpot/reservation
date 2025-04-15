package reservation

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/utility/utils"
	"strings"
	"time"
)

type IReservationStore interface {
	CountReservationsInRange(ctx context.Context, staffId uint64, start, end time.Time, statuses ...reservation.ReservationEnum) (int, error)
	Save(ctx context.Context, r *reservation.Reservation) error
	ReservationById(ctx context.Context, reservationId uint64) (*reservation.Reservation, error)
	UpdateReservationStatus(ctx context.Context, reservationId uint64, status reservation.ReservationEnum) (int64, error)
	BookingsInRange(ctx context.Context, staffId uint64, from time.Time, to time.Time) ([]*reservation.CRMBookingsResponse, error)
}

type reservationStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewReservationStore(l utils.ILogger, db pkg.Db) IReservationStore {
	return &reservationStore{logger: l, db: db}
}

func (dep *reservationStore) BookingsInRange(ctx context.Context, staffId uint64, from time.Time, to time.Time) ([]*reservation.CRMBookingsResponse, error) {
	q := `
		SELECT
		    p.firstname as staff_name,
		    p.image_key,
			r.reservation_id,
			r.name,
			r.email,
			r.description,
			r.phone,
			r.price,
			r.status,
			r.scheduled_for,
			r.expire_at,
			json_agg(s.name) as services
		FROM reservation r
		INNER JOIN reservation_service rs ON rs.reservation_id = r.reservation_id
		INNER JOIN service_type s ON s.service_id = rs.service_id
		INNER JOIN staff st ON st.staff_id = r.staff_id
		INNER JOIN profile p ON p.profile_id = st.profile_id
		WHERE r.staff_id = $1 AND (r.scheduled_for BETWEEN $2 AND $3)
		GROUP BY p.firstname, p.image_key, r.reservation_id, r.name, r.email, r.description, r.phone, r.price, r.status, r.scheduled_for, r.expire_at
	`

	rows, err := dep.db.QueryContext(ctx, q, staffId, from, to)
	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return nil, &utils.NotFoundError{Message: "error retrieving bookings in range"}
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			dep.logger.Error(ctx, err.Error())
			err = &utils.ServerError{Message: "error closing stream after bookings in range"}
		}
	}(rows)

	var arr []*reservation.CRMBookingsResponse

	for rows.Next() {
		var o reservation.CRMBookingsResponse
		var data json.RawMessage

		if err = rows.Scan(
			&o.StaffFirstname,
			&o.ImageKey,
			&o.ReservationId,
			&o.Name,
			&o.Email,
			&o.Description,
			&o.Phone,
			&o.Price,
			&o.Status,
			&o.ScheduledFor,
			&o.ExpireAt,
			&data,
		); err != nil {
			dep.logger.Error(ctx, err.Error())
			return nil, &utils.ServerError{Message: "error scanning database rows for bookings in range"}
		}

		if err = json.Unmarshal(data, &o.Services); err != nil {
			dep.logger.Error(ctx, err.Error())
			return nil, &utils.ServerError{Message: "error unmarshalling services from bookings in range"}
		}

		arr = append(arr, &o)
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(ctx, err.Error())
		return nil, &utils.NotFoundError{Message: "error iterating through bookings rows for bookings in range"}
	}

	return arr, err
}

func (dep *reservationStore) UpdateReservationStatus(ctx context.Context, reservationId uint64, status reservation.ReservationEnum) (int64, error) {
	q := "UPDATE reservation SET status = $2 WHERE reservation_id = $1"

	r, err := dep.db.ExecContext(ctx, q, reservationId, status)
	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return 0, &utils.InsertionError{Message: "error updating reservation status"}
	}

	num, err := r.RowsAffected()
	if err != nil {
		return 0, &utils.InsertionError{Message: "error displaying rows affected when updating reservation status"}
	}

	return num, nil
}

func (dep *reservationStore) ReservationById(ctx context.Context, reservationId uint64) (*reservation.Reservation, error) {
	var r reservation.Reservation

	q := "SELECT * FROM reservation WHERE reservation_id = $1"
	row := dep.db.QueryRowContext(ctx, q, reservationId)

	err := row.Scan(
		&r.ReservationId, &r.Name, &r.Email, &r.Description, &r.Phone, &r.Price, &r.Status, &r.CreatedAt, &r.ScheduledFor, &r.ExpireAt, &r.StaffId)

	if err != nil {
		dep.logger.Error(ctx, err.Error())
		return nil, &utils.NotFoundError{Message: "error retrieving reservation with id"}
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
		dep.logger.Error(ctx, err.Error())
		return 0, &utils.NotFoundError{Message: "error counting reservations in range"}
	}

	return count, nil
}

func (dep *reservationStore) Save(ctx context.Context, r *reservation.Reservation) error {
	if r == nil {
		return &utils.ServerError{Message: "reservation cannot be nil"}
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
		dep.logger.Error(ctx, err.Error())
		return &utils.InsertionError{Message: "error saving reservation"}
	}

	return nil
}