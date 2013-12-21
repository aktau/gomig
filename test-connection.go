package main

import (
	"fmt"
	"github.com/aktau/gomig/db"
	"github.com/aktau/gomig/db/common"
	"launchpad.net/goyaml"
)

type TestConnectionCommand struct {
	/* config file */
	File string `short:"f" long:"file" description:"The path of the configuration file to use" default:"config.yml"`
}

func (x *TestConnectionCommand) Execute(args []string) error {
	verbosity := len(options.Verbose)
	common.DBEXEC_VERBOSE = false

	conf := LoadConfigOrDie(x.File)

	haveError := false

	fmt.Println("Testing connection to both source and destination db (if specified)")

	/* try connecting to the source */
	if verbosity > 0 {
		rawSrcParams, _ := goyaml.Marshal(conf.Mysql)
		srcParams := string(rawSrcParams)
		fmt.Printf("source:\n%v\n", IndentWith(srcParams, "  "))
	}
	fmt.Print("connecting...")
	reader, err := db.OpenReader("mysql", conf.Mysql)
	if err != nil {
		fmt.Printf("ERROR (%v)\n", err)
		haveError = true
	} else {
		fmt.Println("OK")
		defer reader.Close()
	}

	if verbosity > 0 {
		fmt.Println("")
	}

	/* try connecting to the destination */
	if verbosity > 0 {
		rawDstParams, _ := goyaml.Marshal(conf.Destination)
		dstParams := string(rawDstParams)
		fmt.Printf("destination:\n%v\n", IndentWith(dstParams, "  "))
	}
	fmt.Print("connecting...")
	if conf.Destination.File != "" {
		fmt.Println("IS A FILE")
	} else {
		writer, err := db.OpenWriter("postgres", conf.Destination.Postgres)
		if err != nil {
			fmt.Printf("ERROR (%v)\n", err)
			haveError = true
		} else {
			fmt.Println("OK")
			defer writer.Close()
		}
	}

	if haveError {
		return fmt.Errorf("could not connect to all databases")
	}

	return nil
}

func init() {
	var x TestConnectionCommand
	parser.AddCommand("test",
		"Test if a connection to the source and destination databases can be established",
		"Test if a connection to the source and destination databases can be established",
		&x)
}
