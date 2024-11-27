package reservation

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IReservationStore interface {
}

type reservationStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewReservationStore(l utils.ILogger, db pkg.Db) IReservationStore {
	return &reservationStore{logger: l, db: db}
}
