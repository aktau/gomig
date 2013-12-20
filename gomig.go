package main

import (
	"github.com/jessevdk/go-flags"
	"os"
)

type Options struct {
	/* verbosity level */
	Verbose []bool `short:"v" long:"verbose" description:"verbose output"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
