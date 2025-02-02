package reservation

import (
	"context"
	"fmt"
	reservationModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	serviceModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/service_type"
	staffModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/service_type"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"testing"
	"time"
)

func TestReservationService(t *testing.T) {
	t.Parallel()

	mockLog := utils.NewMockLogger()

	t.Run("creating a reservation", func(t *testing.T) {
		t.Parallel()

		t.Run("reject creation staff not found", func(t *testing.T) {
			t.Parallel()

			serviceStore := &service_type.MockServiceTypeStore{}
			adapters := &stores.Adapters{
				StaffStore:   &staff.MockStaffStore{StaffByUUIDError: fmt.Errorf("staff not found")},
				ServiceStore: serviceStore,
			}
			s := NewReservationService(mockLog, adapters, nil, nil)

			// method to test
			err := s.Create(context.Background(), &reservationModel.ReservationPayload{})

			// assert
			if err == nil {
				t.Errorf("expect error but given nil")
			}

			if err.Error() != "invalid staff id" {
				t.Errorf("expect %s given %s", "staff not found", err.Error())
			}

			if serviceStore.ServicesByStaffIdCalled {
				t.Errorf("ServicesByStaffIdCalled expect false given true")
			}
		})

		t.Run("reject creation matchStaffServices not found", func(t *testing.T) {
			t.Parallel()

			scheduleStore := &schedule.MockScheduleStore{}
			adapters := &stores.Adapters{
				StaffStore:    &staff.MockStaffStore{StaffByUUIDReturn: &staffModel.StaffEntity{}},
				ServiceStore:  &service_type.MockServiceTypeStore{ServicesByStaffIdError: fmt.Errorf("no matchStaffServices")},
				ScheduleStore: scheduleStore,
			}
			s := NewReservationService(mockLog, adapters, nil, nil)

			// method to test
			err := s.Create(context.Background(), &reservationModel.ReservationPayload{})

			// assert
			if err == nil {
				t.Errorf("expect error but given nil")
			}

			if err.Error() != "invalid service for staff" {
				t.Errorf("expect %s given %s", "invalid service for staff", err.Error())
			}

			if scheduleStore.CountSchedulesInRangeAndVisibilityCalled {
				t.Errorf("CountSchedulesInRangeAndVisibilityCalled expect false given true")
			}
		})

		t.Run("reject creation staff does not offer any service", func(t *testing.T) {
			t.Parallel()

			// given
			erp := "erp"
			scheduleStore := &schedule.MockScheduleStore{}
			adapters := &stores.Adapters{
				StaffStore:    &staff.MockStaffStore{StaffByUUIDReturn: &staffModel.StaffEntity{}},
				ServiceStore:  &service_type.MockServiceTypeStore{ServicesByStaffIdReturn: []*serviceModel.ServiceTypeEntity{}},
				ScheduleStore: scheduleStore,
			}
			s := NewReservationService(mockLog, adapters, nil, nil)

			// method to test
			err := s.Create(context.Background(), &reservationModel.ReservationPayload{Services: []string{erp}})

			// assert
			if err == nil {
				t.Errorf("expect error but given nil")
			}

			errMess := "1 or more services were not found for selected staff"
			if err.Error() != errMess {
				t.Errorf("expect %s given %s", errMess, err.Error())
			}

			if scheduleStore.CountSchedulesInRangeAndVisibilityCalled {
				t.Errorf("CountSchedulesInRangeAndVisibilityCalled expect false given true")
			}
		})

		t.Run("reject creation 1 or more service staff does not offer", func(t *testing.T) {
			t.Parallel()

			// given
			erp := "erp"
			accounting := "accounting"
			payload := &reservationModel.ReservationPayload{
				Services: []string{erp, accounting},
			}
			scheduleStore := &schedule.MockScheduleStore{}
			adapters := &stores.Adapters{
				StaffStore: &staff.MockStaffStore{StaffByUUIDReturn: &staffModel.StaffEntity{}},
				ServiceStore: &service_type.MockServiceTypeStore{
					ServicesByStaffIdReturn: []*serviceModel.ServiceTypeEntity{{Name: "erp"}}},
				ScheduleStore: scheduleStore,
			}
			s := NewReservationService(mockLog, adapters, nil, nil)

			// method to test
			err := s.Create(context.Background(), payload)

			// assert
			if err == nil {
				t.Errorf("expect error but given nil")
			}

			errMess := "1 or more services were not found for selected staff"
			if err.Error() != errMess {
				t.Errorf("expect %s given %s", errMess, err.Error())
			}

			if scheduleStore.CountSchedulesInRangeAndVisibilityCalled {
				t.Errorf("CountSchedulesInRangeAndVisibilityCalled expect false given true")
			}
		})

		t.Run("reject creation. attempting to reserve a day in the past", func(t *testing.T) {
			t.Parallel()

			// given
			erp := "erp"
			payload := &reservationModel.ReservationPayload{
				Services: []string{erp},
				Time:     fmt.Sprintf("%v", time.Now().Add(-1*time.Hour).UnixMilli()),
			}
			scheduleStore := &schedule.MockScheduleStore{CountSchedulesInRangeAndVisibilityError: fmt.Errorf("err")}
			adapters := &stores.Adapters{
				StaffStore: &staff.MockStaffStore{StaffByUUIDReturn: &staffModel.StaffEntity{}},
				ServiceStore: &service_type.MockServiceTypeStore{
					ServicesByStaffIdReturn: []*serviceModel.ServiceTypeEntity{{Name: erp}}},
				ScheduleStore: scheduleStore,
			}
			s := NewReservationService(mockLog, adapters, nil, nil)

			// method to test
			err := s.Create(context.Background(), payload)

			// assert
			if err == nil {
				t.Errorf("expect error but given nil")
			}

			errMes := "cannot make a reservation for a past day"
			if err.Error() != errMes {
				t.Errorf("expect %s given %s", errMes, err.Error())
			}
		})

		t.Run("reject creation. reservation time does not work with start schedule", func(t *testing.T) {
			t.Parallel()

			// given
			erp := "erp"
			payload := &reservationModel.ReservationPayload{
				Services: []string{erp},
				Time:     fmt.Sprintf("%v", time.Now().Add(24*time.Hour).UnixMilli()),
			}
			reservationStore := &reservation.MockReservationStore{}
			adapters := &stores.Adapters{
				StaffStore: &staff.MockStaffStore{StaffByUUIDReturn: &staffModel.StaffEntity{}},
				ServiceStore: &service_type.MockServiceTypeStore{
					ServicesByStaffIdReturn: []*serviceModel.ServiceTypeEntity{{Name: erp}}},
				ScheduleStore:    &schedule.MockScheduleStore{CountSchedulesInRangeAndVisibilityReturn: 0},
				ReservationStore: reservationStore,
			}
			s := NewReservationService(mockLog, adapters, nil, nil)

			// method to test
			err := s.Create(context.Background(), payload)

			// assert
			if err == nil {
				t.Errorf("expect error but given nil")
			}

			mess := "invalid reservation time"
			if err.Error() != mess {
				t.Errorf("expect %s given %s", mess, err.Error())
			}

			if reservationStore.SelectForUpdateSaveCalled {
				t.Errorf("Save expect false given true")
			}
		})
	})
}
