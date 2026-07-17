# SICONFI

**Nome cliente:** SICONFI

**Descrição:** Cliente Go para integração com a API de Dados Abertos do SICONFI (Sistema de Informações Contábeis e Fiscais do Setor Público Brasileiro) da Secretaria do Tesouro Nacional, fornecendo acesso a entes, DCA, RGF, RREO, matrizes de saldos contábeis (MSC), extrato de entregas e anexos de relatórios.

## Doc Client

**Documentação de integração client:** https://apidatalake.tesouro.gov.br/docs/siconfi/
**Base URL:** https://apidatalake.tesouro.gov.br/ords/siconfi/tt

### Restrições da API

- **Rate limit:** 1 requisição por segundo.
- **Paginação padrão:** 5000 itens por página.
- **Tratamento de erros:** A API pode retornar erro de conta bloqueada (`AccountIsLocked`), que é tratado pelo cliente com o erro `ErrSICONFIIndisponivel`.

## APIs Integradas

### Entes

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarEntes | `/entes` | — | `[]Ente` | Informações básicas de cadastro dos entes da federação. |

### DCA

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarDCA | `/dca` | `anExercicio int64`, `idEnte int`, `noAnexo ...string` | `[]DCAItem` | Lista dos dados contidos nos quadros das contas anuais. |

### RGF

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarRGF | `/rgf` | `params RGFParams` | `[]RGFItem` | Lista dos dados contidos nos anexos do RGF. |

### RREO

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarRREO | `/rreo` | `params RREOParams` | `[]RREOItem` | Lista dos dados contidos nos anexos do RREO. |

### MSC

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarMSCPatrimonial | `/msc_patrimonial` | `params MSCParams` | `[]MSCItem` | Detalhamento dos registros informados nas contas contábeis que recebem lançamentos de natureza patrimonial. |
| BuscarMSCOrcamentaria | `/msc_orcamentaria` | `params MSCParams` | `[]MSCItem` | Detalhamento dos registros informados nas contas contábeis que recebem lançamentos de natureza orçamentária. |
| BuscarMSCControle | `/msc_controle` | `params MSCParams` | `[]MSCItem` | Detalhamento dos registros informados nas contas contábeis que recebem lançamentos de natureza de controle. |

### Extrato de Entregas

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BuscarExtratoEntregas | `/extrato_entregas` | `idEnte int`, `anReferencia int64` | `[]ExtratoEntregasItem` | Extrato contendo informações sobre relatórios homologados e retificados e matrizes entregues. |

### Anexos

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarAnexosRelatorios | `/anexos-relatorios` | — | `[]AnexoRelatorio` | Tabela de apoio dos anexos dos relatórios. |
