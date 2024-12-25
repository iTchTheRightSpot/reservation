package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
)

func DeleteAll(db *sql.DB) error {
	if _, err := db.Exec("TRUNCATE schedule, staff, profile, service_type, staff_service, reservation, reservation_service CASCADE"); err != nil {
		return err
	}
	return nil
}

func PreSaveStaff(ctx context.Context, a *stores.Adapters) (*staff.Staff, error) {
	p := profile.Profile{
		Firstname: "erp",
		Lastname:  "erp",
		Email:     fmt.Sprintf("%s@email.com", uuid.NewString()),
	}

	if _, err := a.ProfileStore.Save(ctx, &p); err != nil {
		return nil, err
	}

	s := staff.Staff{
		UUID:      uuid.New(),
		ProfileId: &p.ProfileId,
	}

	if _, err := a.StaffStore.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}
