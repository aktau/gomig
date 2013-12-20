package main

import (
	"github.com/jessevdk/go-flags"
	"os"
)

const (
	DEFAULT_CONFIG_PATH = "config.yml"
	SAMPLE_CONFIG_PATH  = "config.sample.yml"
)

type Options struct {
	/* verbosity level */
	Verbose []bool `short:"v" long:"verbose" description:"verbose output"`
	Version bool   `long:"version" description:"shows the program version and available backends"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
