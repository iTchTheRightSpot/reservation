package stores

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/utility/utils"
)

type ITransactionProvider interface {
	RunInTransaction(ctx context.Context, txFunc func(adapters *Adapters) error) error
}

type transactionProvider struct {
	logger utils.ILogger
	db     *sql.DB
}

func NewTransactionProvider(l utils.ILogger, db *sql.DB) ITransactionProvider {
	return &transactionProvider{logger: l, db: db}
}

func (p *transactionProvider) RunInTransaction(ctx context.Context, txFunc func(*Adapters) error) error {
	return utils.RunInTx(ctx, p.logger, p.db, func(tx *sql.Tx) error { return txFunc(NewAdapters(p.logger, p.db, nil)) })
}
