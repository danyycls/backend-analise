# Mapeamento de Rotas — PODP (LICEU API)

## Introdução

### APIs Internas

| API | Tipo | Descrição | Doc |
|-----|------|-----------|-----|
| **TSE (CSV → PostgreSQL)** | Interna (dados importados) | Dados eleitorais do TSE importados de CSVs para PostgreSQL relacional | [`docs/db-tse.md`](./db-tse.md) · [`docs/tse-importacao.md`](./tse-importacao.md) |

### APIs Externas Consultadas

| API | Tipo | Base URL (configurável via env) | Doc |
|-----|------|----------------------------------|-----|
| **Câmara dos Deputados** | REST | `DEPUTADOS_BASE_URL` | [`docs/clientes/camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| **Senado Federal** | REST | `SENADO_BASE_URL` | [`docs/clientes/senado-federal.md`](./clientes/senado-federal.md) |
| **TCU — Tribunal de Contas da União** | REST | `TCU_BASE_URL` | [`docs/clientes/tcu.md`](./clientes/tcu.md) |
| **PNCP — Portal Nacional de Contratações Públicas** | REST | `PNCP_BASE_URL` | [`docs/clientes/pncp.md`](./clientes/pncp.md) |
| **Portal da Transparência** | REST | `PORTAL_TRANSPARENCIA_BASE_URL` + `PORTAL_TRANSPARENCIA_API_KEY` | [`docs/clientes/portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| **IBGE — Instituto Brasileiro de Geografia e Estatística** | REST | `IBGE_BASE_URL` + `IBGE_AGREGADOS_BASE_URL` | [`docs/clientes/ibge.md`](./clientes/ibge.md) |
| **OpenCNPJ** | REST | `OPENCNPJ_BASE_URL` | [`docs/clientes/opencnpj.md`](./clientes/opencnpj.md) |
| **SICONFI — Sistema de Informações Contábeis e Fiscais** | REST | `SICONFI_BASE_URL` | [`docs/clientes/siconfi.md`](./clientes/siconfi.md) |

**Total: 1 API interna + 8 APIs externas = 9 fontes de dados**

---

## Rotas de Consulta Dados TSE

Base de dados PostgreSQL populada via importação de CSVs do TSE (2006-2024).

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 1 | Importar CSV | POST | `/import` | PostgreSQL (TSE) | Importa planilhas CSV do TSE para o banco relacional | [`tse-importacao.md`](./tse-importacao.md) |
| 2 | Listar Cargos | GET | `/busca/cargos` | PostgreSQL (TSE) | Lista cargos eletivos disponíveis | [`db-tse.md`](./db-tse.md) |
| 3 | Listar Partidos | GET | `/busca/partidos` | PostgreSQL (TSE) | Lista partidos políticos | [`db-tse.md`](./db-tse.md) |
| 4 | Buscar Candidatos | POST | `/busca/candidatos` | PostgreSQL (TSE) | Consulta candidatos com filtros | [`db-tse.md`](./db-tse.md) |
| 5 | Buscar Doadores | POST | `/busca/doadores` | PostgreSQL (TSE) | Consulta doadores de campanha | [`db-tse.md`](./db-tse.md) |
| 6 | Buscar Fornecedores | POST | `/busca/fornecedores` | PostgreSQL (TSE) | Consulta fornecedores de campanha | [`db-tse.md`](./db-tse.md) |
| 7 | Buscar Relações | POST | `/busca/relacoes` | PostgreSQL (TSE) | Busca relações entre entidades eleitorais | [`db-tse.md`](./db-tse.md) |
| 8 | Consultar Entidade | POST | `/entidade` | PostgreSQL (TSE) | Consulta detalhes de uma entidade | [`db-tse.md`](./db-tse.md) |

---

## Ligação Política

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 9 | Analisar Contexto | POST | `/busca/contexto` | OpenCNPJ, TCU | Análise de ligação política entre entidades (sócio, parentesco, etc.) | [`ligacao-politica.md`](./ligacao-politica.md) |

---

## PNCP — Análise de Órgãos e Publicações

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 10 | Analisar Órgão PNCP | POST | `/orgao/analise` | PNCP, OpenCNPJ, Redis | Análise completa de um órgão no PNCP (contratos, CNPJ) | [`pncp.md`](./clientes/pncp.md) |
| 11 | Stream Órgão | GET | `/orgao/analise/stream/:jobId` | — (WebSocket) | Resultados em tempo real da análise de órgão | [`pncp.md`](./clientes/pncp.md) |
| 12 | Batch Órgão | GET | `/orgao/analise/batch/:jobId` | PNCP, OpenCNPJ, Redis | Resultados consolidados em lote | [`pncp.md`](./clientes/pncp.md) |
| 13 | Analisar Publicação | POST | `/publicacao/analise` | PNCP, OpenCNPJ, Redis | Análise de publicações no PNCP | [`pncp.md`](./clientes/pncp.md) |
| 14 | Stream Publicação | GET | `/publicacao/analise/stream/:jobId` | — (WebSocket) | Resultados em tempo real da análise de publicação | [`pncp.md`](./clientes/pncp.md) |
| 15 | Batch Publicação | GET | `/publicacao/analise/batch/:jobId` | PNCP, OpenCNPJ, Redis | Resultados consolidados em lote | [`pncp.md`](./clientes/pncp.md) |
| 16 | Listar Municípios | GET | `/ibge/municipios/:uf` | IBGE | Lista municípios de uma UF | [`ibge.md`](./clientes/ibge.md) |

---

## TCU — Tribunal de Contas da União

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 17 | Contas Irregulares | POST | `/tcu/contas-irregulares` | TCU | Consulta contas julgadas irregulares | [`tcu.md`](./clientes/tcu.md) |
| 18 | Fins Eleitorais | POST | `/tcu/fins-eleitorais` | TCU | Consulta dados com fins eleitorais | [`tcu.md`](./clientes/tcu.md) |
| 19 | Inabilitados | POST | `/tcu/inabilitados` | TCU | Consulta inabilitados para cargos | [`tcu.md`](./clientes/tcu.md) |
| 20 | Inidôneos | POST | `/tcu/inidoneos` | TCU | Consulta empresas inidôneas | [`tcu.md`](./clientes/tcu.md) |

---

## Deputados (Câmara dos Deputados)

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 21 | Deputados Ativos | GET | `/deputados` | Câmara dos Deputados | Lista todos deputados em exercício | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 22 | Detalhes Deputado | GET | `/deputados/:id/completo` | Câmara dos Deputados | Detalhes completos de um deputado | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 23 | Despesas Deputado | GET | `/deputados/:id/despesas` | Câmara dos Deputados | Despesas parlamentares de um deputado | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 24 | Órgãos Deputado | GET | `/deputados/:id/orgaos` | Câmara dos Deputados | Órgãos associados ao deputado | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |

---

## Senado Federal

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 25 | Listar Senadores | GET | `/senado/senadores` | Senado Federal | Lista todos os senadores | [`senado-federal.md`](./clientes/senado-federal.md) |
| 26 | Detalhes Senador | GET | `/senado/senadores/:codigo/completo` | Senado Federal | Detalhes completos de um senador | [`senado-federal.md`](./clientes/senado-federal.md) |
| 27 | Cargos Senador | GET | `/senado/senadores/:codigo/cargos` | Senado Federal | Cargos exercidos pelo senador | [`senado-federal.md`](./clientes/senado-federal.md) |
| 28 | Comissões Senador | GET | `/senado/senadores/:codigo/comissoes` | Senado Federal | Comissões que o senador participa | [`senado-federal.md`](./clientes/senado-federal.md) |
| 29 | Mandatos Senador | GET | `/senado/senadores/:codigo/mandatos` | Senado Federal | Mandatos do senador | [`senado-federal.md`](./clientes/senado-federal.md) |
| 30 | Orçamento | GET | `/senado/orcamento` | Senado Federal | Dados orçamentários do Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 31 | Processos | GET | `/senado/processors` | Senado Federal | Lista processos legislativos | [`senado-federal.md`](./clientes/senado-federal.md) |
| 32 | Assuntos Processo | GET | `/senado/processo/assuntos` | Senado Federal | Assuntos dos processos | [`senado-federal.md`](./clientes/senado-federal.md) |
| 33 | Emendas Processo | GET | `/senado/processo/emendas` | Senado Federal | Emendas dos processos | [`senado-federal.md`](./clientes/senado-federal.md) |
| 34 | Detalhes Processo | GET | `/senado/processo/:id` | Senado Federal | Detalhes de um processo específico | [`senado-federal.md`](./clientes/senado-federal.md) |
| 35 | Votações | GET | `/senado/votacoes` | Senado Federal | Lista votações do plenário | [`senado-federal.md`](./clientes/senado-federal.md) |
| 36 | Votações Comissão | GET | `/senado/votacoes/comissao/:sigla` | Senado Federal | Votações em comissão por sigla | [`senado-federal.md`](./clientes/senado-federal.md) |
| 37 | Votações Parlamentar | GET | `/senado/votacoes/parlamentar/:codigo` | Senado Federal | Votações por parlamentar | [`senado-federal.md`](./clientes/senado-federal.md) |
| 38 | Tramitação Matéria | GET | `/senado/materia/tramitacao` | Senado Federal | Tramitação de matérias legislativas | [`senado-federal.md`](./clientes/senado-federal.md) |
| 39 | Agenda Dia | GET | `/senado/agenda/dia/:data` | Senado Federal | Agenda do dia no Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 40 | Agenda Mês | GET | `/senado/agenda/mes/:data` | Senado Federal | Agenda do mês no Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 41 | Encontro | GET | `/senado/encontro/:codigo` | Senado Federal | Detalhes de um encontro | [`senado-federal.md`](./clientes/senado-federal.md) |
| 42 | Comissões | GET | `/senado/comissoes` | Senado Federal | Lista todas comissões do Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 43 | Detalhes Comissão | GET | `/senado/comissoes/:codigo` | Senado Federal | Detalhes de uma comissão | [`senado-federal.md`](./clientes/senado-federal.md) |

---

## Estados e IBGE

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 44 | Listar Estados | GET | `/ibge/estados` | IBGE | Lista todos os estados brasileiros | [`ibge.md`](./clientes/ibge.md) |
| 45 | Dados Completos Estado | GET | `/estado/:uf/dados-completos` | IBGE, Câmara Deputados, Senado Federal, DB (TSE) | Dados consolidados do estado (geral + políticos) | [`ibge.md`](./clientes/ibge.md) · [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) · [`senado-federal.md`](./clientes/senado-federal.md) · [`db-tse.md`](./db-tse.md) |
| 46 | Dados Básicos Estado | GET | `/estado/:uf/basico` | IBGE | Dados básicos do estado (população, área, etc.) | [`ibge.md`](./clientes/ibge.md) |
| 47 | Candidatos Estado | GET | `/estado/:uf/candidatos` | PostgreSQL (TSE) | Candidatos do estado por UF | [`db-tse.md`](./db-tse.md) |
| 48 | Deputados Estado | GET | `/estado/:uf/deputados` | Câmara dos Deputados | Deputados federais do estado | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 49 | Senadores Estado | GET | `/estado/:uf/senadores` | Senado Federal | Senadores do estado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 50 | Municípios/População | GET | `/ibge/municipios-populacao/:uf` | IBGE | Municípios com dados populacionais | [`ibge.md`](./clientes/ibge.md) |
| 51 | Financeiro Estado | GET | `/estado/:uf/financeiro/stream` | SICONFI, IBGE, Redis, Portal Transparência | Dados financeiros do estado via WebSocket | [`siconfi.md`](./clientes/siconfi.md) · [`ibge.md`](./clientes/ibge.md) · [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

---

## Município

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 52 | Detalhes Município | GET | `/municipio/:codigoIBGE/detalhes/stream` | SICONFI, PNCP | Detalhes do município via WebSocket | [`siconfi.md`](./clientes/siconfi.md) · [`pncp.md`](./clientes/pncp.md) · [`detalhes-municipio.md`](./detalhes-municipio.md) |

---

## Portal da Transparência

### Órgãos

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 53 | Órgãos SIAPE | GET | `/portal-transparencia/orgaos/siape` | Portal Transparência | Consulta órgãos por código SIAPE | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 54 | Órgãos SIAFI | GET | `/portal-transparencia/orgaos/siafi` | Portal Transparência | Consulta órgãos por código SIAFI | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Pessoas

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 55 | Pessoa Física | GET | `/portal-transparencia/pessoas/fisica` | Portal Transparência | Consulta pessoa física | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 56 | Pessoa Jurídica | GET | `/portal-transparencia/pessoas/juridica` | Portal Transparência | Consulta pessoa jurídica | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Cartões Corporativos

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 57 | Cartões | GET | `/portal-transparencia/cartoes` | Portal Transparência | Consulta gastos com cartões corporativos | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Servidores

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 58 | Servidores | GET | `/portal-transparencia/servidores` | Portal Transparência | Lista servidores públicos | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 59 | Servidor por ID | GET | `/portal-transparencia/servidores/:id` | Portal Transparência | Detalhes de um servidor | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 60 | Remuneração | GET | `/portal-transparencia/servidores/remuneracao` | Portal Transparência | Remuneração de servidores | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 61 | Servidores por Órgão | GET | `/portal-transparencia/servidores/por-orgao` | Portal Transparência | Servidores por órgão | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 62 | Funções e Cargos | GET | `/portal-transparencia/servidores/funcoes-e-cargos` | Portal Transparência | Funções e cargos comissionados | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 63 | PEPs | GET | `/portal-transparencia/servidores/peps` | Portal Transparência | Programas de integridade (PEPs) | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Despesas

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 64 | Tipos de Transferência | GET | `/portal-transparencia/despesas/tipo-transferencia` | Portal Transparência | Tipos de transferência de recursos | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 65 | Recursos Recebidos | GET | `/portal-transparencia/despesas/recursos-recebidos` | Portal Transparência | Recursos recebidos por ente | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 66 | Despesas por Órgão | GET | `/portal-transparencia/despesas/por-orgao` | Portal Transparência | Despesas por órgão | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 67 | Funcional Programática | GET | `/portal-transparencia/despesas/por-funcional-programatica` | Portal Transparência | Despesas por classificação funcional programática | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 68 | Movimentação Líquida | GET | `/portal-transparencia/despesas/por-funcional-programatica/movimentacao-liquida` | Portal Transparência | Movimentação líquida da despesa | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 69 | Plano Orçamentário | GET | `/portal-transparencia/despesas/plano-orcamentario` | Portal Transparência | Plano orçamentário | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 70 | Itens de Empenho | GET | `/portal-transparencia/despesas/itens-de-empenho` | Portal Transparência | Itens de empenho | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 71 | Histórico Empenho | GET | `/portal-transparencia/despesas/itens-de-empenho/historico` | Portal Transparência | Histórico de itens de empenho | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 72 | Subfunções | GET | `/portal-transparencia/despesas/funcional-programatica/subfuncoes` | Portal Transparência | Subfunções da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 73 | Programas | GET | `/portal-transparencia/despesas/funcional-programatica/programs` | Portal Transparência | Programas da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 74 | Listar Funcional Programática | GET | `/portal-transparencia/despesas/funcional-programatica/listar` | Portal Transparência | Lista classificação funcional programática | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 75 | Funções | GET | `/portal-transparencia/despesas/funcional-programatica/funcoes` | Portal Transparência | Funções da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 76 | Ações | GET | `/portal-transparencia/despesas/funcional-programatica/acoes` | Portal Transparência | Ações da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 77 | Favorecidos Finais | GET | `/portal-transparencia/despesas/favorecidos-finais-por-documento` | Portal Transparência | Favorecidos finais por documento | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 78 | Empenhos Impactados | GET | `/portal-transparencia/despesas/empenhos-impactados` | Portal Transparência | Empenhos impactados | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 79 | Documentos | GET | `/portal-transparencia/despesas/documentos` | Portal Transparência | Documentos de despesa | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 80 | Documento por Código | GET | `/portal-transparencia/despesas/documentos/:codigo` | Portal Transparência | Documento de despesa por código | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 81 | Documentos Relacionados | GET | `/portal-transparencia/despesas/documentos-relacionados` | Portal Transparência | Documentos relacionados | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 82 | Documentos por Favorecido | GET | `/portal-transparencia/despesas/documentos-por-favorecido` | Portal Transparência | Documentos por favorecido | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Emendas

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 83 | Emendas | GET | `/portal-transparencia/emendas` | Portal Transparência | Consulta emendas parlamentares | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 84 | Documentos Emenda | GET | `/portal-transparencia/emendas/documentos/:codigo` | Portal Transparência | Documentos de emenda parlamentar | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

---

## Resumo

| Seção | Qtde Rotas | API/Doc principal |
|-------|-----------|-------------------|
| TSE (consulta) | 8 | PostgreSQL (TSE) |
| Ligação Política | 1 | OpenCNPJ + TCU |
| PNCP | 7 | PNCP + OpenCNPJ |
| TCU | 4 | TCU |
| Deputados | 4 | Câmara dos Deputados |
| Senado Federal | 19 | Senado Federal |
| Estados / IBGE | 8 | IBGE + múltiplas |
| Município | 1 | SICONFI + PNCP |
| Portal Transparência — Órgãos | 2 | Portal Transparência |
| Portal Transparência — Pessoas | 2 | Portal Transparência |
| Portal Transparência — Cartões | 1 | Portal Transparência |
| Portal Transparência — Servidores | 6 | Portal Transparência |
| Portal Transparência — Despesas | 19 | Portal Transparência |
| Portal Transparência — Emendas | 2 | Portal Transparência |
| **Total** | **84** | **9 fontes de dados** |

---

> **Legenda:** Rotas definidas em `internal/app/routes.go` (Gin framework).  
> **WebSocket:** rotas `/stream/:jobId` e `/detalhes/stream` utilizam WebSocket (`gorilla/websocket`) para streaming de dados.  
> **Cache:** Redis utilizado nas rotas de PNCP, Ligação Política, Financeiro Estado e Deputados.
