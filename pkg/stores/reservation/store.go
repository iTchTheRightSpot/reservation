package reservation

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/utils"
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
	panic("implement me")
}

func (dep *reservationStore) SelectForUpdateSave(ctx context.Context, r *reservation.Reservation) error {
	//TODO implement me
	panic("implement me")
}
