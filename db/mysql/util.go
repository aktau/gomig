package mysql

import (
	. "github.com/aktau/gomig/db/common"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func MysqlToGenericType(mysqlType string) *Type {
	rt := mysqlType
	switch {
	case strings.Contains(rt, "float"):
		return FloatType()
	case strings.Contains(rt, "double"):
		return DoubleType()
	case strings.Contains(rt, "numeric"), strings.Contains(rt, "decimal"):
		scale, precision := ExtractPrecisionAndScale(rt)
		return NumericType(scale, precision)
	case strings.Contains(rt, "int"):
		return IntType()
	case strings.Contains(rt, "blob"), strings.Contains(rt, "binary"):
		return BlobType()
	case strings.HasPrefix(rt, "char"):
		t := PaddedTextType()
		t.Max = ExtractLength(rt)
		return t
	case strings.Contains(rt, "varchar"), strings.Contains(rt, "text"):
		t := TextType()
		t.Max = ExtractLength(rt)
		return t
	case rt == "bit(1)", rt == "tinyint(1)", rt == "tinyint(1) unsigned":
		return BoolType()
	default:
		log.Println("WARNING: mysql: encountered an unknown type, ", rt)
		return SimpleType(rt)
	}
}

/* returns 0 if no length could be determined */
func ExtractLength(mysqlType string) uint {
	/* matches should be: [mysqlType, length] */
	matches := regexp.MustCompile(`\w+\((\d+)\)`).FindStringSubmatch(mysqlType)

	if len(matches) != 2 {
		return 0
	}

	i, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	return uint(i)
}

/* returns a precision, scale tuple */
func ExtractPrecisionAndScale(mysqlType string) (uint, uint) {
	/* we should get something like: TYPE(precision, scale) */
	/* matches should be: [mysqlType, precision, scale] */
	matches := regexp.MustCompile(`\w+\(\s*(\d+)\s*,\s*(\d+)\s*\)`).FindStringSubmatch(mysqlType)

	if len(matches) != 3 {
		return 0, 0
	}

	precision, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0
	}
	scale, err := strconv.Atoi(matches[2])

	return uint(precision), uint(scale)
}
