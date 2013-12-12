package common

import (
	"io"
)

const (
	CapCopy = iota
)

type Executor interface {
	io.Closer

	Transaction(name string, statements []string) error
	Single(name string, statement string) error

	/* for e.g. PostgreSQL COPY support */
	HasCapability(capability int) bool
}
