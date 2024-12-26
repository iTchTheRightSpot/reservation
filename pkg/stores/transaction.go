package stores

import (
	"database/sql"
	"errors"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

type ITransactionProvider interface {
	RunInTransaction(txFunc func(adapters *Adapters) error) error
}

type transactionProvider struct {
	logger utils.ILogger
	db     *sql.DB
}

func NewTransactionProvider(l utils.ILogger, db *sql.DB) ITransactionProvider {
	return &transactionProvider{logger: l, db: db}
}

func (p *transactionProvider) RunInTransaction(txFunc func(*Adapters) error) error {
	p.logger.Log("transaction beginning isolation level options nil")

	err := p.runInTx(p.db, func(tx *sql.Tx) error { return txFunc(NewAdapters(p.logger, tx, nil)) })
	if err != nil {
		p.logger.Error("transaction not committed", err.Error())
		return err
	}

	p.logger.Log("transaction commited successfully")
	return nil
}

func (p *transactionProvider) runInTx(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
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
