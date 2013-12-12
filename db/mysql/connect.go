package mysql

import (
	"database/sql"
	"fmt"
	. "github.com/aktau/gomig/db/common"
	_ "github.com/go-sql-driver/mysql"
)

func openDB(conf *Config) (*sql.DB, error) {
	protocol := "unix"
	address := conf.Socket
	if address == "" {
		protocol = "tcp"
		port := 3306
		if conf.Port != 0 {
			port = conf.Port
		}
		address = fmt.Sprintf("%v:%v", conf.Hostname, port)
	}

	/* root:pw@unix(/tmp/mysql.sock)/myDatabase?loc=Local */
	uri := fmt.Sprintf("%v:%v@%v(%v)/%v", conf.Username, conf.Password,
		protocol, address, conf.Database)
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	/* try to ping, let's fail fast */
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
