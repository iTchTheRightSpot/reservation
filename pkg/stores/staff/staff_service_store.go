package staff

import (
	"context"
	"errors"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IStaffServiceStore interface {
	Save(ctx context.Context, s *staff.StaffServiceEntity) error
	CountByStaffIdAndServiceId(ctx context.Context, staffId, serviceId uint64) (int, error)
}

type staffServiceStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewStaffServiceStore(l utils.ILogger, db pkg.Db) IStaffServiceStore {
	return &staffServiceStore{logger: l, db: db}
}

func (dep *staffServiceStore) Save(ctx context.Context, s *staff.StaffServiceEntity) error {
	if s == nil {
		return errors.New("StaffServiceEntity cannot be nil")
	}
	q := `
        INSERT INTO staff_service (staff_id, service_id)
        VALUES ($1, $2)
        RETURNING junction_id, staff_id, service_id
    `

	row := dep.db.QueryRowContext(ctx, q, s.StaffId, s.ServiceId)
	if err := row.Scan(&s.JunctionId, &s.StaffId, &s.ServiceId); err != nil {
		dep.logger.Error(err.Error())
		return errors.New("exception linking service to staff")
	}

	dep.logger.Log("saved to staff_service table")
	return nil
}

func (dep *staffServiceStore) CountByStaffIdAndServiceId(ctx context.Context, staffId, serviceId uint64) (int, error) {
	var count int

	q := "SELECT COUNT(*) FROM staff_service WHERE staff_id = $1 AND service_id = $2"

	row := dep.db.QueryRowContext(ctx, q, staffId, serviceId)
	if err := row.Scan(&count); err != nil {
		dep.logger.Error(err.Error())
		return 0, errors.New("error count staff service")
	}

	return count, nil
}
