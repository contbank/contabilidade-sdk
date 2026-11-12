package contabilidade

import (
	"strings"
	"testing"
)

// Prod (12/11/2026): o log "HTTP Request Contabilidade" levava o PFX (EmissorBase64) e a senha (EmissorSenha) da emissão.
func TestSanitizeBodyRemovesCertificateAndPassword(t *testing.T) {
	body := []byte(`{"Certificado":{"EmissorDocumento":"33011770000199","EmissorBase64":"MIIOGwIBAzCCDdcG","EmissorSenha":"segredo"},` +
		`"Nota":{"RPS":{"Serie":"A","Numero":"339"},"Emissor":{"InscricaoMunicipal":"62141830"}}}`)

	got := sanitizeBody(body)

	for _, secret := range []string{"MIIOGwIBAzCCDdcG", "segredo"} {
		if strings.Contains(got, secret) {
			t.Fatalf("log ainda tem %q: %s", secret, got)
		}
	}
	for _, kept := range []string{`"EmissorDocumento":"33011770000199"`, `"Numero":"339"`, `"InscricaoMunicipal":"62141830"`} {
		if !strings.Contains(got, kept) {
			t.Fatalf("log perdeu %s: %s", kept, got)
		}
	}
}

func TestSanitizeBodyRemovesBase64InListsAndKeepsErrors(t *testing.T) {
	body := []byte(`{"Sucesso":false,"Erros":["E0120 - IM do prestador"],"Notas":[{"Numero":1,"XmlBase64":"PD94bWwg"}]}`)

	got := sanitizeBody(body)

	if strings.Contains(got, "PD94bWwg") {
		t.Fatalf("XmlBase64 ficou no log: %s", got)
	}
	if !strings.Contains(got, "E0120 - IM do prestador") {
		t.Fatalf("mensagem de erro sumiu do log: %s", got)
	}
}

func TestSanitizeBodyNonJSON(t *testing.T) {
	if got := sanitizeBody([]byte("%PDF-1.7 binario")); got != "<corpo não JSON, 16 bytes>" {
		t.Fatalf("corpo não JSON: %q", got)
	}
	if got := sanitizeBody(nil); got != "" {
		t.Fatalf("corpo vazio: %q", got)
	}
}
