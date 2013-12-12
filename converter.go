package main

import (
	"github.com/aktau/gomig/db/common"
	"log"
)

var (
	VERBOSE = true
)

type tempViews struct {
	r     common.Reader
	views map[string]string
}

func createViews(r common.Reader, views map[string]string) *tempViews {
	t := &tempViews{r, views}

	/* create all views */
	t.Create()

	return t
}

func (t *tempViews) Create() {
	for name, body := range t.views {
		if VERBOSE {
			log.Printf("creating view '%v' with body\n\t%v\n", name, body)
		}

		err := t.r.CreateView(name, body)
		if err != nil {
			log.Println("converter: error while creating view", name, body, err)
		}
	}
}

func (t *tempViews) Erase() {
	for name, _ := range t.views {
		err := t.r.DropView(name)
		if err != nil {
			log.Println("converter: error while creating view", name, err)
		}
	}
}

func Convert(r common.ReadCloser, w common.WriteCloser, options *Config) error {
	if len(options.Views) > 0 {
		if VERBOSE {
			log.Println("converter: creating views...")
		}
		tempViews := createViews(r, options.Views)
		defer tempViews.Erase()
	}

	if !options.SuppressDdl {
		createTables(r, w, options)
	}
	if options.Truncate {
		truncateTables(r, w, options)
	}
	if !options.SuppressData {
		if options.Merge {
			mergeData(r, w, options)
		} else {
			writeData(r, w, options)
		}
	}

	createIndices(r, w, options)
	createConstraints(r, w, options)

	return nil
}

func createTables(r common.ReadCloser, w common.WriteCloser, options *Config) error {
	return nil
}

func truncateTables(r common.ReadCloser, w common.WriteCloser, options *Config) error {
	return nil
}

/* the only one I need for the moment */
func mergeData(r common.ReadCloser, w common.WriteCloser, options *Config) error {
	tables := activeTables(r.ListTables(), options)

	for _, table := range tables {
		w.MergeTable(&table, r)
	}

	return nil
}

func activeTables(alltables []common.Table, options *Config) []common.Table {
	tables := make([]common.Table, 0, 4)
	for _, table := range alltables {
		_, incl := options.OnlyTables[table.Name]
		_, excl := options.ExcludeTables[table.Name]

		if incl && !excl {
			tables = append(tables, table)
		}
	}

	return tables
}

func writeData(r common.ReadCloser, w common.WriteCloser, options *Config) error {
	return nil
}

func createIndices(r common.ReadCloser, w common.WriteCloser, options *Config) error {
	return nil
}

func createConstraints(r common.ReadCloser, w common.WriteCloser, options *Config) error {
	return nil
}
