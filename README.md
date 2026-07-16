# contabilidade-sdk

SDK Go para o gateway NFS-e da Contabilidade.com (`https://nfse.contabilidade.com`).

Segue a mesma filosofia do `celcoin-sdk`: package flat, `Config` → `NewSession` → HTTP autenticado → `NewClient`.

## Autenticação

Swagger: `securitySchemes.Bearer` — token estático via header `Authorization: Bearer <token>`.

## Uso (DI no container do microserviço)

```go
import contabilidade "github.com/contbank/contabilidade-sdk"

cfg := contabilidade.Config{
    APIEndpoint: settings.Contabilidade.APIEndpoint,
    BearerToken: settings.Contabilidade.BearerToken,
}
session, err := contabilidade.NewSession(cfg)
httpClient := contabilidade.CreateBearerHTTPClient(session)
client := contabilidade.NewClient(httpClient, *session)

resp, err := client.EmitirNfse(ctx, contabilidade.EmitirNfseRequest{ /* ... */ })
```

Ver `config.example.yaml` para a seção YAML.
