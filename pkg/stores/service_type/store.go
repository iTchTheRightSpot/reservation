package service_type

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IServiceTypeStore interface {
	Save(ctx context.Context, s *service.ServiceTypeEntity) error
	ServiceByName(ctx context.Context, name string) (*service.ServiceTypeEntity, error)
	ServicesByStaffId(ctx context.Context, staffId uint64, visible bool) ([]*service.ServiceTypeEntity, error)
}

type serviceStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewServiceTypeStore(l utils.ILogger, db pkg.Db) IServiceTypeStore {
	return &serviceStore{logger: l, db: db}
}

func (dep *serviceStore) Save(ctx context.Context, s *service.ServiceTypeEntity) error {
	if s == nil {
		return errors.New("service type entity is nil")
	}

	q := `
		INSERT INTO service_type (name, price, is_visible, is_reoccurring, duration, clean_up_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING service_id, name, price, is_visible, is_reoccurring, duration, clean_up_time
	`

	row := dep.db.QueryRowContext(ctx, q, s.Name, s.Price, s.IsVisible, s.IsReoccurring, s.Duration, s.CleanUpTime)
	err := row.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.IsReoccurring, &s.Duration, &s.CleanUpTime)

	if err != nil {
		dep.logger.Error(err)
		return errors.New("exception saving to service table")
	}
	dep.logger.Log("new service saved")
	return nil
}

func (dep *serviceStore) ServiceByName(ctx context.Context, name string) (*service.ServiceTypeEntity, error) {
	var s service.ServiceTypeEntity

	var q = `
		SELECT
		    service_id, name, price, is_visible, is_reoccurring, duration, clean_up_time
		FROM service_type WHERE name = $1
	`

	row := dep.db.QueryRowContext(ctx, q, name)
	err := row.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.IsReoccurring, &s.Duration, &s.CleanUpTime)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception retrieving ServiceByName %s", name)
	}

	return &s, nil
}

func (dep *serviceStore) ServicesByStaffId(ctx context.Context, staffId uint64, visible bool) ([]*service.ServiceTypeEntity, error) {
	var arr []*service.ServiceTypeEntity

	var q = `
	 SELECT s.* FROM service_type s
	 INNER JOIN staff_service ss ON ss.service_id = s.service_id
	 WHERE ss.staff_id = $1 AND s.is_visible = $2
	`

	rows, err := dep.db.QueryContext(ctx, q, staffId, visible)
	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception retrieving services")
	}

	defer func(rows *sql.Rows) { err = rows.Close() }(rows)

	for rows.Next() {
		var s service.ServiceTypeEntity

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
