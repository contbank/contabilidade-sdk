package contabilidade

import "time"

// API path constants (Contabilidade.com NFS-e gateway).
const (
	PathEmitirNfse           = "/api/nfse/emitir"
	PathCancelarNfse         = "/api/nfse/cancelar"
	PathBaixarXmlNfse        = "/api/nfse/xml/%s/%d/%s" // im, numero, codigo
	PathBaixarDanfse         = "/api/nfse/danfse/%s"    // chave acesso (Portal Nacional)
	PathConsultarServicoCnae = "/api/servico/consultar-por-cnae"
	PathListarMunicipios     = "/api/servico/municipios"
	PathStatus               = "/api/status"
)

// --- Enums (Swagger) ---

// LocalIncidenciaEnum ...
type LocalIncidenciaEnum int32

const (
	LocalIncidenciaPrestador LocalIncidenciaEnum = 0
	LocalIncidenciaTomador   LocalIncidenciaEnum = 1
	LocalIncidenciaOutro     LocalIncidenciaEnum = 2
)

// OpcaoSimplesNFSeEnum ...
type OpcaoSimplesNFSeEnum int32

const (
	OpcaoSimplesNaoOptante              OpcaoSimplesNFSeEnum = 1
	OpcaoSimplesOptanteMEI              OpcaoSimplesNFSeEnum = 2
	OpcaoSimplesOptanteME_EPP           OpcaoSimplesNFSeEnum = 3
	OpcaoSimplesOptanteME_EPP_Sublimite OpcaoSimplesNFSeEnum = 4
	OpcaoSimplesMei                     OpcaoSimplesNFSeEnum = 5
	OpcaoSimplesOutros                  OpcaoSimplesNFSeEnum = 7
	OpcaoSimplesNaoInformado            OpcaoSimplesNFSeEnum = 99
)

// RegimeEspecialTributacaoEnum ...
type RegimeEspecialTributacaoEnum int32

// TipoTributacaoNFSeEnum ...
type TipoTributacaoNFSeEnum int32

const (
	TipoTributacaoNoMunicipio         TipoTributacaoNFSeEnum = 1
	TipoTributacaoForaMunicipio       TipoTributacaoNFSeEnum = 2
	TipoTributacaoIsento              TipoTributacaoNFSeEnum = 3
	TipoTributacaoImune               TipoTributacaoNFSeEnum = 4
	TipoTributacaoExigibilidadeSusp   TipoTributacaoNFSeEnum = 5
	TipoTributacaoExigibilidadeSuspJD TipoTributacaoNFSeEnum = 6
)

// --- Certificate ---

// CertificadoDto is the PFX certificate payload for PMSP digital signature.
type CertificadoDto struct {
	EmissorDocumento *string `json:"EmissorDocumento,omitempty"`
	EmissorBase64    *string `json:"EmissorBase64,omitempty"`
	EmissorSenha     *string `json:"EmissorSenha,omitempty"`
}

// --- Emitir NFS-e ---

// EmitirNfseRequest is the body for POST /api/nfse/emitir.
type EmitirNfseRequest struct {
	Certificado *CertificadoDto  `json:"Certificado,omitempty"`
	Nota        *NotaFiscalInput `json:"Nota,omitempty"`
}

// EmitirNfseResponse is the result of an emission attempt.
type EmitirNfseResponse struct {
	Sucesso           bool     `json:"Sucesso"`
	Numero            int64    `json:"Numero"`
	CodigoVerificacao *string  `json:"CodigoVerificacao,omitempty"`
	Link              *string  `json:"Link,omitempty"`
	PdfUrl            *string  `json:"PdfUrl,omitempty"`
	// ChaveAcesso chave de 44 dígitos retornada pelo Portal Nacional (fora de SP).
	// Necessária para download do DANFSe e cancelamento nacional.
	ChaveAcesso *string  `json:"ChaveAcesso,omitempty"`
	Erros       []string `json:"Erros,omitempty"`
}

