package main

import (
	"github.com/aktau/gomig/db/common"
	"sort"
)

/* this file deals with trying to migrate tables in the right order, which
 * is very primitive for now. I generally like Go, but the sort interface
 * is... convoluted. */

/* By is the type of a "less" function that defines the ordering of its
 * Table arguments. */
type By func(t1, t2 *common.Table) bool

/* Sort is a method on the function type, By, that sorts the argument slice
 * according to the function. */
func (by By) Sort(tables []*common.Table) {
	/* the Sort method's receiver is the function (closure)
	 * that defines the sort order. */
	ps := &tableSorter{
		tables: tables,
		by:     by,
	}
	sort.Sort(ps)
}

/* planetSorter joins a By function and a slice of Tables to be sorted. */
type tableSorter struct {
	tables []*common.Table
	by     By /* Closure used in the Less method. */
}

/* implement the sort.Interface. Less is implemented by calling the "by"
 * closure in the sorter. */
func (s *tableSorter) Len() int { return len(s.tables) }
func (s *tableSorter) Swap(i, j int) {
	s.tables[i], s.tables[j] = s.tables[j], s.tables[i]
}
func (s *tableSorter) Less(i, j int) bool {
	return s.by(s.tables[i], s.tables[j])
}

/* Sort the src list of tables in a such a way that the order of the names
 * list is respected. */
func OrderTableByNamesList(src []*common.Table, names []string) {
	if len(src) == 0 || len(names) == 0 {
		return
	}

	lookup := make(map[string]int)
	for idx, name := range names {
		lookup[name] = idx
	}

	sorter := func(t1, t2 *common.Table) bool {
		idx1 := lookup[t1.Name]
		idx2 := lookup[t2.Name]

		return idx1 < idx2
	}

	By(sorter).Sort(src)
}
