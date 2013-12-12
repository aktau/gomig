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

	tables := r.FilteredTables(options.OnlyTables, options.ExcludeTables)

	if !options.SuppressDdl {
		createTables(tables, w)
	}
	if options.Truncate {
		truncateTables(tables, w)
	}
	if !options.SuppressData {
		if options.Merge {
			for _, table := range tables {
				w.MergeTable(table, r)
			}
		} else {
			writeData(tables, w)
		}
	}

	createIndices(tables, w)
	createConstraints(tables, w)

	return nil
}

func createTables(tables []*common.Table, w common.Writer) error {
	return nil
}

func truncateTables(tables []*common.Table, w common.Writer) error {
	return nil
}

func writeData(tables []*common.Table, w common.Writer) error {
	return nil
}

func createIndices(tables []*common.Table, w common.Writer) error {
	return nil
}

func createConstraints(tables []*common.Table, w common.Writer) error {
	return nil
}
