package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (p *transactionProvider) RunInTransaction(ctx context.Context, iso *sql.TxOptions, txFunc func(*Adapters) error) error {
	if iso != nil {
		p.logger.Log(fmt.Sprintf("transaction beginning isolation level options %v", iso))
	} else {
		p.logger.Log("transaction beginning isolation level options nil")
	}

	err := p.runInTx(ctx, p.db, iso, func(tx *sql.Tx) error { return txFunc(NewAdapters(p.logger, tx, nil)) })
	if err != nil {
		p.logger.Error("transaction not committed", err.Error())
		return err
	}

	p.logger.Log("transaction commited successfully")
	return nil
}

func (p *transactionProvider) runInTx(ctx context.Context, db *sql.DB, iso *sql.TxOptions, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, iso)
	if err != nil {
		p.logger.Error("error starting transaction", err.Error())
		return err
	}

	if err = fn(tx); err == nil {
		p.logger.Log("committing transaction")
		return tx.Commit()
	}

	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		p.logger.Error("failed to rollback transaction", "original error", err, "rollback error", rollbackErr)
		return errors.Join(err, rollbackErr)
	}

	p.logger.Error("transaction rolled back due to error", err)
	return err
}
