package stores

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type Adapters struct {
	ProfileStore profile.IProfileStore
	Transaction  ITransactionProvider
}

func NewAdapters(l utils.ILogger, db utils.Db, p ITransactionProvider) *Adapters {
	return &Adapters{
		ProfileStore: profile.NewProfileStore(l, db),
		Transaction:  p,
	}
}
