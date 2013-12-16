package common

import (
	"database/sql"
	"log"
)

var (
	DBEXEC_VERBOSE = true
)

type DbExecutor struct {
	db *sql.DB
	tx *sql.Tx
}

func NewDbExecutor(db *sql.DB) (*DbExecutor, error) {
	return &DbExecutor{db, nil}, nil
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
		return nil
	}

	e.tx = tx
	return nil
}

func (e *DbExecutor) Commit() error {
	if e.tx == nil {
		return ErrNoTxInProgress
	}

	/* end transaction */
	tx := e.tx
	e.tx = nil
	return tx.Commit()
}

func (e *DbExecutor) submitSimple(stmt string) error {
	_, err := e.db.Exec(stmt)
	return err
}

func (e *DbExecutor) submitTransactional(stmt string) error {
	tx := e.tx
	_, err := tx.Exec(stmt)
	if err != nil {
		rerr := tx.Rollback()
		if rerr != nil && DBEXEC_VERBOSE {
			log.Printf("DbExecutor: error while rolling back: %v", rerr)
		}
	}
	return err
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

func (e *DbExecutor) BulkInit(table string) error {
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

/* warning: closes the db connection that was passed to the constructor */
func (e *DbExecutor) Close() error {
	return e.db.Close()
}
