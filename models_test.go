package contabilidade_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	contabilidade "github.com/contbank/contabilidade-sdk"
	"github.com/contbank/contabilidade-sdk/pkg/models"
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
	if got.Anexo != 5 {
		t.Fatalf("Anexo = %d, want 5", got.Anexo)
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
	if got.Anexo != 0 {
		t.Fatalf("Anexo = %d, want 0", got.Anexo)
	}
}

func TestServiceSuggestion_unmarshalsAnexoRomanArrayAndNumber(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want models.FlexAnnex
	}{
		{name: "roman", raw: `{"Anexo":"III"}`, want: 3},
		{name: "array", raw: `{"Anexo":["IV"]}`, want: 4},
		{name: "number", raw: `{"Anexo":5}`, want: 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got contabilidade.CodigoServicoDto
			if err := json.Unmarshal([]byte(tc.raw), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.Anexo != tc.want {
				t.Fatalf("Anexo = %d, want %d", got.Anexo, tc.want)
			}
		})
	}
}

func TestConsultaNotasStatusResponse_unmarshalsDotNetTimestamps(t *testing.T) {
	raw := []byte(`{
		"Sucesso":true,
		"Token":"2e9dd291-e10f-47d9-ab12-11a17db2acf9",
		"Finalizado":false,
		"Status":"Erro",
		"DataCriacao":"2026-08-25T14:34:39.6866667",
		"DataFinalizacao":"2026-08-25T14:34:41.1866667",
		"Erros":[]
	}`)

	var got contabilidade.ConsultaNotasStatusResponse
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Status == nil || *got.Status != contabilidade.ConsultaNotasStatusErro {
		t.Fatalf("Status = %#v", got.Status)
	}
	if got.DataCriacao == nil {
		t.Fatal("DataCriacao nil")
	}
	wantCriacao := time.Date(2026, 8, 25, 14, 34, 39, 686666700, time.UTC)
	if !got.DataCriacao.Equal(wantCriacao) {
		t.Fatalf("DataCriacao = %v, want %v", got.DataCriacao.Time, wantCriacao)
	}
	if got.DataFinalizacao == nil {
		t.Fatal("DataFinalizacao nil")
	}
}

func TestAPITime_unmarshalsRFC3339(t *testing.T) {
	raw := []byte(`"2026-08-25T14:34:39.6866667Z"`)
	var got contabilidade.APITime
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Location() != time.UTC {
		t.Fatalf("location = %v", got.Location())
	}
}

func TestCancelarNfseRequest_marshalsNacionalFields(t *testing.T) {
	codigo := int32(contabilidade.MotivoCancelamentoErroEmissao)
	req := contabilidade.CancelarNfseRequest{
		Numero:                      131,
		CodigoVerificacao:           contabilidade.String("31062002209474631000117000000000013126086084662553"),
		InscricaoMunicipal:          contabilidade.String("02234900012"),
		MunicipioCodigoIBGE:         contabilidade.String("3106200"),
		ChaveAcesso:                 contabilidade.String("31062002209474631000117000000000013126086084662553"),
		MotivoCancelamentoCodigo:    &codigo,
		MotivoCancelamentoDescricao: contabilidade.String("Nota emitida com valor incorreto"),
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(raw)
	for _, want := range []string{
		`"MunicipioCodigoIBGE":"3106200"`,
		`"ChaveAcesso":"31062002209474631000117000000000013126086084662553"`,
		`"MotivoCancelamentoCodigo":1`,
		`"MotivoCancelamentoDescricao":"Nota emitida com valor incorreto"`,
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("payload missing %s\ngot %s", want, s)
		}
	}
}
