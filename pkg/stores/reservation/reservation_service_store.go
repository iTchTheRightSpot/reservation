package reservation

import (
	"context"
	"github.com/iTchTheRightSpot/reservation/pkg"
	"github.com/iTchTheRightSpot/reservation/pkg/models/reservation"
	"github.com/iTchTheRightSpot/utility/utils"
)

type IReservationServiceStore interface {
	Save(ctx context.Context, entity *reservation.ReservationServiceEntity) error
}

type reservationServiceStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewReservationServiceStore(l utils.ILogger, db pkg.Db) IReservationServiceStore {
	return &reservationServiceStore{logger: l, db: db}
}

func (dep *reservationServiceStore) Save(ctx context.Context, e *reservation.ReservationServiceEntity) error {
	q := `
		INSERT INTO reservation_service (reservation_id, service_id)
		VALUES ($1, $2)
		RETURNING junction_id, reservation_id, service_id
	`

	row := dep.db.QueryRowContext(ctx, q, e.ReservationId, e.ServiceId)
	if err := row.Scan(&e.JunctionId, &e.ReservationId, &e.ServiceId); err != nil {
		dep.logger.Error(ctx, err.Error())
		return &utils.InsertionError{}
	}
	return nil
}