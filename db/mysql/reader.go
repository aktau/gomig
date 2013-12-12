package mysql

import (
	"database/sql"
	"fmt"
	. "github.com/aktau/gomig/db/common"
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
		/* columns := fetchColumns() */

		/* create table struct */
		table := &Table{Name: tableName}

		tables = append(tables, table)
	}

	return tables
}

/* caller is responsible for cleaning up the sql.Rows object */
func (r *MysqlReader) Read(table *Table) (*sql.Rows, error) {
	rows, err := r.Query(fmt.Sprint("SELECT * FROM %v;", table.Name))
	if err != nil {
		return nil, err
	}

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
