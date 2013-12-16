package postgres

import (
	"database/sql"
	"errors"
	"github.com/aktau/gomig/db/common"
	"github.com/lib/pq"
	"log"
)

var (
	PG_DB_EXECUTOR_VERBOSE = true
)

type PgDbExecutor struct {
	common.DbExecutor
	bulkStmt *sql.Stmt
}

func NewPgDbExecutor(db *sql.DB) (*PgDbExecutor, error) {
	base, err := common.NewDbExecutor(db)
	if err != nil {
		return nil, err
	}

	return &PgDbExecutor{*base, nil}, nil
}

func (e *PgDbExecutor) BulkInit(table string, columns ...string) error {
	db := e.GetDb()
	if db == nil {
		return errors.New("executor did not have a valid database")
	}

	var (
		stmt *sql.Stmt
		err  error
	)
	tx := e.GetTx()
	copySql := pq.CopyIn(table, columns...)
	if tx == nil {
		stmt, err = db.Prepare(copySql)
	} else {
		stmt, err = tx.Prepare(copySql)
	}
	if err != nil {
		return err
	}
	e.bulkStmt = stmt

	return nil
}

func (e *PgDbExecutor) BulkAddRecord(args ...interface{}) error {
	/* TODO: does not check if bulkStmt exists yet */
	_, err := e.bulkStmt.Exec(args...)
	return err
}

func (e *PgDbExecutor) BulkFinish() error {
	stmt := e.bulkStmt
	defer func() { e.bulkStmt = nil }()

	_, err := stmt.Exec()
	if err != nil {
		cerr := stmt.Close()
		if cerr != nil {
			log.Println("pg_executor: could not properly close bulk statement", cerr)
		}
		return err
	}

	err = stmt.Close()
	return err
}

func (e *PgDbExecutor) HasCapability(capability int) bool {
	return capability == common.CapBulkTransfer
}
