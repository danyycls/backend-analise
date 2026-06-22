# Dev Roadmap

## Status das Integrações

| Integração | Status | Documentação | Código | Observações |
|---|---|---|---|---|
| PNCP | ✅ Integrado | `docs/clientes/pncp.md` | `internal/shared/clients/pncp/` | |
| Senado Federal | ✅ Integrado | `docs/clientes/senado-federal.md` | `internal/shared/clients/senado/` | |
| Câmara dos Deputados | ✅ Integrado | `docs/clientes/camara-dos-deputados.md` | `internal/shared/clients/deputados/` | |
| IBGE | ✅ Integrado | `docs/clientes/ibge.md` | `internal/shared/clients/ibge/` | |
| TCU | ✅ Integrado | `docs/clientes/tcu.md` | `internal/shared/clients/tcu/` | |
| Portal da Transparência | ✅ Integrado | `docs/clientes/portal-da-transparencia.md` | `internal/shared/clients/portaltransparencia/` | Requer chave de API |
| SICONFI | ⚠️ Problema | `docs/clientes/siconfi.md` | `internal/shared/clients/siconfi/` | Erro de permissão em endpoints públicos |
| TSE (Prestação de Contas / Candidatos) | ⚠️ Problema | `docs/tse-importacao.md` | `internal/esferas-brasileiras/tse/` | Candidato eleito não tem vínculo direto com município/UF |
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
