package contabilidade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingHTTPClient wraps *http.Client and logs request/response payloads.
type LoggingHTTPClient struct {
	client *http.Client
}

// NewLoggingHTTPClient creates a LoggingHTTPClient.
func NewLoggingHTTPClient(client *http.Client) *LoggingHTTPClient {
	return &LoggingHTTPClient{client: client}
}

// Do executes the HTTP request with structured logging.
func (c *LoggingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	var reqBody []byte
	if req != nil && req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	logrus.WithFields(logrus.Fields{
		"contabilidade_request": logrus.Fields{
			"method":     req.Method,
			"url":        req.URL.String(),
			"header":     sanitizeHeaders(req.Header),
			"body":       sanitizeBody(reqBody),
			"user-agent": req.UserAgent(),
		},
	}).Info("HTTP Request Contabilidade")

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("HTTP request failed")
		return nil, err
	}

	duration := time.Since(start)

	var respBody []byte
	if resp != nil && resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}

	logrus.WithFields(logrus.Fields{
		"contabilidade_response": logrus.Fields{
			"header":   resp.Header,
			"status":   resp.StatusCode,
			"duration": duration,
			"body":     sanitizeBody(respBody),
		},
	}).Info("HTTP Response Contabilidade")

	return resp, nil
}

// redactedValue substitui o valor de campos sensíveis no log.
const redactedValue = "***"

// sensitiveKeyParts nome de campo que contém qualquer um destes trechos não vai para o log:
// certificado (EmissorBase64), senha (EmissorSenha) e artefatos em base64 (XmlBase64).
var sensitiveKeyParts = []string{"senha", "password", "pfx", "base64"}

// sanitizeBody devolve o corpo para o log sem certificado, senha e base64. Corpo que não é JSON
// vai só com o tamanho: não há como saber o que tem dentro.
func sanitizeBody(body []byte) string {
	if len(bytes.TrimSpace(body)) == 0 {
		return ""
	}
	var parsed interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Sprintf("<corpo não JSON, %d bytes>", len(body))
	}
	out, err := json.Marshal(redactSensitive(parsed))
	if err != nil {
		return fmt.Sprintf("<corpo não JSON, %d bytes>", len(body))
	}
	return string(out)
}

func redactSensitive(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, item := range v {
			if isSensitiveKey(key) {
				if item != nil {
					v[key] = redactedValue
				}
				continue
			}
			v[key] = redactSensitive(item)
		}
		return v
	case []interface{}:
		for i, item := range v {
			v[i] = redactSensitive(item)
		}
		return v
	default:
		return v
	}
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, part := range sensitiveKeyParts {
		if strings.Contains(lower, part) {
			return true
		}
	}
	return false
}

func sanitizeHeaders(h http.Header) http.Header {
	cloned := h.Clone()
	if cloned.Get("Authorization") != "" {
		cloned.Set("Authorization", "Bearer ***")
	}
	return cloned
}
