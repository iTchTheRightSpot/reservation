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

func (p *transactionProvider) RunInTransaction(txFunc func(adapters *Adapters) error) error {
	return p.runInTx(p.db, func(tx *sql.Tx) error {
		return txFunc(NewAdapters(p.logger, tx, nil))
	})
}

func (p *transactionProvider) runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
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
