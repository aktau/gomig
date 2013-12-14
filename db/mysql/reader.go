package mysql

import (
	"database/sql"
	"fmt"
	. "github.com/aktau/gomig/db/common"
	"log"
	"strings"
)

type MysqlReader struct {
	*sql.DB
}

func OpenReader(conf *Config) (*MysqlReader, error) {
	db, err := openDB(conf)
	if err != nil {
		return nil, err
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

	for _, tableName := range filteredTableNames {
		/* query table information */
		columns, err := r.columns(tableName)
		if err != nil {
			log.Println("MysqlReader: could not fetch columns of table", tableName, "error:", err)
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
		Type:         r.normalizeType(t),
		RawType:      t,
		Length:       length,
		Null:         rc.null == "YES" || strings.HasPrefix(t, "enum") || t == "date" || t == "datetime" || t == "timestamp",
		PrimaryKey:   rc.key == "PRI",
		AutoIncr:     rc.extra == "auto_increment",
		Default:      rc.defval,
		NeedsQuoting: strings.Contains(t, "text") || strings.Contains(t, "varchar"),
	}, nil
}

func (r *MysqlReader) normalizeType(rt string) string {
	switch {
	case strings.Contains(rt, "char"), strings.Contains(rt, "text"):
		return "text"
	case rt == "bit(1)", rt == "tinyint(1)", rt == "tinyint(1) unsigned":
		return "boolean"
	default:
		return "integer"
	}
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
