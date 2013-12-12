package common

import (
	"bufio"
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
	_, err := e.w.WriteString(fmt.Sprintf("-- %v", name))
	if err != nil {
		return err
	}

	/* start transaction */
	e.w.WriteString(SBegin)

	/* write out all statements */
	for _, stmt := range statements {
		_, err := e.w.WriteString(stmt)
		if err != nil {
			return err
		}
	}

	/* end transaction */
	e.w.WriteString(SBegin)

	return nil
}

func (e *FileExecutor) Single(name string, statement string) error {
	/* write comment */
	_, err := e.w.WriteString(fmt.Sprintf("-- %v", name))
	if err != nil {
		return err
	}

	_, err = e.w.WriteString(statement)
	return err
}

func (e *FileExecutor) HasCapability(capability int) bool {
	return false
}

func (e *FileExecutor) Close() error {
	err := e.w.Flush()
	if err != nil {
		defer e.f.Close()
		return err
	}

	return e.f.Close()
}
