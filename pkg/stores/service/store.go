package service

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IServiceStore interface {
	Save(ctx context.Context, s *service.Service) (*service.Service, error)
}

type serviceStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewServiceStore(l utils.ILogger, db pkg.Db) IServiceStore {
	return &serviceStore{logger: l, db: db}
}

func (dep *serviceStore) Save(ctx context.Context, s *service.Service) (*service.Service, error) {
	if s == nil {
		return nil, fmt.Errorf("shift object is nil")
	}

	q := `
		INSERT INTO service (name, price, is_visible, duration, clean_up_time)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING service_id, name, price, is_visible, duration, clean_up_time
	`

	row := dep.db.QueryRowContext(ctx, q, s.Name, s.Price, s.IsVisible, s.Duration, s.CleanUpTime)
	err := row.Scan(&s.ServiceId, &s.Name, &s.Price, &s.IsVisible, &s.Duration, &s.CleanUpTime)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to service table")
	}

	return s, nil
}
