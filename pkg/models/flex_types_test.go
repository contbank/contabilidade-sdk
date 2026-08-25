package models_test

import (
	"encoding/json"
	"testing"

	"github.com/contbank/contabilidade-sdk/pkg/models"
)

func TestFlexAnnex_unmarshalsNumber(t *testing.T) {
	var got models.FlexAnnex
	if err := json.Unmarshal([]byte(`5`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != 5 {
		t.Fatalf("got %d, want 5", got)
	}
}

func TestFlexAnnex_unmarshalsNumericString(t *testing.T) {
	var got models.FlexAnnex
	if err := json.Unmarshal([]byte(`"5"`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != 5 {
		t.Fatalf("got %d, want 5", got)
	}
}

func TestFlexAnnex_unmarshalsRoman(t *testing.T) {
	cases := map[string]models.FlexAnnex{
		`"III"`: 3,
		`"iv"`:  4,
		`"V"`:   5,
	}
	for raw, want := range cases {
		var got models.FlexAnnex
		if err := json.Unmarshal([]byte(raw), &got); err != nil {
			t.Fatalf("%s: %v", raw, err)
		}
		if got != want {
			t.Fatalf("%s: got %d, want %d", raw, got, want)
		}
	}
}

func TestFlexAnnex_unmarshalsArray(t *testing.T) {
	var got models.FlexAnnex
	if err := json.Unmarshal([]byte(`["III"]`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != 3 {
		t.Fatalf("got %d, want 3", got)
	}
}

func TestFlexAnnex_unmarshalsEmptyStringAndNull(t *testing.T) {
	var got models.FlexAnnex
	if err := json.Unmarshal([]byte(`""`), &got); err != nil {
		t.Fatalf("empty: %v", err)
	}
	if got != 0 {
		t.Fatalf("empty: got %d, want 0", got)
	}
	if err := json.Unmarshal([]byte(`null`), &got); err != nil {
		t.Fatalf("null: %v", err)
	}
	if got != 0 {
		t.Fatalf("null: got %d, want 0", got)
	}
}

func TestFlexAnnex_unmarshalsGarbageWithoutError(t *testing.T) {
	var got models.FlexAnnex
	if err := json.Unmarshal([]byte(`{"x":1}`), &got); err != nil {
		t.Fatalf("object must not fail parent: %v", err)
	}
	if got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
}
