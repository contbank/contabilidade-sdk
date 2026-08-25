package models

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// FlexAnnex is the Simples Nacional annex (3, 4 or 5) tolerant to Contabilidade.com
// JSON shapes: number, string ("5", "III", ""), or array (["III"]).
type FlexAnnex int

// UnmarshalJSON accepts int | string | []any without failing the parent payload.
func (a *FlexAnnex) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*a = 0
		return nil
	}

	switch b[0] {
	case '[':
		var arr []json.RawMessage
		if err := json.Unmarshal(b, &arr); err != nil || len(arr) == 0 {
			*a = 0
			return nil
		}
		return a.UnmarshalJSON(arr[0])

	case '"':
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			*a = 0
			return nil
		}
		*a = FlexAnnex(parseAnnexToken(s))
		return nil

	default:
		// Number (int or float JSON).
		var n json.Number
		if err := json.Unmarshal(b, &n); err == nil {
			if i, err := n.Int64(); err == nil {
				*a = FlexAnnex(int(i))
				return nil
			}
			if f, err := n.Float64(); err == nil {
				*a = FlexAnnex(int(f))
				return nil
			}
		}
		var f float64
		if err := json.Unmarshal(b, &f); err == nil {
			*a = FlexAnnex(int(f))
			return nil
		}
		*a = 0
		return nil
	}
}

// MarshalJSON encodes as a JSON number (or 0).
func (a FlexAnnex) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(a))
}

// Int returns the annex as int (0 when unknown/empty).
func (a FlexAnnex) Int() int {
	return int(a)
}

func parseAnnexToken(raw string) int {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0
	}

	upper := strings.ToUpper(s)
	switch upper {
	case "III":
		return 3
	case "IV":
		return 4
	case "V":
		return 5
	}

	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return 0
}
