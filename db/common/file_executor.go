package common

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
)

const (
	SBegin    = "BEGIN"
	SCommit   = "COMMIT"
	SRollback = "ROLLBACK"
)

type FileExecutor struct {
	f *os.File
	w *bufio.Writer

	txInProgress bool
}

func NewFileExecutor(filename string) (*FileExecutor, error) {
	fo, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &FileExecutor{fo, bufio.NewWriter(fo), false}, nil
}

func (e *FileExecutor) Begin(name string) error {
	if e.txInProgress {
		return ErrTxInProgress
	}
	e.txInProgress = true

	/* write comment */
	_, err := e.w.WriteString(fmt.Sprintf("-- %v\n", name))
	if err != nil {
		return err
	}

	/* start transaction */
	_, err = e.w.WriteString(SBegin + ";\n\n")
	return err
}

func (e *FileExecutor) Commit() error {
	if !e.txInProgress {
		return ErrNoTxInProgress
	}
	e.txInProgress = false

	/* end transaction */
	_, err := e.w.WriteString(SCommit + ";\n\n")
	return err
}

func (e *FileExecutor) Submit(stmt string) error {
	_, err := e.w.WriteString(stmt + "\n")
	return err
}

func (e *FileExecutor) Transaction(name string, statements []string) error {
	err := e.Begin(name)
	if err != nil {
		return err
	}

	/* write out all statements */
	for _, stmt := range statements {
		err := e.Submit(stmt)
		if err != nil {
			return err
		}
	}

	return e.Commit()
}

func (e *FileExecutor) Multiple(name string, statements []string) []error {
	errors := make([]error, 0, len(statements))

	/* write comment */
	_, err := e.w.WriteString(fmt.Sprintf("-- %v\n", name))
	if err != nil {
		errors = append(errors, err)
	}

	/* write out all statements, rollback in case of error */
	for _, stmt := range statements {
		err := e.Submit(stmt)
		if err != nil {
			errors = append(errors, err)
		}
	}

	_, err = e.w.WriteString("\n")
	if err != nil {
		errors = append(errors, err)
	}

	return errors
}

func (e *FileExecutor) Single(name string, statement string) error {
	/* write comment */
	_, err := e.w.WriteString(fmt.Sprintf("-- %v\n", name))
	if err != nil {
		return err
	}

	_, err = e.w.WriteString(statement)
	return err
}

func (e *FileExecutor) BulkInit(table string) error {
	return ErrCapNotSupported
}

func (e *FileExecutor) BulkAddRecord(args ...interface{}) error {
	return ErrCapNotSupported
}

func (e *FileExecutor) BulkFinish() error {
	return ErrCapNotSupported
}

func (e *FileExecutor) HasCapability(capability int) bool {
	return false
}

func (e *FileExecutor) GetDb() *sql.DB {
	return nil
}

func (e *FileExecutor) Close() error {
	err := e.w.Flush()
	if err != nil {
		defer e.f.Close()
		return err
	}

	return e.f.Close()
}
