package reservation

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/utils"
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

func (dep *reservationServiceStore) Save(ctx context.Context, entity *reservation.ReservationServiceEntity) error {
	//TODO implement me
	panic("implement me")
}
