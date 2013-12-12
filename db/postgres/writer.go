package postgres

import (
	"fmt"
	. "github.com/aktau/gomig/db/common"
)

type genericPostgresWriter struct {
	e Executor
}

/* how to do an UPSERT/MERGE in PostgreSQL
 * http://stackoverflow.com/questions/17267417/how-do-i-do-an-upsert-merge-insert-on-duplicate-update-in-postgresq */
func (w *genericPostgresWriter) MergeTable(table *Table, r Reader) error {
	stmts := make([]string, 0, 5)

	/* create temporary table */
	stmts = append(stmts, "CREATE TEMPORARY TABLE newvals(id integer, somedata text);")

	/* bulk insert values */
	// TODO: loop over reader

	/* lock the target table */
	stmts = append(stmts, fmt.Sprint("LOCK TABLE %v IN EXCLUSIVE MODE;", table.Name))

	/* UPDATE from temp table to target table based on PK */
	stmts = append(stmts, fmt.Sprint(
		`UPDATE testtable
		 SET somedata = newvals.somedata
		 FROM newvals
		 WHERE newvals.id = testtable.id;`))

	/* INSERT from temp table to target table based on PK */
	stmts = append(stmts, fmt.Sprint(
		`INSERT INTO testtable
		 SELECT newvals.id, newvals.somedata
		 FROM newvals
		 LEFT OUTER JOIN testtable ON (testtable.id = newvals.id)
		 WHERE testtable.id IS NULL;`))

	err := w.e.Transaction("merge some stuff", []string{"SELECT...", "INSERT INTO..."})
	return err
}

func (w *genericPostgresWriter) Close() error {
	return w.e.Close()
}

type PostgresWriter struct {
	genericPostgresWriter
}

func NewPostgresWriter(conf *Config) (*PostgresWriter, error) {
	db, err := openDB(conf)
	if err != nil {
		return nil, err
	}

	executor, err := NewDbExecutor(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresWriter{genericPostgresWriter{executor}}, nil
}

type PostgresFileWriter struct {
	genericPostgresWriter
}

func NewPostgresFileWriter(filename string) (*PostgresFileWriter, error) {
	executor, err := NewFileExecutor(filename)
	if err != nil {
		return nil, err
	}
	return &PostgresFileWriter{genericPostgresWriter{executor}}, err
}
