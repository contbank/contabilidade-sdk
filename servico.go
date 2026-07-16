package contabilidade

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// ConsultarServicoPorCnae sends POST /api/servico/consultar-por-cnae.
func (c *Client) ConsultarServicoPorCnae(
	ctx context.Context,
	req ConsultarServicoPorCnaeRequest,
) (*ConsultarServicoPorCnaeResponse, error) {
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + PathConsultarServicoCnae

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
		"cnaes":      req.Cnaes,
	}).Info("ConsultarServicoPorCnae")

	var response ConsultarServicoPorCnaeResponse
	status, err := c.doJSON(ctx, http.MethodPost, endpoint, req, &response)
	if err != nil {
		logrus.WithError(err).Error("ConsultarServicoPorCnae: http error")
		return nil, err
	}

	if status == http.StatusUnauthorized {
		return nil, newAPIError(ErrUnauthorized, "UNAUTHORIZED", nil)
	}

	if status != http.StatusOK {
		return nil, newAPIError(ErrDefaultConsultarServico, "CONSULTAR_SERVICO_ERROR", []string{
			fmt.Sprintf("unexpected status %d", status),
		})
	}

	return &response, nil
}

// ListarMunicipios sends GET /api/servico/municipios.
func (c *Client) ListarMunicipios(ctx context.Context) ([]string, error) {
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + PathListarMunicipios

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
	}).Info("ListarMunicipios")

	var municipios []string
	status, err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &municipios)
	if err != nil {
		logrus.WithError(err).Error("ListarMunicipios: http error")
		return nil, err
	}

	if status == http.StatusUnauthorized {
		return nil, newAPIError(ErrUnauthorized, "UNAUTHORIZED", nil)
	}

	if status != http.StatusOK {
		return nil, newAPIError(ErrDefaultListarMunicipios, "LISTAR_MUNICIPIOS_ERROR", []string{
			fmt.Sprintf("unexpected status %d", status),
		})
	}

	return municipios, nil
}

// Status sends GET /api/status (no business auth required by the gateway, but
// the shared client still attaches the Bearer token).
func (c *Client) Status(ctx context.Context) error {
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + PathStatus

	status, err := c.doJSON(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("contabilidade status check failed with http %d", status)
	}

	return nil
}
