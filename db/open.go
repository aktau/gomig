package db

import (
	"fmt"
	. "github.com/aktau/gomig/db/common"
	"github.com/aktau/gomig/db/mysql"
	"github.com/aktau/gomig/db/postgres"
)

func OpenReader(driverName string, conf *Config) (ReadCloser, error) {
	switch driverName {
	case "mysql":
		return mysql.OpenReader(conf)
	}

	return nil, fmt.Errorf("db: OpenReader: unknown driver type: %v", driverName)
}

func OpenFileWriter(driverName string, filename string) (WriteCloser, error) {
	switch driverName {
	case "postgres":
		return postgres.NewPostgresFileWriter(filename)
	}

	return nil, fmt.Errorf("db: OpenFileWriter: unknown driver type: %v", driverName)
}

func OpenWriter(driverName string, conf *Config) (WriteCloser, error) {
	switch driverName {
	case "postgres":
		return postgres.NewPostgresWriter(conf)
	}

	return nil, fmt.Errorf("db: OpenWriter: unknown driver type: %v", driverName)
}
