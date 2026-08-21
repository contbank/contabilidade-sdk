package contabilidade

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/contbank/grok"
)

var (
	// ErrInvalidAPIEndpoint ...
	ErrInvalidAPIEndpoint = grok.NewError(http.StatusBadRequest, "INVALID_API_ENDPOINT", "invalid api endpoint")
	// ErrUnauthorized ...
	ErrUnauthorized = grok.NewError(http.StatusUnauthorized, "UNAUTHORIZED", "invalid or missing bearer token")
	// ErrDefaultEmitirNfse ...
	ErrDefaultEmitirNfse = grok.NewError(http.StatusConflict, "EMITIR_NFSE_ERROR", "error emitting nfse")
	// ErrEmitirNfseRejected ...
	ErrEmitirNfseRejected = grok.NewError(http.StatusBadRequest, "EMITIR_NFSE_REJECTED", "nfse emission rejected by municipality")
	// ErrDefaultCancelarNfse ...
	ErrDefaultCancelarNfse = grok.NewError(http.StatusConflict, "CANCELAR_NFSE_ERROR", "error cancelling nfse")
	// ErrCancelarNfseRejected ...
	ErrCancelarNfseRejected = grok.NewError(http.StatusBadRequest, "CANCELAR_NFSE_REJECTED", "nfse cancellation rejected by municipality")
	// ErrDefaultBaixarXml ...
	ErrDefaultBaixarXml = grok.NewError(http.StatusConflict, "BAIXAR_XML_NFSE_ERROR", "error downloading nfse xml")
	// ErrDefaultBaixarDanfse ...
	ErrDefaultBaixarDanfse = grok.NewError(http.StatusConflict, "BAIXAR_DANFSE_ERROR", "error downloading national danfse pdf")
	// ErrNfseNotFound ...
	ErrNfseNotFound = grok.NewError(http.StatusNotFound, "NFSE_NOT_FOUND", "nfse not found")
	// ErrDefaultConsultarServico ...
	ErrDefaultConsultarServico = grok.NewError(http.StatusConflict, "CONSULTAR_SERVICO_ERROR", "error consulting service by cnae")
	// ErrDefaultListarMunicipios ...
	ErrDefaultListarMunicipios = grok.NewError(http.StatusConflict, "LISTAR_MUNICIPIOS_ERROR", "error listing municipalities")
	// ErrDefaultSolicitarConsultaNotas ...
	ErrDefaultSolicitarConsultaNotas = grok.NewError(http.StatusConflict, "SOLICITAR_CONSULTA_NOTAS_ERROR", "error requesting notes consultation")
	// ErrSolicitarConsultaNotasRejected ...
	ErrSolicitarConsultaNotasRejected = grok.NewError(http.StatusBadRequest, "SOLICITAR_CONSULTA_NOTAS_REJECTED", "notes consultation request rejected")
	// ErrDefaultConsultaNotasStatus ...
	ErrDefaultConsultaNotasStatus = grok.NewError(http.StatusConflict, "CONSULTA_NOTAS_STATUS_ERROR", "error consulting notes consultation status")
	// ErrConsultaNotasNotFound ...
	ErrConsultaNotasNotFound = grok.NewError(http.StatusNotFound, "CONSULTA_NOTAS_NOT_FOUND", "notes consultation token not found")
	// ErrDefaultBaixarNotas ...
	ErrDefaultBaixarNotas = grok.NewError(http.StatusConflict, "BAIXAR_NOTAS_ERROR", "error listing consulted notes")
	// ErrInvalidRequest ...
	ErrInvalidRequest = grok.NewError(http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
)

// Error wraps a grok.Error with a stable key for consumers.
type Error struct {
	ErrorKey  string
	GrokError *grok.Error
	Messages  []string
}

// Error implements the error interface.
func (e *Error) Error() string {
	msgs := e.Messages
	if e.GrokError != nil && len(e.GrokError.Messages) > 0 {
		msgs = e.GrokError.Messages
	}
	return fmt.Sprintf("Key: %s - Messages: %s", e.ErrorKey, strings.Join(msgs, "\n"))
}

// ParseErr extracts a contabilidade.Error from a generic error.
func ParseErr(err error) (*Error, bool) {
	ce, ok := err.(*Error)
	return ce, ok
}

func newAPIError(base *grok.Error, key string, messages []string) *Error {
	cloned := *base
	if len(messages) > 0 {
		cloned.Messages = messages
	}
	return &Error{
		ErrorKey:  key,
		GrokError: &cloned,
		Messages:  messages,
	}
}
