# Dev Roadmap

## Status das Integrações

| Integração | Status | Documentação | Código | Observações |
|---|---|---|---|---|
| PNCP | ✅ Integrado | `docs/clientes/pncp.md` | `internal/sources/pncp/client/` | |
| Senado Federal | ✅ Integrado | `docs/clientes/senado-federal.md` | `internal/sources/senado/client/` | |
| Câmara dos Deputados | ✅ Integrado | `docs/clientes/camara-dos-deputados.md` | `internal/sources/deputados/client/` | |
| IBGE | ✅ Integrado | `docs/clientes/ibge.md` | `internal/sources/ibge/client/` | |
| TCU | ✅ Integrado | `docs/clientes/tcu.md` | `internal/sources/tcu/client/` | |
| Portal da Transparência | ✅ Integrado | `docs/clientes/portal-da-transparencia.md` | `internal/sources/portaltransparencia/client/` | Requer chave de API |
| OpenCNPJ | ✅ Integrado | `docs/clientes/opencnpj.md` | `internal/sources/opencnpj/client/` | |
| SICONFI | ⚠️ Problema | `docs/clientes/siconfi.md` | `internal/sources/siconfi/client/` | Erro de permissão em endpoints públicos |
| TSE (Prestação de Contas / Candidatos) | ⚠️ Problema | `docs/tse-importacao.md` | `internal/sources/tse/` | Candidato eleito não tem vínculo direto com município/UF |
| SERPRO — Consulta Dívida Ativa | 📋 Pendente | — | — | |
| SIOP | 📋 Pendente | — | — | |
| BNDES | 📋 Pendente | — | — | |
| DATAJUD | 📋 Pendente | — | — | |
| CVM — Fundos de Investimento | 📋 Pendente | — | — | |
| Banco de Sanções (CEIS/CNEP) | 🔬 Estudo | — | — | |

## Problemas Conhecidos

### SICONFI — Erro de permissão em endpoints públicos

A API do SICONFI é documentada como acesso público sem necessidade de credenciais, porém retorna erro de permissão (403) em algumas requisições.

**Referência:** https://apidatalake.tesouro.gov.br/docs/siconfi/

### TSE — Candidato eleito sem vínculo direto com município/UF

Os dados de candidatos do TSE não registram de forma direta a qual município ou UF o candidato eleito está vinculado, dificultando a conexão entre os dados eleitorais e as demais esferas (municipal, estadual, federal).

**Referências:**
- https://dadosabertos.tse.jus.br/dataset/candidatos-2024
- https://dadosabertos.tse.jus.br/dataset/prestacao-de-contas-eleitorais-2024
