package contabilidade

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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
	s = strings.TrimSpace(s)
	if s == "" {
		t.Time = time.Time{}
		return nil
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			t.Time = parsed
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("cannot parse time %q: %w", s, lastErr)
}

// MarshalJSON encodes as RFC3339Nano (UTC) for requests we send.
func (t APITime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.UTC().Format(time.RFC3339Nano))
}
