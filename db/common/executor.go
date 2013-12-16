package common

import (
	"database/sql"
	"errors"
	"io"
)

const (
	CapBulkTransfer = iota
)

var (
	ErrTxInProgress    = errors.New("another transaction is already in progress")
	ErrNoTxInProgress  = errors.New("no transaction is in progress")
	ErrCapNotSupported = errors.New("capbility not supported")
)

type Executor interface {
	io.Closer

	/* begin a transaction, it's an error to begin a transaction
	 * while another is already in progress */
	Begin(name string) error
	Commit() error

	/* if the statement provokes an error, will automatically rollback,
	 * after which the transaction is no longer in progress */
	Submit(stmt string) error

	/* bulk statements for copying large amounts of data, the underlying
	 * implementation will try to use the most efficient way of achieving this,
	 * for example postgres' COPY FROM semantics. */
	BulkInit(table string, columns ...string) error
	BulkAddRecord(args ...interface{}) error
	BulkFinish() error

	/* submit a transaction in one go */
	Transaction(name string, statements []string) error

	/* submit multiple statements in one go (without a transaction) */
	Multiple(name string, statements []string) []error

	/* submit a single statement */
	Single(name string, statement string) error

	/* for e.g. PostgreSQL COPY support */
	HasCapability(capability int) bool

	/* will return the underlying sql.DB if it exists, otherwise nil */
	GetDb() *sql.DB

	/* will return the underlying sql.Tx if it exists, otherwise nil */
	GetTx() *sql.Tx
}
