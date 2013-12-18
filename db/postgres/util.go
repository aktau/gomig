package postgres

import (
	"fmt"
	"github.com/aktau/gomig/db/common"
)

func PostgresToGenericType(postgresType string) string {
	return postgresType
}

func GenericToPostgresType(genericType *common.Type) string {
	gen := genericType
	name := gen.Name
	max := gen.Max
	precision := gen.Precision
	scale := gen.Scale

	switch name {
	case "text":
		/* the typical varchar type if its maximum is lower than 200, we
		 * assume they actually meant it */
		if gen.HasMax() && max < 200 {
			return fmt.Sprintf("character varying(%v)", max)
		}

		/* if the text type has no maximum (or a maximum above 200, we assume text) */
		return "text"
	case "char":
		return fmt.Sprintf("character(%v)", max)
	case "float":
		return "real"
	case "double":
		return "double precision"
	case "numeric":
		return fmt.Sprintf("numeric(%v, %v)", precision, scale)
	case "blob":
		return "bytea"
	default:
		return name
	}
}

/* converts a RawBytes field into something you can
 * put into a regular insert statement (wrapping strings in $$
 * et cetera) */
func RawToPostgres(val []byte, origType *common.Type) (string, error) {
	if val == nil {
		return "NULL", nil
	} else {
		switch origType.Name {
		case "text", "char":
			return "$$" + string(val) + "$$", nil
		case "boolean":
			/* ascii(48) = "0" and ascii(49) = "1" */
			switch val[0] {
			case 48:
				return "f", nil
			case 49:
				return "t", nil
			default:
				return "", fmt.Errorf("postgres: did not recognize bool value: string(%v) = %v, val[0] = %v", val, string(val), val[0])
			}
		case "integer", "float", "double":
			return string(val), nil
		default:
			return string(val), nil
		}
	}
}
