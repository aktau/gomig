package postgres

import (
	"database/sql"
	"fmt"
	. "github.com/aktau/gomig/db/common"
	"log"
	"strings"
)

var PG_W_VERBOSE = true

var (
	postgresInit = []string{
		"SET client_encoding = 'UTF8';",
		"SET standard_conforming_strings = off;",
		"SET check_function_bodies = false;",
		"SET client_min_messages = warning;",
	}
)

const (
	explainQuery = `
SELECT col.column_name AS field,
       CASE
        WHEN col.character_maximum_length IS NOT NULL THEN col.data_type || '(' || col.character_maximum_length || ')'
        ELSE col.data_type
       END AS type,
       col.is_nullable AS null,
       CASE
        WHEN tc.constraint_type = 'PRIMARY KEY' THEN 'PRI'
        ELSE ''
       END AS key,
       '' AS default,
       '' AS extra
       --kcu.constraint_name AS constraint_name
       --kcu.*,
       --tc.*
FROM   information_schema.columns col
LEFT JOIN   information_schema.key_column_usage kcu ON (kcu.table_name = col.table_name AND kcu.column_name = col.column_name)
LEFT JOIN   information_schema.table_constraints AS tc ON (kcu.constraint_name = tc.constraint_name)
WHERE  col.table_name = '%v'
ORDER BY col.ordinal_position;`
)

type genericPostgresWriter struct {
	e               Executor
	insertBulkLimit int
}

func (w *genericPostgresWriter) bulkTransfer(src *Table, dstName string, rows *sql.Rows) (err error) {
	ex := w.e

	colnames := make([]string, 0, len(src.Columns))
	for _, col := range src.Columns {
		colnames = append(colnames, col.Name)
	}

	if err = ex.BulkInit(dstName, colnames...); err != nil {
		return
	}
	defer func() {
		berr := ex.BulkFinish()
		if err == nil {
			/* if there was no earlier error, set the one from BulkFinish */
			err = berr
		}
	}()

	/* create a slice with the right types to extract into, and let the SQL
	 * driver take care of the conversion */
	vals := NewTypedSlice(src)

	for rows.Next() {
		if err = rows.Scan(vals...); err != nil {
			return fmt.Errorf("postgres: error while reading from source:", err)
		}

		if err = ex.BulkAddRecord(vals...); err != nil {
			return fmt.Errorf("postgres: error during bulk insert:", err)
		}
	}

	return
}

func (w *genericPostgresWriter) normalTransfer(src *Table, dstName string, rows *sql.Rows, linesPerStatement int) error {
	/* an alternate way to do this, with type assertions
	 * but possibly less accurately: http://go-database-sql.org/varcols.html */
	pointers := make([]interface{}, len(src.Columns))
	containers := make([]sql.RawBytes, len(src.Columns))
	for i, _ := range pointers {
		pointers[i] = &containers[i]
	}
	stringrep := make([]string, 0, len(src.Columns))
	insertLines := make([]string, 0, 32)
	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			log.Println("postgres: error while reading from source:", err)
			return err
		}

		for idx, val := range containers {
			str, err := RawToPostgres(val, src.Columns[idx].Type)
			if err != nil {
				return err
			}
			stringrep = append(stringrep, str)
		}

		insertLines = append(insertLines, "("+strings.Join(stringrep, ",")+")")
		stringrep = stringrep[:0]

		if len(insertLines) > w.insertBulkLimit {
			err = w.e.Submit(fmt.Sprintf("INSERT INTO %v VALUES\n\t%v;\n",
				dstName, strings.Join(insertLines, ",\n\t")))
			if err != nil {
				return err
			}

			insertLines = insertLines[:0]
		}
	}

	if len(insertLines) > 0 {
		err := w.e.Submit(fmt.Sprintf("INSERT INTO %v VALUES\n\t%v;\n",
			dstName, strings.Join(insertLines, ",\n\t")))
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *genericPostgresWriter) transferTable(src *Table, dstName string, r Reader) error {
	/* bulk insert values */
	rows, err := r.Read(src)
	if err != nil {
		return err
	}
	defer rows.Close()

	if PG_W_VERBOSE {
		log.Print("postgres: query done, scanning rows...")
	}

	if w.e.HasCapability(CapBulkTransfer) {
		if PG_W_VERBOSE {
			log.Print("postgres: bulk capability detected, performing bulk transfer...")
		}

		err = w.bulkTransfer(src, dstName, rows)
	} else {
		if PG_W_VERBOSE {
			log.Print("postgres: no bulk capability detected, performing normal transfer...")
		}

		err = w.normalTransfer(src, dstName, rows, w.insertBulkLimit)
	}
	if err != nil {
		return err
	}

	return rows.Err()
}

/* how to do an UPSERT/MERGE in PostgreSQL
 * http://stackoverflow.com/questions/17267417/how-do-i-do-an-upsert-merge-insert-on-duplicate-update-in-postgresq */
