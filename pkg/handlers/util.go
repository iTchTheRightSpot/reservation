package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
)

func DeleteAll(db *sql.DB) error {
	if _, err := db.Exec("TRUNCATE permission, role, profile, schedule, staff_service, reservation_service, payment_detail, reservation, service_type, staff"); err != nil {
		return err
	}
	return nil
}

func PreSaveStaff(ctx context.Context, a *stores.Adapters) (*staff.Staff, error) {
	p := models.ProfileEntity{
		Firstname: "reservation",
		Lastname:  "reservation",
		Email:     fmt.Sprintf("%s@email.com", uuid.NewString()),
		Password:  "password",
	}

	if err := a.ProfileStore.Save(ctx, &p); err != nil {
		return nil, err
	}

	s := staff.Staff{
		UUID:      uuid.New(),
		ProfileId: &p.ProfileId,
	}

	if err := a.StaffStore.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func CountResponseStatus(arr []int, status int) int {
	var n int
	for _, num := range arr {
		if num == status {
			n += 1
		}
	}
	return n
}
