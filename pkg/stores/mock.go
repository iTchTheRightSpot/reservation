package stores

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type mockLiveTransactionProvider struct {
	logger utils.ILogger
	tx     *sql.Tx
}

func MockLiveTransactionProvider(l utils.ILogger, tx *sql.Tx) ITransactionProvider {
	return &mockLiveTransactionProvider{logger: l, tx: tx}
}

func (p *mockLiveTransactionProvider) RunInTransaction(txFunc func(adapters *Adapters) error) error {
	return txFunc(NewAdapters(p.logger, p.tx, nil))
}