func (w *genericPostgresWriter) MergeTable(src *Table, dstName string, r Reader) error {
	tmpName := "gomig_tmp"

	err := w.e.Begin(
		fmt.Sprintf("merge table %v into table %v", src.Name, dstName))
	if err != nil {
		return err
	}

	/* create temporary table */
	err = w.e.Submit(fmt.Sprintf("CREATE TEMPORARY TABLE %v (\n\t%v\n)\nON COMMIT DROP;\n", tmpName, ColumnsSql(src)))
	if err != nil {
		return err
	}

	if PG_W_VERBOSE {
		log.Println("postgres: preparing to read values from source db")
	}

	err = w.transferTable(src, tmpName, r)
	if err != nil {
		return err
	}

	if PG_W_VERBOSE {
		log.Print("postgres: rowscan done, creating merge statements")
	}

	/* analyze the temp table, for performance */
	err = w.e.Submit(fmt.Sprintf("ANALYZE %v;\n", tmpName))
	if err != nil {
		return err
	}

	/* lock the target table */
	err = w.e.Submit(fmt.Sprintf("LOCK TABLE %v IN EXCLUSIVE MODE;", dstName))
	if err != nil {
		return err
	}

	colnames := make([]string, 0, len(src.Columns))
	srccol := make([]string, 0, len(src.Columns))
	pkWhere := make([]string, 0, len(src.Columns))
	pkIsNull := make([]string, 0, len(src.Columns))
	colassign := make([]string, 0, len(src.Columns))
	for _, col := range src.Columns {
		colnames = append(colnames, col.Name)
		srccol = append(srccol, "src."+col.Name)
		if col.PrimaryKey {
			pkWhere = append(pkWhere, fmt.Sprintf("dst.%[1]v = src.%[1]v", col.Name))
			pkIsNull = append(pkIsNull, fmt.Sprintf("dst.%[1]v IS NULL", col.Name))
		} else {
			colassign = append(colassign, fmt.Sprintf("%[1]v = src.%[1]v", col.Name))
		}
	}
	pkWherePart := strings.Join(pkWhere, "\nAND    ")
	pkIsNullPart := strings.Join(pkIsNull, "\nAND    ")
	srccolPart := strings.Join(srccol, ",\n       ")

	/* UPDATE from temp table to target table based on PK */
	err = w.e.Submit(fmt.Sprintf(`
UPDATE %v AS dst
SET    %v
FROM   %v AS src
WHERE  %v;`, dstName, strings.Join(colassign, ",\n       "), tmpName, pkWherePart))
	if err != nil {
		return err
	}

	/* INSERT from temp table to target table based on PK */
	err = w.e.Submit(fmt.Sprintf(`
INSERT INTO %[1]v (%[3]v)
SELECT %[4]v
FROM   %[2]v AS src
LEFT OUTER JOIN %[1]v AS dst ON (
	   %[5]v
)
WHERE  %[6]v;
`, dstName, tmpName, strings.Join(colnames, ", "), srccolPart, pkWherePart, pkIsNullPart))
	if err != nil {
		return err
	}

	if PG_W_VERBOSE {
		log.Print("postgres: statements completed, executing transaction")
	}

	err = w.e.Commit()
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

	executor, err := NewPgDbExecutor(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	errors := executor.Multiple("initializing DB connection (WARNING: connection pooling might mess with this)", postgresInit)
	if len(errors) > 0 {
		executor.Close()
		for _, err := range errors {
			log.Println("postgres error:", err)
		}
		return nil, errors[0]
	}

	return &PostgresWriter{genericPostgresWriter{executor, 64}}, nil
}

type PostgresFileWriter struct {
	genericPostgresWriter
}

func NewPostgresFileWriter(filename string) (*PostgresFileWriter, error) {
	executor, err := NewFileExecutor(filename)
	if err != nil {
		return nil, err
	}

	errors := executor.Multiple("initializing DB connection", postgresInit)
	if len(errors) > 0 {
		executor.Close()
		for _, err := range errors {
			log.Println("postgres error:", err)
		}
		return nil, errors[0]
	}

	return &PostgresFileWriter{genericPostgresWriter{executor, 256}}, err
}

func ColumnsSql(table *Table) string {
	colSql := make([]string, 0, len(table.Columns))

	for _, col := range table.Columns {
		colSql = append(colSql, fmt.Sprintf("%v %v", col.Name, GenericToPostgresType(col.Type)))
	}

	pkCols := make([]string, 0, len(table.Columns))
	for _, col := range table.Columns {
		if col.PrimaryKey {
			pkCols = append(pkCols, col.Name)
		}
	}

	/* add the primary key */
	colSql = append(colSql, fmt.Sprintf("PRIMARY KEY (%v)",
		strings.Join(pkCols, ", ")))

	return strings.Join(colSql, ",\n\t")
}
