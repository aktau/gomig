package common

func FilterInclExcl(list []string, incl map[string]bool, excl map[string]bool) []string {
	filtered := make([]string, 0, len(list)/4)

	allIncl := (incl == nil || len(incl) == 0)
	noExcl := (excl == nil || len(excl) == 0)

	for _, str := range list {
		isIncl := allIncl || incl[str]
		isExcl := !noExcl && excl[str]

		if !isExcl && isIncl {
			filtered = append(filtered, str)
		}
	}

	return filtered
}
