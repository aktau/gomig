package mysql

import (
	"database/sql"
	"fmt"
	. "github.com/aktau/gomig/db/common"
	"log"
	"strings"
)

var (
	READER_VERBOSE = false
)

var (
	mysqlInit = []string{
		"SET collation_connection = utf8_general_ci",
		"SET NAMES utf8",
	}
)

type MysqlReader struct {
	*sql.DB
}

func OpenReader(conf *Config) (*MysqlReader, error) {
	db, err := openDB(conf)
	if err != nil {
		return nil, err
	}

	log.Printf("mysql/openreader: initializing")
	for _, stmt := range mysqlInit {
		log.Printf("%v", stmt)
		if _, err := db.Exec(stmt); err != nil {
			defer db.Close()
			return nil, err
		}
	}

	return &MysqlReader{db}, nil
}

func (r *MysqlReader) TableNames() []string {
	rows, err := r.Query("SHOW TABLES;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tables := make([]string, 0, 8)

	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}

		tables = append(tables, name)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return tables
}

func (r *MysqlReader) Tables() []*Table {
	return r.FilteredTables(nil, nil)
}

func (r *MysqlReader) FilteredTables(incl, excl map[string]bool) []*Table {
	tableNames := r.TableNames()
	filteredTableNames := FilterInclExcl(tableNames, incl, excl)
	tables := make([]*Table, 0, len(filteredTableNames))

	if READER_VERBOSE {
		log.Printf("mysql: all tables = %v, filtered = %v\n", tableNames, filteredTableNames)
	}

	for _, tableName := range filteredTableNames {
		/* query table information */
		columns, err := r.columns(tableName)
		if err != nil {
			log.Println("mysql: could not fetch columns of table", tableName, "error:", err)
		}

		/* create table struct */
		table := &Table{Name: tableName, DbType: "mysql", Columns: columns}

		tables = append(tables, table)
	}

	return tables
}

type rawCol struct {
	name    string
	rawtype string
	null    string
	key     string
	defval  sql.NullString
	extra   string
}

func (r *MysqlReader) columns(table string) ([]*Column, error) {
	rows, err := r.Query(fmt.Sprintf("EXPLAIN `%v`;", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]*Column, 0, 8)

	var rc rawCol
	for rows.Next() {
		err = rows.Scan(&rc.name, &rc.rawtype, &rc.null, &rc.key, &rc.defval, &rc.extra)
		if err != nil {
			return nil, err
		}

		col, err := r.processCol(table, &rc)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cols, nil
}

func (r *MysqlReader) processCol(table string, rc *rawCol) (*Column, error) {
	t := rc.rawtype
	length := 255

	return &Column{
		TableName:    table,
		Name:         rc.name,
		Type:         MysqlToGenericType(t),
		RawType:      t,
		Length:       length,
		Null:         rc.null == "YES" || strings.HasPrefix(t, "enum") || t == "date" || t == "datetime" || t == "timestamp",
		PrimaryKey:   rc.key == "PRI",
		AutoIncr:     rc.extra == "auto_increment",
		Default:      rc.defval,
		NeedsQuoting: strings.Contains(t, "text") || strings.Contains(t, "varchar"),
	}, nil
}

/* caller is responsible for cleaning up the sql.Rows object */
func (r *MysqlReader) Read(table *Table) (*sql.Rows, error) {
	rows, err := r.Query(fmt.Sprintf("SELECT * FROM %v;", table.Name))
	if err != nil {
		return nil, err
	}

	/* vals := make([]interface{}, len(src.Columns)) */
	/* for rows.Next() { */
	/* 	err = rows.Scan(vals...) */
	/* 	if err != nil { */
	/* 		panic(err) */
	/* 	} */

	/* 	tables = append(tables, name) */
	/* } */

	return rows, nil
}

func (r *MysqlReader) CreateView(name string, body string) error {
	stmt := fmt.Sprintf("CREATE VIEW %v AS %v;", name, body)

	_, err := r.Exec(stmt)
	return err
}

func (r *MysqlReader) DropView(name string) error {
	stmt := fmt.Sprintf("DROP VIEW %v;", name)

	_, err := r.Exec(stmt)
	return err
}

/* can't use temporary tables, as they don't appear in SHOW TABLES output */
func (r *MysqlReader) CreateProjection(name string, body string, engine string, pk []string, uks [][]string) error {
	var createPk string
	if len(pk) > 0 {
		createPk = " ( " + "PRIMARY KEY (" + strings.Join(pk, ", ") + ")" + " )"
	} else {
		createPk = ""
	}

	var engineSQL string
	if engine != "" {
		engineSQL = " ENGINE=" + strings.ToUpper(engine)
	}

	collation := " CHARACTER SET utf8 COLLATE utf8_general_ci"

	stmt := fmt.Sprintf("CREATE TABLE %v%v%v%v AS (\n%v\n);",
		name, createPk, engineSQL, collation, body)

	if READER_VERBOSE {
		log.Printf("mysql: creating projection:\n%v\n", stmt)
	}
	_, err := r.Exec(stmt)
	return err
}

func (r *MysqlReader) DropProjection(name string) error {
	stmt := fmt.Sprintf("DROP TABLE %v;", name)

	_, err := r.Exec(stmt)
	return err
}
