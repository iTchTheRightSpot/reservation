package reservation

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"strings"
	"time"
)

type IReservationStore interface {
	CountReservationsInRangeByStaffTimeAndStatuses(ctx context.Context, staffId uint64, start, end time.Time, statuses ...reservation.ReservationEnum) (int, error)
	SelectForUpdateSave(ctx context.Context, r *reservation.Reservation) error
}

type reservationStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewReservationStore(l utils.ILogger, db pkg.Db) IReservationStore {
	return &reservationStore{logger: l, db: db}
}

func (dep *reservationStore) CountReservationsInRangeByStaffTimeAndStatuses(ctx context.Context, staffId uint64, start, end time.Time, statuses ...reservation.ReservationEnum) (int, error) {
	arr := make([]string, len(statuses))
	for i, status := range statuses {
		arr[i] = string(status)
	}

	q := fmt.Sprintf(`
		SELECT COUNT(*) FROM reservation
		WHERE staff_id = $1 AND scheduled_for = $2
		AND expire_at = $3 AND status IN %s
	`, fmt.Sprintf("(%s)", strings.Join(arr, ", ")))

	var count int

	row := dep.db.QueryRowContext(ctx, q, staffId, start, end)
	if err := row.Scan(&count); err != nil {
		dep.logger.Error(err.Error())
		return 0, fmt.Errorf("error counting reservations in range by staff id, time and statues")
	}

	return count, nil
}

func (dep *reservationStore) SelectForUpdateSave(ctx context.Context, r *reservation.Reservation) error {
	arr := make([]string, len(statuses))
	for i, status := range statuses {
		arr[i] = string(status)
	}

	q := `
        WITH conflicting_reservations AS (
            SELECT reservation_id FROM reservation
            WHERE staff_id = $1
              AND scheduled_for BETWEEN $11 AND $12
              AND status IN ${statuses}
            FOR UPDATE  -- lock rows if any conflict exists
        )
        INSERT INTO reservation (staff_id, name, email, description, address, phone, image_key, price, status, created_at, scheduled_for, expire_at)
        SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
        WHERE NOT EXISTS (SELECT 1 FROM conflicting_reservations)
        RETURNING reservation_id, staff_id, name, email, description, address, phone, image_key, price, status, created_at, scheduled_for, expire_at;
    `

	return nil
}
