package postgres

import (
	"database/sql"
	"fmt"
	. "github.com/aktau/gomig/db/common"
	_ "github.com/lib/pq"
)

func openDB(conf *Config) (*sql.DB, error) {
	address := conf.Socket
	if address == "" {
		port := 3306
		if conf.Port != 0 {
			port = conf.Port
		}
		address = fmt.Sprintf("%v:%v", conf.Hostname, port)
	}

	uri := fmt.Sprintf("%v:%v@%v/%v", conf.Username, conf.Password,
		address, conf.Database)
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	/* try to ping, let's fail fast */
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