// NotaFiscalInput holds the full NFS-e payload.
type NotaFiscalInput struct {
	RPS                            *RpsInput                     `json:"RPS,omitempty"`
	Emissor                        *EmissorInput                 `json:"Emissor,omitempty"`
	Tomador                        *TomadorInput                 `json:"Tomador,omitempty"`
	Valores                        *ValoresInput                 `json:"Valores,omitempty"`
	CodigoServico                  *string                       `json:"CodigoServico,omitempty"`
	CodigoLC116                    *string                       `json:"CodigoLC116,omitempty"`
	CodigoCNAEServico              *string                       `json:"CodigoCNAEServico,omitempty"`
	// CodigoTributacaoNacional (cTribNac) — obrigatório no Portal Nacional NFS-e.
	// Exemplos: "010104" (desenvolvimento de software), "010101" (análise/desenvolvimento de sistemas).
	CodigoTributacaoNacional       *string                       `json:"CodigoTributacaoNacional,omitempty"`
	DiscriminacaoServico           *string                       `json:"DiscriminacaoServico,omitempty"`
	DataEmissao                    *time.Time                    `json:"DataEmissao,omitempty"`
	DataFatoGerador                *time.Time                    `json:"DataFatoGerador,omitempty"`
	IDTipoTributacaoNFSeEnum       *TipoTributacaoNFSeEnum       `json:"IDTipoTributacaoNFSeEnum,omitempty"`
	IDOpcaoSimplesNFSeEnum         *OpcaoSimplesNFSeEnum         `json:"IDOpcaoSimplesNFSeEnum,omitempty"`
	IDRegimeEspecialTributacaoEnum *RegimeEspecialTributacaoEnum `json:"IDRegimeEspecialTributacaoEnum,omitempty"`
	ISSRetido                      bool                          `json:"ISSRetido"`
	ISSRetidoIntermediador         bool                          `json:"ISSRetidoIntermediador"`
	LocalIncidencia                *LocalIncidenciaEnum          `json:"LocalIncidencia,omitempty"`
	MunicipioCodigoIBGEPrestacao   *string                       `json:"MunicipioCodigoIBGEPrestacao,omitempty"`
	NumeroLote                     *string                       `json:"NumeroLote,omitempty"`
	NumeroNFSeSubstituicao         *int64                        `json:"NumeroNFSeSubstituicao,omitempty"`
	Intermediador                  *IntermediadorInput           `json:"Intermediador,omitempty"`
	Obra                           *ObraInput                    `json:"Obra,omitempty"`
	Guia                           *GuiaInput                    `json:"Guia,omitempty"`
	Tributacao                     *TributacaoInput              `json:"Tributacao,omitempty"`
	Observacao                     *ObservacaoInput              `json:"Observacao,omitempty"`
}

// RpsInput is the provisional service receipt identifier.
type RpsInput struct {
	Serie       *string    `json:"Serie,omitempty"`
	Tipo        *string    `json:"Tipo,omitempty"`
	Numero      *string    `json:"Numero,omitempty"`
	DataEmissao *time.Time `json:"DataEmissao,omitempty"`
}

// EmissorInput is the service provider (prestador).
type EmissorInput struct {
	CNPJ                 *string `json:"CNPJ,omitempty"`
	RazaoSocial          *string `json:"RazaoSocial,omitempty"`
	NomeFantasia         *string `json:"NomeFantasia,omitempty"`
	InscricaoMunicipal   *string `json:"InscricaoMunicipal,omitempty"`
	InscricaoEstadual    *string `json:"InscricaoEstadual,omitempty"`
	Email                *string `json:"Email,omitempty"`
	CodigoPostal         *string `json:"CodigoPostal,omitempty"`
	Logradouro           *string `json:"Logradouro,omitempty"`
	LogradouroTipo       *string `json:"LogradouroTipo,omitempty"`
	LogradouroComplemento *string `json:"LogradouroComplemento,omitempty"`
	LogradouroNumero     *string `json:"LogradouroNumero,omitempty"`
	Bairro               *string `json:"Bairro,omitempty"`
	MunicipioCodigoIBGE  *string `json:"MunicipioCodigoIBGE,omitempty"`
	MunicipioNome        *string `json:"MunicipioNome,omitempty"`
	MunicipioSiglaUF     *string `json:"MunicipioSiglaUF,omitempty"`
	Cnae                 *string `json:"Cnae,omitempty"`
}

// TomadorInput is the service taker (cliente).
type TomadorInput struct {
	Estrangeiro           bool    `json:"Estrangeiro"`
	PessoaFisica          bool    `json:"PessoaFisica"`
	Documento             *string `json:"Documento,omitempty"`
	RazaoSocialNome       *string `json:"RazaoSocialNome,omitempty"`
	NomeFantasiaApelido   *string `json:"NomeFantasiaApelido,omitempty"`
	InscricaoMunicipal    *string `json:"InscricaoMunicipal,omitempty"`
	InscricaoEstadual     *string `json:"InscricaoEstadual,omitempty"`
	Email                 *string `json:"Email,omitempty"`
	CodigoPostal          *string `json:"CodigoPostal,omitempty"`
	Logradouro            *string `json:"Logradouro,omitempty"`
	LogradouroTipo        *string `json:"LogradouroTipo,omitempty"`
	LogradouroComplemento *string `json:"LogradouroComplemento,omitempty"`
	LogradouroNumero      *string `json:"LogradouroNumero,omitempty"`
	Bairro                *string `json:"Bairro,omitempty"`
	MunicipioCodigoIBGE   *string `json:"MunicipioCodigoIBGE,omitempty"`
	MunicipioNome         *string `json:"MunicipioNome,omitempty"`
	MunicipioSiglaUF      *string `json:"MunicipioSiglaUF,omitempty"`
	CodigoPais            *string `json:"CodigoPais,omitempty"`
}

