package common

import (
	"io"
)

type Writer interface {
	/*
		CreateTable(t *Table) error
		Truncate(t *Table) error
	*/

	/* merge the contents of table */
	MergeTable(src *Table, dstName, extraDstCond string, r Reader) error

	/* (over)write the contents of table */
	/* WriteTable(t *Table) error */

	/*
		CreateIndices(t *Table) error
		CreateConstraints(t *Table) error
	*/
}

type WriteCloser interface {
	io.Closer
	Writer
}
