package reservation

import (
	"context"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"time"
)

type MockReservationStore struct {
	CountReservationsInRangeByStaffTimeAndStatusesReturn int
	CountReservationsInRangeByStaffTimeAndStatusesError  error
	CountReservationsInRangeByStaffTimeAndStatusesCalled bool
	SelectForUpdateSaveError                             error
	SelectForUpdateSaveCalled                            bool
}

func (dep *MockReservationStore) CountReservationsInRange(context.Context, uint64, time.Time, time.Time, ...reservation.ReservationEnum) (int, error) {
	dep.CountReservationsInRangeByStaffTimeAndStatusesCalled = true
	return dep.CountReservationsInRangeByStaffTimeAndStatusesReturn, dep.CountReservationsInRangeByStaffTimeAndStatusesError
}

func (dep *MockReservationStore) Save(context.Context, *reservation.Reservation) error {
	dep.SelectForUpdateSaveCalled = true
	return dep.SelectForUpdateSaveError
}

type MockReservationServiceStore struct {
	SaveReservationService       *reservation.ReservationServiceEntity
	SaveReservationServiceErr    error
	SaveReservationServiceCalled bool
}

func (dep *MockReservationServiceStore) Save(_ context.Context, r *reservation.ReservationServiceEntity) error {
	dep.SaveReservationServiceCalled = true
	r = dep.SaveReservationService
	return dep.SaveReservationServiceErr
}
