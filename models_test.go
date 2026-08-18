package contabilidade_test

import (
	"encoding/json"
	"testing"

	contabilidade "github.com/contbank/contabilidade-sdk"
)

func TestServiceSuggestion_unmarshalsAnexo(t *testing.T) {
	raw := []byte(`{"Codigo":"02666","Descricao":"Programação","Anexo":5,"EmiteNFSe":true}`)

	var got contabilidade.ServiceSuggestion
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Anexo != 5 {
		t.Fatalf("Anexo = %d, want 5", got.Anexo)
	}
	if got.Codigo == nil || *got.Codigo != "02666" {
		t.Fatalf("Codigo = %#v", got.Codigo)
	}
}
