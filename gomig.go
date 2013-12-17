package main

import (
	"fmt"
	"github.com/aktau/gomig/db"
	"github.com/aktau/gomig/db/common"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strings"
)

const (
	DEFAULT_CONFIG_PATH = "config.yml"
	SAMPLE_CONFIG_PATH  = "config.sample.yml"
)

type Options struct {
	/* verbosity level */
	Verbose []bool `short:"v" long:"verbose" description:"verbose output"`

	/* config file */
	File string `short:"f" long:"file" description:"the location of the configuration file" optional:"yes" default:"config.yml"`
}

type Backend struct {
	Source      string
	Destination string
}

func (b *Backend) String() string {
	return b.Source + " -> " + b.Destination
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

func description() {
	backends := []Backend{Backend{"MySQL", "Postgres"}}
	stringized := make([]string, 0, len(backends))
	for _, backend := range backends {
		stringized = append(stringized, "- "+backend.String())
	}

	fmt.Printf(
		"Welcome to gomig v.%v.%v.%v, sync data between SQL data sources, supported backends:\n%v\n\n",
		GOMIG_MAJ_VERSION, GOMIG_MIN_VERSION, GOMIG_MIC_VERSION,
		strings.Join(stringized, "\n"))
}

func main() {
	description()

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	verbosity := len(options.Verbose)

	conf, err := LoadConfig(options.File, DEFAULT_CONFIG_PATH, SAMPLE_CONFIG_PATH)
	if err != nil {
		fmt.Print("error while loading config file, ", err)
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
		log.Println("gomig: error while creating reader,", err)
		return
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
		log.Println("gomig: error while creating writer,", err)
		return
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
}
