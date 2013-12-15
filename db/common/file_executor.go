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
}

func NewFileExecutor(filename string) (*FileExecutor, error) {
	fo, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &FileExecutor{fo, bufio.NewWriter(fo)}, nil
}

func (e *FileExecutor) Transaction(name string, statements []string) error {
	/* write comment */
	_, err := e.w.WriteString(fmt.Sprintf("-- %v\n", name))
	if err != nil {
		return err
	}

	/* start transaction */
	e.w.WriteString(SBegin + ";\n\n")

	/* write out all statements */
	for _, stmt := range statements {
		_, err := e.w.WriteString(stmt + "\n")
		if err != nil {
			return err
		}
	}

	/* end transaction */
	e.w.WriteString(SCommit + ";\n\n")

	return nil
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
		_, err := e.w.WriteString(stmt + "\n")
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
