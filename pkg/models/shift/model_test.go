package shift

import (
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"testing"
	"time"
)

func TestShift(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()
	timezone := logger.Timezone()

	t.Run("return error as dto contains duplicate times", func(t *testing.T) {
		t.Parallel()

		now := logger.Date()
		start1 := now.Add(time.Duration(30) * time.Minute).Format(utils.TimeFormat)

		timeSlots := []ShiftSegment{
			{IsVisible: true, Start: now.Format(utils.TimeFormat), Duration: 3600},
			{IsVisible: false, Start: start1, Duration: 3600},
		}

		dto := ShiftDto{
			StaffUUID: "test-staff-uuid",
			Times:     &timeSlots,
		}

		_, err := dto.CheckForOverlappingSegments(timezone)
		if err == nil {
			t.Fatalf("expected an error due to duplicate time slots but got none")
		}
	})

	t.Run("no duplicate times", func(t *testing.T) {
		t.Parallel()

		now := logger.Date()
		start1 := now.Add(time.Duration(1) * time.Hour).Format(utils.TimeFormat)

		timeSlots := []ShiftSegment{
			{IsVisible: true, Start: now.Format(utils.TimeFormat), Duration: 3600},
			{IsVisible: false, Start: start1, Duration: 3600},
		}

		dto := ShiftDto{
			StaffUUID: "test-staff-uuid",
			Times:     &timeSlots,
		}

		result, err := dto.CheckForOverlappingSegments(timezone)
		if err != nil {
			t.Fatalf("did not expect an error but got: %v", err)
		}
		if len(result) != len(timeSlots) {
			t.Errorf("expected %d time slots in result, got %d", len(timeSlots), len(result))
		}
	})

	t.Run("return error as time slot contains invalid format", func(t *testing.T) {
		t.Parallel()

		timeSlots := []ShiftSegment{
			{IsVisible: true, Start: "invalidTime", Duration: 3600},
		}

		dto := ShiftDto{
			StaffUUID: "test-staff-uuid",
			Times:     &timeSlots,
		}

		_, err := dto.CheckForOverlappingSegments(timezone)
		if err == nil {
			t.Fatalf("expected an error due to invalid time format but got none")
		}
	})

	t.Run("return error as time slots is nil", func(t *testing.T) {
		t.Parallel()

		dto := ShiftDto{
			StaffUUID: "test-staff-uuid",
			Times:     nil,
		}

		_, err := dto.CheckForOverlappingSegments(timezone)
		if err == nil {
			t.Fatalf("expected an error due to nil time slots but got none")
		}
	})
}
