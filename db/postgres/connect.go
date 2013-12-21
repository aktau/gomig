package postgres

import (
	"database/sql"
	"fmt"
	. "github.com/aktau/gomig/db/common"
	_ "github.com/lib/pq"
	"strings"
)

func openDB(conf *Config) (*sql.DB, error) {
	params := make([]string, 0, 4)

	if conf.Username != "" {
		params = append(params, "user="+conf.Username)
	}
	if conf.Password != "" {
		params = append(params, fmt.Sprintf("password='%v'", conf.Password))
	}
	if conf.Socket != "" {
		params = append(params, fmt.Sprintf("host='%v'", conf.Socket))
		params = append(params, "sslmode=disable")
	} else {
		port := 3306
		if conf.Port != 0 {
			port = conf.Port
		}

		params = append(params, fmt.Sprintf("host='%v'", conf.Hostname))
		params = append(params, fmt.Sprintf("port=%v", port))
	}
	if conf.Database != "" {
		params = append(params, "dbname="+conf.Database)
	}

	uri := strings.Join(params, " ")
	db, err := sql.Open("postgres", uri)
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
