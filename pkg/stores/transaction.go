package stores

import (
	"context"
	"database/sql"
	"errors"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type ITransactionProvider interface {
	RunInTransaction(ctx context.Context, iso *sql.TxOptions, txFunc func(adapters *Adapters) error) error
}

type transactionProvider struct {
	logger utils.ILogger
	db     *sql.DB
}

func NewTransactionProvider(l utils.ILogger, db *sql.DB) ITransactionProvider {
	return &transactionProvider{logger: l, db: db}
}

func (p *transactionProvider) RunInTransaction(ctx context.Context, iso *sql.TxOptions, txFunc func(adapters *Adapters) error) error {
	return p.runInTx(ctx, p.db, iso, func(tx *sql.Tx) error {
		return txFunc(NewAdapters(p.logger, tx, nil))
	})
}

func (p *transactionProvider) runInTx(ctx context.Context, db *sql.DB, iso *sql.TxOptions, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, iso)
	if err != nil {
		return err
	}

	if err = fn(tx); err == nil {
		return tx.Commit()
	}

	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
