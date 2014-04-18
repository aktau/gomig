package postgres

import (
	"database/sql"
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
	modifier := gen.Modifier

	switch name {
	case common.TypeText:
		/* the typical varchar type if its maximum is lower than 200, we
		 * assume they actually meant it */
		if gen.HasMax() && max < 200 {
			return fmt.Sprintf("character varying(%v)", max)
		}

		/* if the text type has no maximum (or a maximum above 200, we assume text) */
		return "text"
	case common.TypeChar:
		return fmt.Sprintf("character(%v)", max)
	case common.TypeFloat:
		return "real"
	case common.TypeDouble:
		return "double precision"
	case common.TypeNumeric:
		return fmt.Sprintf("numeric(%v, %v)", precision, scale)
	case common.TypeBit:
		return fmt.Sprintf("bit varying(%v)", max)
	case common.TypeBlob:
		return "bytea"
	case common.TypeInteger:
		switch modifier {
		case common.TypeSmall:
			return "smallint"
		case common.TypeNormal:
			return "integer"
		case common.TypeLarge:
			return "bigint"
		case common.TypeHuge:
			/* from: http://en.wikibooks.org/wiki/Converting_MySQL_to_PostgreSQL */
			return "numeric(20)"
		default:
			return "integer"
		}
	case common.TypeSet:
		return "text[]"
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
		case common.TypeText, common.TypeChar:
			return "$$" + string(val) + "$$", nil
		case common.TypeBool:
			/* ascii(48) = "0" and ascii(49) = "1" */
			switch val[0] {
			case 48:
				return "f", nil
			case 49:
				return "t", nil
			default:
				return "", fmt.Errorf("postgres: did not recognize bool value: string(%v) = %v, val[0] = %v", val, string(val), val[0])
			}
		case common.TypeNumeric, common.TypeInteger, common.TypeFloat, common.TypeDouble:
			return string(val), nil
		default:
			return string(val), nil
		}
	}
}

func NewTypedSlice(src *common.Table) []interface{} {
	vals := make([]interface{}, len(src.Columns))
	for i, col := range src.Columns {
		switch col.Type.Name {
		case common.TypeBool:
			if col.Null {
				vals[i] = new(sql.NullBool)
			} else {
				vals[i] = new(bool)
			}
		case common.TypeNumeric, common.TypeFloat, common.TypeDouble:
			if col.Null {
				vals[i] = new(sql.NullFloat64)
			} else {
				vals[i] = new(float64)
			}
		case common.TypeInteger:
			if col.Null {
				vals[i] = new(sql.NullInt64)
			} else {
				vals[i] = new(int64)
			}
		case common.TypeBlob:
			/* do we have a suitable NullBlob or NullByte somewhere? I bet
			 * this gives problems somehow with NULLable blob fields... */
			vals[i] = new([]byte)
		default:
			if col.Null {
				vals[i] = new(sql.NullString)
			} else {
				vals[i] = new(string)
			}
		}
	}

	return vals
}
