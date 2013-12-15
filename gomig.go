package main

import (
	"fmt"
	"github.com/aktau/gomig/db"
	"github.com/aktau/gomig/db/common"
	"log"
)

func main() {
	fmt.Printf("Welcome to gomig v.%v.%v\n\n", GOMIG_MAJ_VERSION, GOMIG_MIN_VERSION)

	conf, err := LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}
	log.Printf("gomig: successfully loaded config:\n\n%v\n\n", conf)

	/* open source */
	log.Println("gomig: connecting to source", conf.Mysql)
	reader, err := db.OpenReader("mysql", conf.Mysql)
	if err != nil {
		log.Println("gomig: error while creating reader,", err)
		return
	}
	defer reader.Close()
	log.Println("gomig: succesfully connected to source")

	/* open destination */
	log.Println("gomig: connecting to destination", conf.Destination)
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
	log.Println("gomig: succesfully connected to destination")

	log.Println("gomig: converting")
	err = Convert(reader, writer, conf)
	if err != nil {
		fmt.Println("gomig: could not complete conversion, error:", err)
	} else {
		log.Println("gomig: done")
	}
}
