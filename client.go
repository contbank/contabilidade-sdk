package contabilidade

import (
	"context"
	"net/http"
)

// ContabilidadeClient defines the main Contabilidade.com NFS-e gateway operations.
// Prefer injecting this interface in microservices for testability.
type ContabilidadeClient interface {
	// EmitirNfse emits an NFS-e via PMSP (POST /api/nfse/emitir).
	EmitirNfse(ctx context.Context, req EmitirNfseRequest) (*EmitirNfseResponse, error)

	// CancelarNfse cancels a previously issued NFS-e (POST /api/nfse/cancelar).
	CancelarNfse(ctx context.Context, req CancelarNfseRequest) (*CancelarNfseResponse, error)

	// BaixarXmlNfse downloads the XML of an issued NFS-e
	// (POST /api/nfse/xml/{im}/{numero}/{codigo}).
	BaixarXmlNfse(ctx context.Context, inscricaoMunicipal string, numero int64, codigoVerificacao string, cert CertificadoDto) ([]byte, error)

	// ConsultarServicoPorCnae returns ISS service codes for one or more CNAEs
	// (POST /api/servico/consultar-por-cnae).
	ConsultarServicoPorCnae(ctx context.Context, req ConsultarServicoPorCnaeRequest) (*ConsultarServicoPorCnaeResponse, error)

	// ListarMunicipios lists municipalities with service data (GET /api/servico/municipios).
	ListarMunicipios(ctx context.Context) ([]string, error)

	// Status performs the unauthenticated health check (GET /api/status).
	Status(ctx context.Context) error
}

// Client is the concrete Contabilidade.com SDK client.
type Client struct {
	session    Session
	httpClient *LoggingHTTPClient
}

// NewClient creates a ContabilidadeClient using the given authenticated httpClient and session.
// Signature mirrors celcoin-sdk domain constructors: NewXxx(httpClient, session).
func NewClient(httpClient *http.Client, session Session) ContabilidadeClient {
	return &Client{
		session:    session,
		httpClient: NewLoggingHTTPClient(httpClient),
	}
}
