package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IServiceStore interface {
	Save(ctx context.Context, s *service.ServiceEntity) (*service.ServiceEntity, error)
	ServiceByName(ctx context.Context, name string) (*service.ServiceEntity, error)
	ServicesByStaffId(ctx context.Context, staffId uint64) ([]*service.ServiceEntity, error)
}

type serviceStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewServiceStore(l utils.ILogger, db pkg.Db) IServiceStore {
	return &serviceStore{logger: l, db: db}
}

func (dep *serviceStore) Save(ctx context.Context, s *service.ServiceEntity) (*service.ServiceEntity, error) {
	if s == nil {
		return nil, fmt.Errorf("schedule object is nil")
	}

	q := `
		INSERT INTO service (name, price, is_visible, is_reoccurring, duration, clean_up_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING service_id, name, price, is_visible, is_reoccurring, duration, clean_up_time
	`

	row := dep.db.QueryRowContext(ctx, q, s.Name, s.Price, s.IsVisible, s.IsReoccurring, s.Duration, s.CleanUpTime)
	err := row.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.IsReoccurring, &s.Duration, &s.CleanUpTime)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to service table")
	}
	dep.logger.Log("new service saved")
	return s, nil
}

func (dep *serviceStore) ServiceByName(ctx context.Context, name string) (*service.ServiceEntity, error) {
	var s service.ServiceEntity

	var q = `
		SELECT
		    service_id, name, price, is_visible, is_reoccurring, duration, clean_up_time
		FROM service WHERE name = $1
	`

	row := dep.db.QueryRowContext(ctx, q, name)
	err := row.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.IsReoccurring, &s.Duration, &s.CleanUpTime)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception retrieving ServiceByName %s", name)
	}

	return &s, nil
}

// ServicesByStaffId TODO test
func (dep *serviceStore) ServicesByStaffId(ctx context.Context, staffId uint64) ([]*service.ServiceEntity, error) {
	var arr []*service.ServiceEntity

	var q = `
	 SELECT s.* FROM service s
	 INNER JOIN staff_service ss ON ss.service_id = s.service_id
	 WHERE ss.staff_id = $1
	`

	rows, err := dep.db.QueryContext(ctx, q, staffId)
	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception retrieving services")
	}

	defer func(rows *sql.Rows) { err = rows.Close() }(rows)

	for rows.Next() {
		var s service.ServiceEntity

		err = rows.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.IsReoccurring, &s.Duration, &s.CleanUpTime)

		if err != nil {
			dep.logger.Error(err)
			return nil, fmt.Errorf("exception scanning service")

		}

		arr = append(arr, &s)
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("error iterating services")
	}

	return arr, err
}
