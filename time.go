package contabilidade

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// dateTimeOffsetRegex detects when a '+' timezone offset was decoded as a space
// (common with url.ParseQuery), e.g. "2026-07-01T05:00:00.000 02:00".
var dateTimeOffsetRegex = regexp.MustCompile(`^(.*T\d{2}:\d{2}:\d{2}(\.\d+)?)\s(\d{2}:?\d{2})$`)

// NormalizeDateTimeOffset restores a positive timezone offset that became a space.
func NormalizeDateTimeOffset(dateTime string) string {
	dateTime = strings.TrimSpace(dateTime)
	return dateTimeOffsetRegex.ReplaceAllString(dateTime, "$1+$3")
}

// APITime unmarshals Contabilidade.com timestamps that may omit timezone
// (e.g. .NET DateTime: "2026-08-25T14:34:39.6866667").
type APITime struct {
	time.Time
}

// UnmarshalJSON accepts RFC3339 and local/untyped ISO-8601 variants.
func (t *APITime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseFlexibleDateTime(s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// MarshalJSON encodes as RFC3339Nano (UTC) for requests we send.
func (t APITime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.UTC().Format(time.RFC3339Nano))
}

// ParseFlexibleDateTime parses Contabilidade / prefeitura timestamps in several shapes.
func ParseFlexibleDateTime(value string) (time.Time, error) {
	value = NormalizeDateTimeOffset(strings.TrimSpace(value))
	if value == "" {
		return time.Time{}, nil
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02/01/2006 15:04:05",
		"02/01/2006",
	}

	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}
	return time.Time{}, fmt.Errorf("cannot parse time %q: %w", value, lastErr)
}
