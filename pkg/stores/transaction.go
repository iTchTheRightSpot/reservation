package stores

import (
	"database/sql"
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
	err := p.runInTx(p.db, func(tx *sql.Tx) error { return txFunc(NewAdapters(p.logger, tx, nil)) })
	if err != nil {
		p.logger.Error("transaction not committed", err.Error())
		return err
	}

	p.logger.Log("transaction commited successfully")
	return nil
}

func (p *transactionProvider) runInTx(db *sql.DB, fn func(*sql.Tx) error) error {
	p.logger.Log("BEGINNING TRANSACTION")

	tx, err := db.Begin()
	if err != nil {
		p.logger.Error("ERROR STARTING TRANSACTION", err.Error())
		return &utils.ServerError{Message: "error starting transaction"}
	}

	if err = fn(tx); err == nil {
		p.logger.Log("COMMITTING TRANSACTION")
		if err = tx.Commit(); err != nil {
			p.logger.Error(err.Error())
			return &utils.InsertionError{Message: "error committing transaction"}
		}
		p.logger.Log("TRANSACTION COMMITTED SUCCESSFULLY")
		return nil
	}

	p.logger.Error("TRANSACTION FAILED", err.Error())
	p.logger.Log("BEGINNING TRANSACTION ROLLBACK")

	if rbe := tx.Rollback(); rbe != nil {
		p.logger.Error("FAILED TO ROLLBACK TRANSACTION", rbe.Error())
		return err
	}

	p.logger.Log("TRANSACTION ROLLED BACK SUCCESSFULLY")
	return err
}
