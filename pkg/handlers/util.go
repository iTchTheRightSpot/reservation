package handlers

import "database/sql"

func DeleteAll(db *sql.DB) error {
	if _, err := db.Exec("TRUNCATE schedule, staff, profile, service_type, staff_service, reservation, reservation_service CASCADE"); err != nil {
		return err
	}
	return nil
}
