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
}

func NewDbExecutor(db *sql.DB) (*DbExecutor, error) {
	return &DbExecutor{db}, nil
}

func (e *DbExecutor) Multiple(name string, statements []string) []error {
	errors := make([]error, 0, len(statements))

	if DBEXEC_VERBOSE {
		log.Printf("DbExecutor: starting multiple statements: %v", name)
	}

	/* write out all statements, rollback in case of error */
	for _, stmt := range statements {
		if DBEXEC_VERBOSE {
			log.Println(stmt)
		}

		_, err := e.db.Exec(stmt)
		if err != nil {
			if DBEXEC_VERBOSE {
				log.Printf("DbExecutor: error while executing: %v", err)
			}

			errors = append(errors, err)
		}
	}

	return errors
}

func (e *DbExecutor) Transaction(name string, statements []string) error {
	if DBEXEC_VERBOSE {
		log.Printf("DbExecutor: starting transaction: %v", name)
	}

	/* start transaction */
	tx, err := e.db.Begin()
	if err != nil {
		return err
	}

	/* write out all statements, rollback in case of error */
	for _, stmt := range statements {
		if DBEXEC_VERBOSE {
			log.Println(stmt)
		}

		_, err := tx.Exec(stmt)
		if err != nil {
			rerr := tx.Rollback()
			if rerr != nil && DBEXEC_VERBOSE {
				log.Printf("DbExecutor: error while rolling back: %v", rerr)
			}
			return err
		}
	}

	/* end transaction */
	return tx.Commit()
}

func (e *DbExecutor) Single(name string, statement string) error {
	/* write comment */
	if DBEXEC_VERBOSE {
		log.Printf("DbExecutor: starting transaction: %v", name)
	}

	_, err := e.db.Exec(statement)
	return err
}

func (e *DbExecutor) HasCapability(capability int) bool {
	return false
}

func (e *DbExecutor) GetDb() *sql.DB {
	return e.db
}

/* warning: closes the db connection that was passed to the constructor */
func (e *DbExecutor) Close() error {
	return e.db.Close()
}
