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

func (r *MysqlReader) ListTables() []Table {
	rows, err := r.Query("SHOW TABLES;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tables := make([]Table, 0, 8)

	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}

		tables = append(tables, Table{Name: name})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return tables
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
