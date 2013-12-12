package postgres

import (
	"database/sql"
	. "github.com/aktau/gomig/db/common"
)

type PostgresReader struct {
	*sql.DB
}

func OpenReader(conf *Config) (*PostgresReader, error) {
	db, err := openDB(conf)
	if err != nil {
		return nil, err
	}

	return &PostgresReader{db}, nil
}
