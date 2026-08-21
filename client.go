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

	// BaixarDanfse downloads the DANFSe PDF of an NFS-e issued via Portal Nacional
	// (POST /api/nfse/danfse/{chave}). certificadoBytes is the raw PFX/P12 content.
	BaixarDanfse(ctx context.Context, chaveAcesso string, certificadoBytes []byte, senha string) ([]byte, error)

	// SolicitarConsulta starts an async notes consultation for a period
	// (POST /api/consulta-notas/solicitar) and returns a tracking token.
	SolicitarConsulta(ctx context.Context, req SolicitarConsultaNotasRequest) (*SolicitarConsultaNotasResponse, error)

	// ConsultarStatus polls whether a notes consultation finished
	// (GET /api/consulta-notas/{token}/status).
	ConsultarStatus(ctx context.Context, token string) (*ConsultaNotasStatusResponse, error)

	// BaixarNotas lists notes found by a consultation (with XmlBase64)
	// (GET /api/consulta-notas/{token}/notas).
	BaixarNotas(ctx context.Context, token string) (*ListarNotasConsultadasResponse, error)

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
