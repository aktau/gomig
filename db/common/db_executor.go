package common

import (
	"database/sql"
	"fmt"
	"log"
)

var (
	DBEXEC_VERBOSE = true
)

type DbExecutor struct {
	db *sql.DB
	tx *sql.Tx

	err func(err error) error // can create for more helpful messages
}

func identity(err error) error { return err }

func NewDbExecutor(db *sql.DB, err func(error) error) (*DbExecutor, error) {
	if err == nil {
		err = identity
	}
	return &DbExecutor{db: db, tx: nil, err: err}, nil
}

func (e *DbExecutor) Begin(name string) error {
	if e.tx != nil {
		return ErrTxInProgress
	}

	if DBEXEC_VERBOSE {
		log.Printf("DbExecutor: starting transaction: %v", name)
	}

	/* start transaction */
	tx, err := e.db.Begin()
	if err != nil {
		return e.err(err)
	}

	e.tx = tx
	return nil
}

func (e *DbExecutor) Commit() error {
	if e.tx == nil {
		return ErrNoTxInProgress
	}
	defer func() { e.tx = nil }()

	return e.err(e.tx.Commit())
}

func (e *DbExecutor) Rollback() error {
	if e.tx == nil {
		return ErrNoTxInProgress
	}
	defer func() { e.tx = nil }()

	rerr := e.err(e.tx.Rollback())
	if rerr != nil && DBEXEC_VERBOSE {
		log.Printf("DbExecutor: error while rolling back: %v", rerr)
	}

	return rerr
}

func (e *DbExecutor) submitSimple(stmt string) error {
	if _, err := e.db.Exec(stmt); err != nil {
		err = e.err(err)
		return fmt.Errorf("'%v' while executing statement\n'%v'", err, stmt)
	}
	return nil
}

func (e *DbExecutor) submitTransactional(stmt string) error {
	_, err := e.tx.Exec(stmt)
	if err != nil {
		err = e.err(err)
		e.Rollback()
		return fmt.Errorf("'%v' while executing statement\n%v in transaction", err, stmt)
	}
	return nil
}

func (e *DbExecutor) Submit(stmt string) error {
	if DBEXEC_VERBOSE {
		log.Println(stmt)
	}

	if e.tx == nil {
		return e.submitSimple(stmt)
	} else {
		return e.submitTransactional(stmt)
	}
}

func (e *DbExecutor) Multiple(name string, statements []string) []error {
	errors := make([]error, 0, len(statements))

	if DBEXEC_VERBOSE {
		log.Printf("DbExecutor: starting multiple statements: %v", name)
	}

	/* write out all statements, rollback in case of error */
	for _, stmt := range statements {
		err := e.Submit(stmt)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (e *DbExecutor) Transaction(name string, statements []string) error {
	/* start transaction */
	err := e.Begin(name)
	if err != nil {
		return err
	}

	/* write out all statements, rollback in case of error */
	for _, stmt := range statements {
		err := e.Submit(stmt)
		if err != nil {
			return err
		}
	}

	/* end transaction */
	return e.Commit()
}

func (e *DbExecutor) Single(name string, statement string) error {
	return e.Submit(statement)
}

func (e *DbExecutor) BulkInit(table string, columns ...string) error {
	return ErrCapNotSupported
}

func (e *DbExecutor) BulkAddRecord(args ...interface{}) error {
	return ErrCapNotSupported
}

func (e *DbExecutor) BulkFinish() error {
	return ErrCapNotSupported
}

func (e *DbExecutor) HasCapability(capability int) bool {
	/* return capability == CapBulkTransfer */
	return false
}

func (e *DbExecutor) GetDb() *sql.DB {
	return e.db
}

func (e *DbExecutor) GetTx() *sql.Tx {
	return e.tx
}

/* warning: closes the db connection that was passed to the constructor */
func (e *DbExecutor) Close() error {
	return e.err(e.db.Close())
}
