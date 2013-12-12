package main

import (
	"fmt"
	"github.com/aktau/gomig/db"
	"github.com/aktau/gomig/db/common"
	"log"
)

func main() {
	fmt.Printf("Welcome to gomig v.%v.%v\n", GOMIG_MAJ_VERSION, GOMIG_MIN_VERSION)

	conf, err := LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("successfully loaded config:\n\n%v\n\n", conf)

	/* open source */
	reader, err := db.OpenReader("mysql", conf.Mysql)
	if err != nil {
		log.Println("gomig: error while creating reader,", err)
		return
	}
	defer reader.Close()

	fmt.Println("ALL TABLES: ", reader.ListTables())

	/* open destination */
	var writer common.WriteCloser
	if conf.Destination.File != "" {
		writer, err = db.OpenFileWriter("postgres", conf.Destination.File)
	} else {
		writer, err = db.OpenWriter("postgres", conf.Destination.Postgres)
	}
	if err != nil {
		log.Println("gomig: error while creating writer,", err)
		return
	}
	defer writer.Close()

	Convert(reader, writer, conf)
}
