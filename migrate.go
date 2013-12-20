package main

import (
	"fmt"
	"github.com/aktau/gomig/db"
	"github.com/aktau/gomig/db/common"
	"log"
	"os"
)

const (
	PATH_CONFIG_DEFAULT = "config.yml"
)

type MigrateCommand struct {
	/* config file */
	File string `short:"f" long:"file" description:"The path of the configuration file to use" default:"config.yml"`
}

func (x *MigrateCommand) Execute(args []string) error {
	verbosity := len(options.Verbose)

	if x.File == "" {
		x.File = PATH_CONFIG_DEFAULT
	}

	conf, err := LoadConfig(x.File)
	if err != nil {
		fmt.Printf("error while loading config file: '%v'\n", err)
		fmt.Println("to generate a sample config file use the generate-config command")
		os.Exit(1)
	}

	if verbosity > 2 {
		fmt.Println("config:", conf)
	}

	/* open source */
	if verbosity > 0 {
		log.Println("gomig: connecting to source", conf.Mysql)
	}

	reader, err := db.OpenReader("mysql", conf.Mysql)
	if err != nil {
		return fmt.Errorf("gomig: error while creating reader, %v", err)
	}
	defer reader.Close()

	if verbosity > 0 {
		log.Println("gomig: succesfully connected to source")
	}

	/* open destination */
	if verbosity > 0 {
		log.Println("gomig: connecting to destination", conf.Destination)
	}

	var writer common.WriteCloser
	if conf.Destination.File != "" {
		writer, err = db.OpenFileWriter("postgres", conf.Destination.File)
	} else {
		writer, err = db.OpenWriter("postgres", conf.Destination.Postgres)
	}
	if err != nil {
		return fmt.Errorf("gomig: error while creating writer, err", err)
	}
	defer writer.Close()

	if verbosity > 0 {
		log.Println("gomig: succesfully connected to destination")
	}

	log.Println("gomig: converting")
	err = Convert(reader, writer, conf, verbosity)
	if err != nil {
		fmt.Println("gomig: could not complete conversion, error:", err)
	} else {
		log.Println("gomig: done")
	}
	return nil
}

func init() {
	var cmd MigrateCommand
	parser.AddCommand("migrate",
		"Migrate data from a source database to a destination file/database",
		"Migrate data from a source database to a destination file/database",
		&cmd)
}