// ValoresInput holds financial values of the invoice.
type ValoresInput struct {
	ValorServico           float64  `json:"ValorServico"`
	Aliquota               float64  `json:"Aliquota"`
	ValorDeducoes          *float64 `json:"ValorDeducoes,omitempty"`
	DescontoIncondicionado *float64 `json:"DescontoIncondicionado,omitempty"`
	DescontoCondicionado   *float64 `json:"DescontoCondicionado,omitempty"`
	ValorISSRetido         *float64 `json:"ValorISSRetido,omitempty"`
	ValorPIS               *float64 `json:"ValorPIS,omitempty"`
	ValorCOFINS            *float64 `json:"ValorCOFINS,omitempty"`
	ValorCSLL              *float64 `json:"ValorCSLL,omitempty"`
	ValorINSS              *float64 `json:"ValorINSS,omitempty"`
	ValorIR                *float64 `json:"ValorIR,omitempty"`
	OutrasRetencoes        *float64 `json:"OutrasRetencoes,omitempty"`
	ValorCredito           *float64 `json:"ValorCredito,omitempty"`
}

// IntermediadorInput ...
type IntermediadorInput struct {
	PessoaFisica       bool    `json:"PessoaFisica"`
	Documento          *string `json:"Documento,omitempty"`
	RazaoSocialNome    *string `json:"RazaoSocialNome,omitempty"`
	InscricaoMunicipal *string `json:"InscricaoMunicipal,omitempty"`
	Email              *string `json:"Email,omitempty"`
}

// ObraInput ...
type ObraInput struct {
	CEI           *string `json:"CEI,omitempty"`
	Matricula     *string `json:"Matricula,omitempty"`
	Encapsulamento *string `json:"Encapsulamento,omitempty"`
}

// GuiaInput ...
type GuiaInput struct {
	Numero       *string    `json:"Numero,omitempty"`
	DataQuitacao *time.Time `json:"DataQuitacao,omitempty"`
}

// TributacaoInput ...
type TributacaoInput struct {
	CargaTributariaValor      *float64 `json:"CargaTributariaValor,omitempty"`
	CargaTributariaPercentual *float64 `json:"CargaTributariaPercentual,omitempty"`
	CargaTributariaFonte      *string  `json:"CargaTributariaFonte,omitempty"`
}

// ObservacaoInput ...
type ObservacaoInput struct {
	DiscriminacaoServico *string `json:"DiscriminacaoServico,omitempty"`
	OutrasInformacoes    *string `json:"OutrasInformacoes,omitempty"`
	IncentivadorCultural *bool   `json:"IncentivadorCultural,omitempty"`
}

// --- Cancelar NFS-e ---

// CancelarNfseRequest is the body for POST /api/nfse/cancelar.
type CancelarNfseRequest struct {
	Certificado        *CertificadoDto `json:"Certificado,omitempty"`
	Numero             int64           `json:"Numero"`
	CodigoVerificacao  *string         `json:"CodigoVerificacao,omitempty"`
	InscricaoMunicipal *string         `json:"InscricaoMunicipal,omitempty"`
}

// CancelarNfseResponse is the result of a cancellation attempt.
type CancelarNfseResponse struct {
	Sucesso           bool       `json:"Sucesso"`
	Numero            int64      `json:"Numero"`
	CodigoVerificacao *string    `json:"CodigoVerificacao,omitempty"`
	DataCancelamento  *time.Time `json:"DataCancelamento,omitempty"`
	Erros             []string   `json:"Erros,omitempty"`
}

// --- Serviço / CNAE ---

// ConsultarServicoPorCnaeRequest ...
type ConsultarServicoPorCnaeRequest struct {
	CodigoIBGEMunicipio *string  `json:"CodigoIBGEMunicipio,omitempty"`
	Cnaes               []string `json:"Cnaes,omitempty"`
}

// ConsultarServicoPorCnaeResponse ...
type ConsultarServicoPorCnaeResponse struct {
	CodigoIBGEMunicipio *string         `json:"CodigoIBGEMunicipio,omitempty"`
	NomeMunicipio       *string         `json:"NomeMunicipio,omitempty"`
	UF                  *string         `json:"UF,omitempty"`
	Resultados          []ResultadoCnae `json:"Resultados,omitempty"`
}

// ResultadoCnae ...
type ResultadoCnae struct {
	CodigoCnae     *string            `json:"CodigoCnae,omitempty"`
	DescricaoCnae  *string            `json:"DescricaoCnae,omitempty"`
	CodigosServico []CodigoServicoDto `json:"CodigosServico,omitempty"`
}

// CodigoServicoDto ...
type CodigoServicoDto struct {
	Codigo    *string `json:"Codigo,omitempty"`
	ItemLC116 *string `json:"ItemLC116,omitempty"`
	Descricao *string `json:"Descricao,omitempty"`
	Natureza  *string `json:"Natureza,omitempty"`
	Aliquota  float64 `json:"Aliquota"`
	EmiteNFSe bool    `json:"EmiteNFSe"`
}

// ProblemDetails follows RFC 7807-style error payloads from the gateway.
type ProblemDetails struct {
	Type     *string `json:"type,omitempty"`
	Title    *string `json:"title,omitempty"`
	Status   *int32  `json:"status,omitempty"`
	Detail   *string `json:"detail,omitempty"`
	Instance *string `json:"instance,omitempty"`
}
