package contabilidade_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
				ChaveAcesso:       contabilidade.String("35260620000000000000065000000000100000000001"),
			})
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/nfse/danfse/"):
			w.Header().Set("Content-Type", "application/pdf")
			_, _ = w.Write([]byte("%PDF-1.4 mock-danfse"))
		case r.Method == http.MethodPost && r.URL.Path == "/api/consulta-notas/solicitar":
			_ = json.NewEncoder(w).Encode(contabilidade.SolicitarConsultaNotasResponse{
				Sucesso: true,
				Token:   contabilidade.String("11111111-2222-3333-4444-555555555555"),
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/consulta-notas/11111111-2222-3333-4444-555555555555/status":
			_ = json.NewEncoder(w).Encode(contabilidade.ConsultaNotasStatusResponse{
				Sucesso:    true,
				Token:      contabilidade.String("11111111-2222-3333-4444-555555555555"),
				Finalizado: true,
				Status:     contabilidade.String(contabilidade.ConsultaNotasStatusFinalizado),
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/consulta-notas/11111111-2222-3333-4444-555555555555/notas":
			_ = json.NewEncoder(w).Encode(contabilidade.ListarNotasConsultadasResponse{
				Sucesso:        true,
				Token:          contabilidade.String("11111111-2222-3333-4444-555555555555"),
				StatusConsulta: contabilidade.String(contabilidade.ConsultaNotasStatusFinalizado),
				Notas: []contabilidade.NotaConsultadaDto{
					{
						Nro:       contabilidade.String("230"),
						Chave:     contabilidade.String("35260620000000000000065000000000100000000001"),
						Status:    contabilidade.String(contabilidade.NotaConsultadaStatusEmitida),
						Direcao:   contabilidade.String(contabilidade.NotaConsultadaDirecaoEmitida),
						XmlBase64: contabilidade.String("PGttbD5tb2NrPC9rbWw+"),
					},
				},
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
			CodigoServico:               contabilidade.String("02881"),
			CodigoTributacaoNacional:    contabilidade.String("171901"),
			CodigoTributacaoMunicipal:   contabilidade.String("001"),
			DiscriminacaoServico: contabilidade.String("Consultoria e desenvolvimento de sistemas"),
			DataEmissao:          &now,
		},
	})

	s.Require().NoError(err)
	s.Equal(int64(230), resp.Numero)
	s.Equal("LRJGC6HA", *resp.CodigoVerificacao)
	s.Require().NotNil(resp.ChaveAcesso)
	s.Equal("35260620000000000000065000000000100000000001", *resp.ChaveAcesso)
	s.True(resp.Sucesso)
}

func (s *ClientSuite) TestBaixarDanfse() {
	pdf, err := s.client.BaixarDanfse(
		context.Background(),
		"35260620000000000000065000000000100000000001",
		[]byte("fake-pfx-bytes"),
		"senha123",
	)
	s.Require().NoError(err)
	s.True(strings.HasPrefix(string(pdf), "%PDF-1.4"))
}

func (s *ClientSuite) TestBaixarDanfse_errorBody() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`E1235 - Falha no esquema XML do DF-e.`))
	}))
	defer server.Close()

	cfg := contabilidade.Config{
		APIEndpoint: contabilidade.String(server.URL),
		BearerToken: contabilidade.String("test-token"),
		Timeout:     func() *time.Duration { d := 5 * time.Second; return &d }(),
	}
	session, err := contabilidade.NewSession(cfg)
	s.Require().NoError(err)
	client := contabilidade.NewClient(contabilidade.CreateBearerHTTPClient(session), *session)

	_, err = client.BaixarDanfse(context.Background(), "chave", []byte("pfx"), "senha")
	s.Require().Error(err)
	ce, ok := contabilidade.ParseErr(err)
	s.True(ok)
	s.Equal("BAIXAR_DANFSE_ERROR", ce.ErrorKey)
	s.Contains(ce.Error(), "E1235")
}

func (s *ClientSuite) TestStatus() {
	err := s.client.Status(context.Background())
	s.NoError(err)
}

func (s *ClientSuite) TestSolicitarConsulta() {
	inicio := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	fim := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)

	resp, err := s.client.SolicitarConsulta(context.Background(), contabilidade.SolicitarConsultaNotasRequest{
		Certificado: &contabilidade.CertificadoDto{
			EmissorDocumento: contabilidade.String("56929454000104"),
			EmissorBase64:    contabilidade.String("MIIKJAIBAzCCCeA..."),
			EmissorSenha:     contabilidade.String("senha123"),
		},
		MunicipioCodigoIBGE: contabilidade.String("3550308"),
		InscricaoMunicipal:  contabilidade.String("14708477"),
		DataInicio:          &inicio,
		DataFim:             &fim,
	})

	s.Require().NoError(err)
	s.True(resp.Sucesso)
	s.Require().NotNil(resp.Token)
	s.Equal("11111111-2222-3333-4444-555555555555", *resp.Token)
}

func (s *ClientSuite) TestConsultarStatus() {
	resp, err := s.client.ConsultarStatus(context.Background(), "11111111-2222-3333-4444-555555555555")
	s.Require().NoError(err)
	s.True(resp.Sucesso)
	s.True(resp.Finalizado)
	s.Require().NotNil(resp.Status)
	s.Equal(contabilidade.ConsultaNotasStatusFinalizado, *resp.Status)
}

func (s *ClientSuite) TestBaixarNotas() {
	resp, err := s.client.BaixarNotas(context.Background(), "11111111-2222-3333-4444-555555555555")
	s.Require().NoError(err)
	s.True(resp.Sucesso)
	s.Require().Len(resp.Notas, 1)
	s.Equal("230", *resp.Notas[0].Nro)
	s.Equal("PGttbD5tb2NrPC9rbWw+", *resp.Notas[0].XmlBase64)
}

func (s *ClientSuite) TestConsultarStatus_emptyToken() {
	_, err := s.client.ConsultarStatus(context.Background(), "  ")
	s.Require().Error(err)
	ce, ok := contabilidade.ParseErr(err)
	s.True(ok)
	s.Equal("INVALID_REQUEST", ce.ErrorKey)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
