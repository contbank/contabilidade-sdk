package contabilidade

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// EmitirNfse sends POST /api/nfse/emitir.
func (c *Client) EmitirNfse(ctx context.Context, req EmitirNfseRequest) (*EmitirNfseResponse, error) {
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + PathEmitirNfse

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
	}).Info("EmitirNfse")

	var response EmitirNfseResponse
	status, err := c.doJSON(ctx, http.MethodPost, endpoint, req, &response)
	if err != nil {
		logrus.WithError(err).Error("EmitirNfse: http error")
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
		return &response, newAPIError(ErrDefaultEmitirNfse, "EMITIR_NFSE_ERROR", msgs)
	}

	if !response.Sucesso {
		return &response, newAPIError(ErrEmitirNfseRejected, "EMITIR_NFSE_REJECTED", response.Erros)
	}

	return &response, nil
}

// CancelarNfse sends POST /api/nfse/cancelar.
func (c *Client) CancelarNfse(ctx context.Context, req CancelarNfseRequest) (*CancelarNfseResponse, error) {
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + PathCancelarNfse

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
		"numero":     req.Numero,
	}).Info("CancelarNfse")

	var response CancelarNfseResponse
	status, err := c.doJSON(ctx, http.MethodPost, endpoint, req, &response)
	if err != nil {
		logrus.WithError(err).Error("CancelarNfse: http error")
		return nil, err
	}

	if status == http.StatusUnauthorized {
		return nil, newAPIError(ErrUnauthorized, "UNAUTHORIZED", nil)
	}

	if status != http.StatusOK {
		msgs := response.Erros
		if len(msgs) == 0 {
			msgs = []string{fmt.Sprintf("unexpected status %d", status)}
		}
		return &response, newAPIError(ErrDefaultCancelarNfse, "CANCELAR_NFSE_ERROR", msgs)
	}

	if !response.Sucesso {
		return &response, newAPIError(ErrCancelarNfseRejected, "CANCELAR_NFSE_REJECTED", response.Erros)
	}

	return &response, nil
}

// BaixarXmlNfse sends POST /api/nfse/xml/{im}/{numero}/{codigo} and returns raw XML bytes.
func (c *Client) BaixarXmlNfse(
	ctx context.Context,
	inscricaoMunicipal string,
	numero int64,
	codigoVerificacao string,
	cert CertificadoDto,
) ([]byte, error) {
	path := fmt.Sprintf(PathBaixarXmlNfse, inscricaoMunicipal, numero, codigoVerificacao)
	endpoint := strings.TrimRight(c.session.APIEndpoint, "/") + path

	logrus.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": ctx.Value("Request-Id"),
		"im":         inscricaoMunicipal,
		"numero":     numero,
	}).Info("BaixarXmlNfse")

	payload, err := json.Marshal(cert)
	if err != nil {
		return nil, fmt.Errorf("BaixarXmlNfse: marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("BaixarXmlNfse: new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/xml, application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("BaixarXmlNfse: read body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	case http.StatusUnauthorized:
		return nil, newAPIError(ErrUnauthorized, "UNAUTHORIZED", nil)
	case http.StatusNotFound:
		return nil, newAPIError(ErrNfseNotFound, "NFSE_NOT_FOUND", []string{string(body)})
	default:
		return nil, newAPIError(ErrDefaultBaixarXml, "BAIXAR_XML_NFSE_ERROR", []string{
			fmt.Sprintf("status %d: %s", resp.StatusCode, string(body)),
		})
	}
}

func (c *Client) doJSON(ctx context.Context, method, endpoint string, in any, out any) (int, error) {
	var bodyReader io.Reader
	if in != nil {
		payload, err := json.Marshal(in)
		if err != nil {
			return 0, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return 0, fmt.Errorf("new request: %w", err)
	}
	if in != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read body: %w", err)
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return resp.StatusCode, fmt.Errorf("unmarshal response: %w; body=%s", err, string(respBody))
		}
	}

	return resp.StatusCode, nil
}
