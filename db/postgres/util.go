package postgres

import (
	"fmt"
)

/* converts a RawBytes field into something you can
 * put into an insert statemtn (wrapping strings in $$
 * et cetera) */
func RawToPostgres(val []byte, origType string) (string, error) {
	if val == nil {
		return "NULL", nil
	} else {
		switch origType {
		case "text":
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
		case "integer":
			return string(val), nil
		case "float":
			return string(val), nil
		default:
			return string(val), nil
		}
	}
}
