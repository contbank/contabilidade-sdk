package contabilidade

import (
	"bytes"
	"io"
	"net/http"
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
			"body":       string(reqBody),
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
			"body":     string(respBody),
		},
	}).Info("HTTP Response Contabilidade")

	return resp, nil
}

func sanitizeHeaders(h http.Header) http.Header {
	cloned := h.Clone()
	if cloned.Get("Authorization") != "" {
		cloned.Set("Authorization", "Bearer ***")
	}
	return cloned
}
