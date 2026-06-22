# PNCP

**Nome cliente:** PNCP

**Descrição:** Cliente Go para integração com a API de Consulta do Portal Nacional de Contratações Públicas (PNCP), fornecendo acesso a contratos e contratações publicadas.

## Doc Client

**Documentação de integração client:** https://pncp.gov.br/api/consulta/swagger-ui/index.html#/
https://www.gov.br/pncp/pt-br/acesso-a-informacao/dados-abertos
**Base URL:** https://pncp.gov.br/pncp-consulta/v1

## APIs Integradas

### Contratos

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarContratos | `/contratos` | `cnpj string`, `dataInicial string`, `dataFinal string`, `pagina int`, `tamanho int` | `[]Contrato` | Consultar Contratos por Data de Publicação. |

### Contratações (Publicação)

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarContratacoesPorMunicipio | `/contratacoes/publicacao` | `codigoMunicipio string`, `dataInicial string`, `dataFinal string`, `codigoModalidade string`, `pagina int`, `tamanho int` | `*PublicacaoResponse` | Consultar Contratações por Data de Publicação (filtro por município IBGE). |
| BuscarContratacoesPorUF | `/contratacoes/publicacao` | `uf string`, `dataInicial string`, `dataFinal string`, `codigoModalidade string`, `pagina int`, `tamanho int` | `*PublicacaoResponse` | Consultar Contratações por Data de Publicação (filtro por UF). |
