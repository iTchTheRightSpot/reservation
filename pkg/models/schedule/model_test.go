package schedule

import (
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"testing"
	"time"
)

func TestScheduleModel(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()
	timezone := logger.Timezone()

	t.Run("return error schedule in the past", func(t *testing.T) {
		t.Parallel()

		// given
		now := logger.Date()
		timeSlots := []ScheduleSegmentPayload{
			{IsVisible: true, Start: now.Format(utils.TimeFormat), Duration: 3600},
		}
		dto := SchedulePayload{
			StaffId: "test-staff-uuid",
			Times:   &timeSlots,
		}
		future := now.Add(time.Duration(10) * time.Minute)

		// method to test
		_, err := dto.CheckForOverlappingSegments(future, timezone)

		// assert
		if err == nil {
			t.Fatalf("expected an error due to duplicate time slots but got none")
		}

		// assert appropriate error
		s := err.Error()
		g := fmt.Sprintf("%s is in the past ", timeSlots[0].Start)

		if s != g {
			t.Errorf("expect %s to equal given %s", s, g)
		}
	})

	t.Run("return error dto contains duplicate times", func(t *testing.T) {
		t.Parallel()

		// given
		now := logger.Date()
		start1 := now.Add(time.Duration(30) * time.Minute).Format(utils.TimeFormat)

		timeSlots := []ScheduleSegmentPayload{
			{IsVisible: true, Start: now.Format(utils.TimeFormat), Duration: 3600},
			{IsVisible: false, Start: start1, Duration: 3600},
		}

		dto := SchedulePayload{
			StaffId: "test-staff-uuid",
			Times:   &timeSlots,
		}

		sub := now.Add(time.Duration(-3) * time.Hour)

		// method to test
		_, err := dto.CheckForOverlappingSegments(sub, timezone)

		// assert
		if err == nil {
			t.Fatalf("expected an error due to duplicate time slots but got none")
		}

		// assert appropriate error
		s := err.Error()
		g := fmt.Sprintf("conflicting date %s & duration %v", timeSlots[1].Start, timeSlots[1].Duration)
		if s != g {
			t.Errorf("expect %s to equal given %s", s, g)
		}
	})

	t.Run("success no duplicate times", func(t *testing.T) {
		t.Parallel()

		// given
		now := logger.Date()
		start1 := now.Add(time.Duration(1) * time.Hour).Format(utils.TimeFormat)
		timeSlots := []ScheduleSegmentPayload{
			{IsVisible: true, Start: now.Format(utils.TimeFormat), Duration: 3600},
			{IsVisible: false, Start: start1, Duration: 3600},
		}
		dto := SchedulePayload{
			StaffId: "test-staff-uuid",
			Times:   &timeSlots,
		}
		sub := now.Add(time.Duration(-3) * time.Hour)

		// method to test
		result, err := dto.CheckForOverlappingSegments(sub, timezone)

		// assert
		if err != nil {
			t.Fatalf("did not expect an error but got: %v", err)
		}

		if len(result) != len(timeSlots) {
			t.Errorf("expected %d time slots in result, got %d", len(timeSlots), len(result))
		}
	})

	t.Run("return error as time slot contains invalid format", func(t *testing.T) {
		t.Parallel()

		timeSlots := []ScheduleSegmentPayload{
			{IsVisible: true, Start: "invalidTime", Duration: 3600},
		}

		dto := SchedulePayload{
			StaffId: "test-staff-uuid",
			Times:   &timeSlots,
		}

		_, err := dto.CheckForOverlappingSegments(logger.Date(), timezone)
		if err == nil {
			t.Fatalf("expected an error due to invalid time format but got none")
		}
	})

	t.Run("return error as time slots is nil", func(t *testing.T) {
		t.Parallel()

		dto := SchedulePayload{
			StaffId: "test-staff-uuid",
			Times:   nil,
		}

		_, err := dto.CheckForOverlappingSegments(logger.Date(), timezone)
		if err == nil {
			t.Fatalf("expected an error due to nil time slots but got none")
		}
	})
}
