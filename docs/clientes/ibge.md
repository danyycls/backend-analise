# IBGE

**Nome cliente:** IBGE

**Descrição:** Cliente Go para integração com as APIs de Dados do IBGE, fornecendo acesso a localidades (estados e municípios) e agregados (estimativas populacionais).

## Doc Client

**Documentação de integração client:** https://servicodados.ibge.gov.br/api/docs/
**Base URL:** https://servicodados.ibge.gov.br/api/v1/localidades (Localidades)
https://servicodados.ibge.gov.br/api/v3/agregados (Agregados)

## APIs Integradas

### Localidades

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarEstados | `/estados` | — | `[]types.EstadoIBGE` | Obtém o conjunto de unidades federativas do Brasil. |
| ListarMunicipios | `/estados/{uf}/municipios` | `uf string` | `[]types.MunicipioIBGE` | Obtém o conjunto de municípios de uma UF. |
| ListarMunicipiosCompleto | `/municipios` | — | `[]types.MunicipioDetalhadoIBGE` | Obtém o conjunto de municípios do Brasil com detalhes (microrregião/mesorregião/UF). |

### Agregados

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarPopulacao | `/6579/periods/-6/variaveis/9324` | `municipioIDs []int` | `map[int]int64` | Estimativas populacionais dos municípios (agregado 6579, variável 9324). Em caso de falha, faz fallback para o agregado 8395 (Censo 2022), período 2022, variável 12494. |
