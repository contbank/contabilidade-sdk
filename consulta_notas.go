package contabilidade

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// SolicitarConsulta sends POST /api/consulta-notas/solicitar and returns a tracking token.
// Processing is asynchronous: poll ConsultarStatus then BaixarNotas when Finalizado.
func (c *Client) SolicitarConsulta(
	ctx context.Context,
	req SolicitarConsultaNotasRequest,
) (*SolicitarConsultaNotasResponse, error) {
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + PathSolicitarConsultaNotas

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
		"ibge":       stringPtrValue(req.MunicipioCodigoIBGE),
	}).Info("SolicitarConsulta")

	var response SolicitarConsultaNotasResponse
	status, err := c.doJSON(ctx, http.MethodPost, endpoint, req, &response)
	if err != nil {
		logrus.WithError(err).Error("SolicitarConsulta: http error")
		return nil, err
	}

	if status == http.StatusUnauthorized {
		return nil, newAPIError(ErrUnauthorized, "UNAUTHORIZED", nil)
	}

	if status != http.StatusOK && status != http.StatusCreated {
		msgs := response.Erros
		if len(msgs) == 0 {
			msgs = []string{fmt.Sprintf("unexpected status %d", status)}
		}
		return &response, newAPIError(ErrDefaultSolicitarConsultaNotas, "SOLICITAR_CONSULTA_NOTAS_ERROR", msgs)
	}

	if !response.Sucesso {
		return &response, newAPIError(ErrSolicitarConsultaNotasRejected, "SOLICITAR_CONSULTA_NOTAS_REJECTED", response.Erros)
	}

	return &response, nil
}

// ConsultarStatus sends GET /api/consulta-notas/{token}/status.
func (c *Client) ConsultarStatus(ctx context.Context, token string) (*ConsultaNotasStatusResponse, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, newAPIError(ErrInvalidRequest, "INVALID_REQUEST", []string{"token is required"})
	}

	path := fmt.Sprintf(PathConsultaNotasStatus, token)
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + path

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
		"token":      token,
	}).Info("ConsultarStatus")

	var response ConsultaNotasStatusResponse
	status, err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &response)
	if err != nil {
		logrus.WithError(err).Error("ConsultarStatus: http error")
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return &response, nil
	case http.StatusUnauthorized:
		return nil, newAPIError(ErrUnauthorized, "UNAUTHORIZED", nil)
	case http.StatusNotFound:
		msgs := response.Erros
		if len(msgs) == 0 {
			msgs = []string{"consulta notas token not found"}
		}
		return &response, newAPIError(ErrConsultaNotasNotFound, "CONSULTA_NOTAS_NOT_FOUND", msgs)
	default:
		msgs := response.Erros
		if len(msgs) == 0 {
			msgs = []string{fmt.Sprintf("unexpected status %d", status)}
		}
		return &response, newAPIError(ErrDefaultConsultaNotasStatus, "CONSULTA_NOTAS_STATUS_ERROR", msgs)
	}
}

// BaixarNotas sends GET /api/consulta-notas/{token}/notas (list with XmlBase64 per note).
// If StatusConsulta is not Finalizado, the list may still be incomplete.
func (c *Client) BaixarNotas(ctx context.Context, token string) (*ListarNotasConsultadasResponse, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, newAPIError(ErrInvalidRequest, "INVALID_REQUEST", []string{"token is required"})
	}

	path := fmt.Sprintf(PathConsultaNotasListar, token)
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + path

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
		"token":      token,
	}).Info("BaixarNotas")

	var response ListarNotasConsultadasResponse
	status, err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &response)
	if err != nil {
		logrus.WithError(err).Error("BaixarNotas: http error")
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return &response, nil
	case http.StatusUnauthorized:
		return nil, newAPIError(ErrUnauthorized, "UNAUTHORIZED", nil)
	case http.StatusNotFound:
		msgs := response.Erros
		if len(msgs) == 0 {
			msgs = []string{"consulta notas token not found"}
		}
		return &response, newAPIError(ErrConsultaNotasNotFound, "CONSULTA_NOTAS_NOT_FOUND", msgs)
	default:
		msgs := response.Erros
		if len(msgs) == 0 {
			msgs = []string{fmt.Sprintf("unexpected status %d", status)}
		}
		return &response, newAPIError(ErrDefaultBaixarNotas, "BAIXAR_NOTAS_ERROR", msgs)
	}
}

func stringPtrValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
