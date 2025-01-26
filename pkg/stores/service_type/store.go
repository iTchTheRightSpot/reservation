package service_type

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service_type"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IServiceTypeStore interface {
	Save(ctx context.Context, s *service_type.ServiceTypeEntity) error
	ServiceTypeByName(ctx context.Context, name string) (*service_type.ServiceTypeEntity, error)
	ServiceTypesByStaffId(ctx context.Context, staffId uint64, visible bool) ([]*service_type.ServiceTypeEntity, error)
	ServiceTypes(ctx context.Context) ([]*service_type.ServiceTypeEntity, error)
	Update(ctx context.Context, s *service_type.ServiceTypeEntity) error
}

type serviceTypeStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewServiceTypeStore(l utils.ILogger, db pkg.Db) IServiceTypeStore {
	return &serviceTypeStore{logger: l, db: db}
}

func (dep *serviceTypeStore) Update(ctx context.Context, s *service_type.ServiceTypeEntity) error {
	q := `
		UPDATE service_type
		SET name = $2, price = $3, is_visible = $4, duration = $5, clean_up_time = $6
		WHERE service_id = $1
	`
	_, err := dep.db.ExecContext(ctx, q, s.ServiceId, s.Name, s.Price, s.IsVisible, s.Duration, s.CleanUpTime)
	if err != nil {
		dep.logger.Error(err.Error())
		return errors.New("error updating service type")
	}
	return nil
}

func (dep *serviceTypeStore) ServiceTypes(ctx context.Context) ([]*service_type.ServiceTypeEntity, error) {
	var arr []*service_type.ServiceTypeEntity

	rows, err := dep.db.QueryContext(ctx, "SELECT * FROM service_type")
	if err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("exception retrieving services by status")
	}

	defer func(rows *sql.Rows) { err = rows.Close() }(rows)

	for rows.Next() {
		var s service_type.ServiceTypeEntity

		err = rows.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.Duration, &s.CleanUpTime)

		if err != nil {
			dep.logger.Error(err)
			return nil, errors.New("exception scanning service by status")

		}

		arr = append(arr, &s)
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(err.Error())
		return nil, errors.New("error iterating services by status")
	}

	return arr, err
}

func (dep *serviceTypeStore) Save(ctx context.Context, s *service_type.ServiceTypeEntity) error {
	if s == nil {
		return errors.New("service type entity is nil")
	}

	q := `
		INSERT INTO service_type (name, price, is_visible, duration, clean_up_time)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING service_id, name, price, is_visible, duration, clean_up_time
	`

	row := dep.db.QueryRowContext(ctx, q, s.Name, s.Price, s.IsVisible, s.Duration, s.CleanUpTime)
	err := row.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.Duration, &s.CleanUpTime)

	if err != nil {
		dep.logger.Error(err)
		return errors.New("exception saving to service table")
	}
	dep.logger.Log("new service saved")
	return nil
}

func (dep *serviceTypeStore) ServiceTypeByName(ctx context.Context, name string) (*service_type.ServiceTypeEntity, error) {
	var s service_type.ServiceTypeEntity

	var q = `
		SELECT
		    service_id, name, price, is_visible, duration, clean_up_time
		FROM service_type WHERE name = $1
	`

	row := dep.db.QueryRowContext(ctx, q, name)
	err := row.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.Duration, &s.CleanUpTime)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception retrieving ServiceTypeByName %s", name)
	}

	return &s, nil
}

func (dep *serviceTypeStore) ServiceTypesByStaffId(ctx context.Context, staffId uint64, visible bool) ([]*service_type.ServiceTypeEntity, error) {
	var arr []*service_type.ServiceTypeEntity

	var q = `
	 SELECT s.* FROM service_type s
	 INNER JOIN staff_service ss ON ss.service_id = s.service_id
	 WHERE ss.staff_id = $1 AND s.is_visible = $2
	`

	rows, err := dep.db.QueryContext(ctx, q, staffId, visible)
	if err != nil {
		dep.logger.Error(err)
		return nil, errors.New("exception retrieving services")
	}

	defer func(rows *sql.Rows) { err = rows.Close() }(rows)

	for rows.Next() {
		var s service_type.ServiceTypeEntity

		err = rows.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.Duration, &s.CleanUpTime)

		if err != nil {
			dep.logger.Error(err)
			return nil, errors.New("exception scanning service")

		}

		arr = append(arr, &s)
	}

	if err = rows.Err(); err != nil {
		dep.logger.Error(err)
		return nil, errors.New("error iterating services")
	}

	return arr, err
}
