package main

import (
	"fmt"
	"strings"
)

const (
	GOMIG_MAJ_VERSION = 0
	GOMIG_MIN_VERSION = 4
	GOMIG_MIC_VERSION = 4
)

type Backend struct {
	Source      string
	Destination string
}

func (b *Backend) String() string {
	return b.Source + " -> " + b.Destination
}

func description() string {
	backends := []Backend{Backend{"MySQL", "Postgres"}}
	stringized := make([]string, 0, len(backends))
	for _, backend := range backends {
		stringized = append(stringized, backend.String())
	}

	return fmt.Sprintf(
		"gomig v.%v.%v.%v, sync data between SQL data sources, supported backends: %v",
		GOMIG_MAJ_VERSION, GOMIG_MIN_VERSION, GOMIG_MIC_VERSION,
		strings.Join(stringized, ", "))
}

type VersionCommand struct{}

func (x *VersionCommand) Execute(args []string) error {
	fmt.Println(description())
	return nil
}

func init() {
	var x VersionCommand
	parser.AddCommand("version",
		"Print the version and supported backends",
		"Print the version and supported backends",
		&x)
}
