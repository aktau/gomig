package main

import (
	"github.com/aktau/gomig/db/common"
	"log"
)

var (
	VERBOSE = true
)

type tempEntities struct {
	r           common.Reader
	views       map[string]string
	projections map[string]ProjectionConfig
}

func createTempEntities(r common.Reader, views map[string]string, projections map[string]ProjectionConfig) *tempEntities {
	t := &tempEntities{r, views, projections}
	t.Create()
	return t
}

func (t *tempEntities) Create() {
	for name, body := range t.views {
		if VERBOSE {
			log.Printf("converter: creating view '%v'\n", name)
		}

		err := t.r.CreateView(name, body)
		if err != nil {
			log.Println("converter: error while creating view", name, body, err)
		}
	}

	for name, proj := range t.projections {
		if VERBOSE {
			log.Printf("converter: creating projection '%v'\n", name)
		}

		err := t.r.CreateProjection(name, proj.Body, proj.Engine, proj.Pk, nil)
		if err != nil {
			log.Println("converter: error while creating projection", name, proj.Body, proj.Pk, err)
		}
	}
}

func (t *tempEntities) Erase() {
	for name, _ := range t.views {
		if VERBOSE {
			log.Printf("converter: dropping view '%v'\n", name)
		}

		err := t.r.DropView(name)
		if err != nil {
			log.Println("converter: error while dropping view", name, err)
		}
	}

	for name, _ := range t.projections {
		if VERBOSE {
			log.Printf("converter: dropping projection '%v'\n", name)
		}

		err := t.r.DropProjection(name)
		if err != nil {
			log.Println("converter: error while dropping projection", name, err)
		}
	}
}

func Convert(r common.ReadCloser, w common.WriteCloser, options *Config, verbosity int) error {
	tempViews := createTempEntities(r, options.Views, options.Projections)
	defer tempViews.Erase()

	tables := r.FilteredTables(options.OnlyTables, options.ExcludeTables)

	/* sort the tables according to only tables if "only tables" was
	 * specified. This is a primitive way to be able to specify some
	 * ordering among the tables. */
	OrderTableByNamesList(tables, options.OnlyTablesList)

	/* override types if specified in the options */
	for _, table := range tables {
		/* is this table a projection? */
		meta, ok := options.Projections[table.Name]
		if !ok {
			continue
		}

		/* see if any of the columns require a different type than the
		 * one we derived */
		for _, col := range table.Columns {
			newtype, ok := meta.Types[col.Name]
			if !ok {
				continue
			}

			col.Type = common.SimpleType(newtype)
			col.RawType = newtype
		}
	}

	if !options.SuppressDdl {
		createTables(tables, w)
	}
	if options.Truncate {
		truncateTables(tables, w)
	}
	if !options.SuppressData {
		if options.Merge {
			for _, srcTable := range tables {
				/* is this table a projection? */
				var extraDstCond string
				if meta, ok := options.Projections[srcTable.Name]; ok {
					extraDstCond = meta.Conditions
				}

				if VERBOSE {
					log.Println("converter: merging table", srcTable.Name)
				}

				dstTableName := strmap(srcTable.Name, options.TableMap)
				err := w.MergeTable(srcTable, dstTableName, extraDstCond, r)
				if err != nil {
					return err
				}
			}
		} else {
			writeData(tables, w)
		}
	}

	createIndices(tables, w)
	createConstraints(tables, w)

	return nil
}

func strmap(srcname string, m map[string]string) string {
	if m == nil {
		return srcname
	}
	mapped, ok := m[srcname]
	if !ok {
		return srcname
	}

	return mapped
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
