package contabilidade

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	// ApiEndpoint is the default Contabilidade.com NFS-e gateway.
	ApiEndpoint string = "https://nfse.contabilidade.com"
)

// Config holds optional overrides used to build a Session (pointer options, AWS-style).
// Map these fields from the microservice config.yaml section `contabilidade:`.
type Config struct {
	APIEndpoint *string
	BearerToken *string
	Timeout     *time.Duration
}

// Session is the materialized runtime configuration shared by all domain clients.
type Session struct {
	APIEndpoint string
	BearerToken string
	Timeout     time.Duration
}

// bearerTransport injects Authorization: Bearer on every request.
type bearerTransport struct {
	underlyingTransport http.RoundTripper
	token               string
	mutex               *sync.Mutex
}

// NewSession materializes a Session from Config, applying defaults and env fallbacks.
func NewSession(config Config) (*Session, error) {
	if config.APIEndpoint == nil {
		config.APIEndpoint = String(ApiEndpoint)
	}

	if config.BearerToken == nil || *config.BearerToken == "" {
		envToken := os.Getenv("CONTABILIDADE_API_TOKEN")
		config.BearerToken = String(envToken)
	}

	if config.BearerToken == nil || *config.BearerToken == "" {
		return nil, fmt.Errorf("contabilidade: bearer token is required (config or CONTABILIDADE_API_TOKEN)")
	}

	timeout := 60 * time.Second
	if config.Timeout != nil {
		timeout = *config.Timeout
	}

	return &Session{
		APIEndpoint: *config.APIEndpoint,
		BearerToken: *config.BearerToken,
		Timeout:     timeout,
	}, nil
}

// CreateBearerHTTPClient returns an *http.Client that injects the Bearer token.
// Contabilidade.com Swagger uses HTTP Bearer auth (static API token), not OAuth2.
func CreateBearerHTTPClient(session *Session) *http.Client {
	return &http.Client{
		Timeout: session.Timeout,
		Transport: &bearerTransport{
			underlyingTransport: http.DefaultTransport,
			token:               session.BearerToken,
			mutex:               &sync.Mutex{},
		},
	}
}

// RoundTrip adds the Authorization Bearer header.
func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mutex.Lock()
	token := t.token
	t.mutex.Unlock()

	req = req.Clone(req.Context())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	return t.underlyingTransport.RoundTrip(req)
}
