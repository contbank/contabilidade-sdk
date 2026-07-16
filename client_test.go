package contabilidade_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	contabilidade "github.com/contbank/contabilidade-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
	server  *httptest.Server
	client  contabilidade.ContabilidadeClient
	session *contabilidade.Session
}

func (s *ClientSuite) SetupTest() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "Bearer test-token", r.Header.Get("Authorization"))

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/nfse/emitir":
			_ = json.NewEncoder(w).Encode(contabilidade.EmitirNfseResponse{
				Sucesso:           true,
				Numero:            230,
				CodigoVerificacao: contabilidade.String("LRJGC6HA"),
				Link:              contabilidade.String("https://nfe.prefeitura.sp.gov.br/contribuinte/notaprint.aspx?nf=230&c=LRJGC6HA"),
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/status":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	cfg := contabilidade.Config{
		APIEndpoint: contabilidade.String(s.server.URL),
		BearerToken: contabilidade.String("test-token"),
		Timeout:     func() *time.Duration { d := 5 * time.Second; return &d }(),
	}

	session, err := contabilidade.NewSession(cfg)
	s.Require().NoError(err)
	s.session = session

	httpClient := contabilidade.CreateBearerHTTPClient(session)
	s.client = contabilidade.NewClient(httpClient, *session)
}

func (s *ClientSuite) TearDownTest() {
	s.server.Close()
}

func (s *ClientSuite) TestEmitirNfse() {
	now := time.Now()
	resp, err := s.client.EmitirNfse(context.Background(), contabilidade.EmitirNfseRequest{
		Certificado: &contabilidade.CertificadoDto{
			EmissorDocumento: contabilidade.String("56929454000104"),
			EmissorBase64:    contabilidade.String("MIIKJAIBAzCCCeA..."),
			EmissorSenha:     contabilidade.String("senha123"),
		},
		Nota: &contabilidade.NotaFiscalInput{
			RPS: &contabilidade.RpsInput{
				Serie:  contabilidade.String("FT26"),
				Tipo:   contabilidade.String("RPS"),
				Numero: contabilidade.String("139"),
			},
			Emissor: &contabilidade.EmissorInput{
				CNPJ:               contabilidade.String("56929454000104"),
				RazaoSocial:        contabilidade.String("GN SOFTWARE E CONSULTING LTDA"),
				InscricaoMunicipal: contabilidade.String("14708477"),
				CodigoPostal:       contabilidade.String("04203050"),
				Logradouro:         contabilidade.String("R BOM PASTOR"),
				LogradouroNumero:   contabilidade.String("529"),
				Bairro:             contabilidade.String("IPIRANGA"),
				MunicipioCodigoIBGE: contabilidade.String("3550308"),
				MunicipioSiglaUF:   contabilidade.String("SP"),
			},
			Tomador: &contabilidade.TomadorInput{
				PessoaFisica:    false,
				Documento:       contabilidade.String("10245240000100"),
				RazaoSocialNome: contabilidade.String("CONTABILIDADE.COM LTDA"),
			},
			Valores: &contabilidade.ValoresInput{
				ValorServico: 1500,
				Aliquota:     2,
			},
			CodigoServico:        contabilidade.String("02881"),
			DiscriminacaoServico: contabilidade.String("Consultoria e desenvolvimento de sistemas"),
			DataEmissao:          &now,
		},
	})

	s.Require().NoError(err)
	s.Equal(int64(230), resp.Numero)
	s.Equal("LRJGC6HA", *resp.CodigoVerificacao)
	s.True(resp.Sucesso)
}

func (s *ClientSuite) TestStatus() {
	err := s.client.Status(context.Background())
	s.NoError(err)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
