package contabilidade_test

import (
	"encoding/json"
	"testing"

	contabilidade "github.com/contbank/contabilidade-sdk"
)

func TestServiceSuggestion_unmarshalsAnexoAsString(t *testing.T) {
	raw := []byte(`{
		"Codigo":"02666",
		"CodigoServicoMunicipal":"02666",
		"Descricao":"Programação",
		"Anexo":"5",
		"EmiteNFSe":true
	}`)

	var got contabilidade.ServiceSuggestion
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Anexo != "5" {
		t.Fatalf("Anexo = %q, want %q", got.Anexo, "5")
	}
	if got.Codigo == nil || *got.Codigo != "02666" {
		t.Fatalf("Codigo = %#v", got.Codigo)
	}
	if got.CodigoServicoMunicipal == nil || *got.CodigoServicoMunicipal != "02666" {
		t.Fatalf("CodigoServicoMunicipal = %#v", got.CodigoServicoMunicipal)
	}
}

func TestServiceSuggestion_unmarshalsEmptyAnexo(t *testing.T) {
	raw := []byte(`{"Codigo":"02666","Anexo":"","EmiteNFSe":true}`)

	var got contabilidade.ServiceSuggestion
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal empty Anexo: %v", err)
	}
	if got.Anexo != "" {
		t.Fatalf("Anexo = %q, want empty", got.Anexo)
	}
}
