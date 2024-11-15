package profile

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type IProfileStore interface {
	Save(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
}

type profileStore struct {
	logger utils.ILogger
	db     pkg.Db
}

func NewProfileStore(l utils.ILogger, db pkg.Db) IProfileStore {
	return &profileStore{logger: l, db: db}
}

func (dep *profileStore) Save(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	if p == nil {
		return nil, fmt.Errorf("profile object is nil")
	}

	q := `
		INSERT INTO profile (firstname, lastname, email, image_key)
        VALUES ($1, $2, $3, $4)
        RETURNING profile_id, firstname, lastname, email, image_key
	`

	row := dep.db.QueryRowContext(ctx, q, p.Firstname, p.Lastname, p.Email, p.ImageKey)
	err := row.Scan(&p.ProfileId, &p.Firstname, &p.Lastname, &p.Email, &p.ImageKey)

	if err != nil {
		dep.logger.Error(err)
		return nil, fmt.Errorf("exception saving to profile table")
	}

	return p, nil
}
