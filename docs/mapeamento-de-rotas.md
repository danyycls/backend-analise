# ODT — Mapeamento de Rotas

## Introdução

### APIs Internas

| API | Tipo | Descrição | Doc |
|-----|------|-----------|-----|
| **TSE (CSV → PostgreSQL)** | Interna (dados importados) | Dados eleitorais do TSE importados de CSVs para PostgreSQL relacional | [`docs/db-tse.md`](./db-tse.md) · [`docs/tse-importacao.md`](./tse-importacao.md) |
| **TSE Repositório** | Interna (consulta direta) | Consultas ao banco TSE: candidatos, fornecedores, doadores, receitas, despesas, relações | [`docs/db-tse.md`](./db-tse.md) |
| **Convênios (CSV → PostgreSQL)** | Interna (dados importados) | Convênios do Portal da Transparência importados de CSVs | [`docs/clientes/portal-transparencia-convenios-importacao.md`](./clientes/portal-transparencia-convenios-importacao.md) |

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

**Total: 3 APIs internas + 8 APIs externas = 11 fontes de dados**

---

## Rotas de Importação

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 1 | Importar CSV | POST | `/import` | PostgreSQL (TSE + Convênios) | Importa planilhas CSV do TSE e convênios para o banco relacional | [`tse-importacao.md`](./tse-importacao.md) |

---

## Rotas de Consulta Dados TSE

Base de dados PostgreSQL populada via importação de CSVs do TSE (2006-2024).

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 2 | Listar Cargos | GET | `/busca/cargos` | PostgreSQL (TSE) | Lista cargos eletivos disponíveis | [`db-tse.md`](./db-tse.md) |
| 3 | Listar Partidos | GET | `/busca/partidos` | PostgreSQL (TSE) | Lista partidos políticos | [`db-tse.md`](./db-tse.md) |
| 4 | Buscar Candidatos | POST | `/busca/candidatos` | PostgreSQL (TSE) | Consulta candidatos com filtros | [`db-tse.md`](./db-tse.md) |
| 5 | Buscar Doadores | POST | `/busca/doadores` | PostgreSQL (TSE) | Consulta doadores de campanha | [`db-tse.md`](./db-tse.md) |
| 6 | Buscar Fornecedores | POST | `/busca/fornecedores` | PostgreSQL (TSE) | Consulta fornecedores de campanha | [`db-tse.md`](./db-tse.md) |
| 7 | Buscar Relações | POST | `/busca/relacoes` | PostgreSQL (TSE) | Busca relações entre entidades eleitorais | [`db-tse.md`](./db-tse.md) |
| 8 | Consultar Entidade | POST | `/entidade` | PostgreSQL (TSE) | Consulta detalhes de uma entidade | [`db-tse.md`](./db-tse.md) |

---

## TSE Repositório — Consultas Diretas ao Banco

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 9 | Cargos Distintos | GET | `/tse/repositorio/cargos-distintos` | PostgreSQL (TSE) | Lista cargos distintos dos candidatos | [`db-tse.md`](./db-tse.md) |
| 10 | Buscar Candidatos | POST | `/tse/repositorio/candidatos` | PostgreSQL (TSE) | Busca candidatos por filtros avançados | [`db-tse.md`](./db-tse.md) |
| 11 | Buscar Candidato por CPF | POST | `/tse/repositorio/candidato/cpf` | PostgreSQL (TSE) | Busca candidato pelo CPF | [`db-tse.md`](./db-tse.md) |
| 12 | Buscar Candidato por ID | POST | `/tse/repositorio/candidato/id` | PostgreSQL (TSE) | Busca candidato pelo ID interno | [`db-tse.md`](./db-tse.md) |
| 13 | Buscar Fornecedores | POST | `/tse/repositorio/fornecedores/documento` | PostgreSQL (TSE) | Busca fornecedores por documento | [`db-tse.md`](./db-tse.md) |
| 14 | Buscar Doadores | POST | `/tse/repositorio/doadores/documento` | PostgreSQL (TSE) | Busca doadores por documento | [`db-tse.md`](./db-tse.md) |
| 15 | Receitas Candidato | POST | `/tse/repositorio/receitas-candidato` | PostgreSQL (TSE) | Busca receitas de um candidato por doador | [`db-tse.md`](./db-tse.md) |
| 16 | Receitas Partido | POST | `/tse/repositorio/receitas-partido` | PostgreSQL (TSE) | Busca receitas de um partido por doador | [`db-tse.md`](./db-tse.md) |
| 17 | Despesas Candidato | POST | `/tse/repositorio/despesas-candidato` | PostgreSQL (TSE) | Busca despesas de um candidato por fornecedor | [`db-tse.md`](./db-tse.md) |
| 18 | Despesas Partido | POST | `/tse/repositorio/despesas-partido` | PostgreSQL (TSE) | Busca despesas de um partido por fornecedor | [`db-tse.md`](./db-tse.md) |
| 19 | Buscar Partidos | POST | `/tse/repositorio/partidos` | PostgreSQL (TSE) | Busca partidos por IDs | [`db-tse.md`](./db-tse.md) |
| 20 | Buscar Eleições | POST | `/tse/repositorio/eleicoes` | PostgreSQL (TSE) | Busca eleições por IDs | [`db-tse.md`](./db-tse.md) |
| 21 | Candidatos Eleitos | POST | `/tse/repositorio/candidatos-eleitos` | PostgreSQL (TSE) | Lista candidatos eleitos por UF | [`db-tse.md`](./db-tse.md) |
| 22 | Buscar Relações | POST | `/tse/repositorio/relacoes` | PostgreSQL (TSE) | Busca relações entre entidades | [`db-tse.md`](./db-tse.md) |

---

## Análise PNCP — Órgãos e Municípios

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 23 | Analisar Órgão PNCP | POST | `/orgao/analise` | PNCP, OpenCNPJ, Redis | Análise completa de um órgão no PNCP | [`pncp.md`](./clientes/pncp.md) |
| 24 | Batch Órgão | GET | `/orgao/analise/batch/:jobId` | PNCP, OpenCNPJ, Redis | Resultados consolidados em lote | [`pncp.md`](./clientes/pncp.md) |
| 25 | Analisar UF/Município | POST | `/uf-municipio/analise` | PNCP, IBGE, OpenCNPJ | Análise de licitações por UF/município | [`pncp.md`](./clientes/pncp.md) |
| 26 | Batch UF/Município | GET | `/uf-municipio/analise/batch/:jobId` | PNCP, IBGE, OpenCNPJ | Resultados consolidados em lote | [`pncp.md`](./clientes/pncp.md) |

---

## TCU — Tribunal de Contas da União

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 27 | Contas Irregulares | POST | `/tcu/contas-irregulares` | TCU | Consulta contas julgadas irregulares | [`tcu.md`](./clientes/tcu.md) |
| 28 | Fins Eleitorais | POST | `/tcu/fins-eleitorais` | TCU | Consulta dados com fins eleitorais | [`tcu.md`](./clientes/tcu.md) |
| 29 | Inabilitados | POST | `/tcu/inabilitados` | TCU | Consulta inabilitados para cargos | [`tcu.md`](./clientes/tcu.md) |
| 30 | Inidôneos | POST | `/tcu/inidoneos` | TCU | Consulta empresas inidôneas | [`tcu.md`](./clientes/tcu.md) |

---

## Deputados (Câmara dos Deputados)

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 31 | Deputados Ativos | GET | `/deputados` | Câmara dos Deputados | Lista todos deputados em exercício | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 32 | Detalhes Deputado | GET | `/deputados/:id/completo` | Câmara dos Deputados | Detalhes completos de um deputado | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 33 | Despesas Deputado | GET | `/deputados/:id/despesas` | Câmara dos Deputados | Despesas parlamentares de um deputado | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 34 | Órgãos Deputado | GET | `/deputados/:id/orgaos` | Câmara dos Deputados | Órgãos associados ao deputado | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 35 | Listar Partidos | GET | `/deputados/partidos` | Câmara dos Deputados | Lista partidos políticos | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 36 | Buscar Partido | GET | `/deputados/partidos/:id` | Câmara dos Deputados | Detalhes de um partido | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 37 | Membros Partido | GET | `/deputados/partidos/:id/membros` | Câmara dos Deputados | Membros de um partido | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 38 | Proposições | GET | `/deputados/proposicoes` | Câmara dos Deputados | Lista proposições legislativas | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 39 | Buscar Proposição | GET | `/deputados/proposicoes/:id` | Câmara dos Deputados | Detalhes de uma proposição | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 40 | Tramitações | GET | `/deputados/proposicoes/:id/tramitacoes` | Câmara dos Deputados | Tramitações de uma proposição | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 41 | Autores | GET | `/deputados/proposicoes/:id/autores` | Câmara dos Deputados | Autores de uma proposição | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 42 | Temas | GET | `/deputados/proposicoes/:id/temas` | Câmara dos Deputados | Temas de uma proposição | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 43 | Relacionadas | GET | `/deputados/proposicoes/:id/relacionadas` | Câmara dos Deputados | Proposições relacionadas | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 44 | Eventos | GET | `/deputados/eventos` | Câmara dos Deputados | Lista eventos | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 45 | Buscar Evento | GET | `/deputados/eventos/:id` | Câmara dos Deputados | Detalhes de um evento | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 46 | Órgãos Câmara | GET | `/deputados/orgaos` | Câmara dos Deputados | Lista órgãos da Câmara | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 47 | Buscar Órgão Câmara | GET | `/deputados/orgaos/:id` | Câmara dos Deputados | Detalhes de um órgão | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 48 | Membros Órgão | GET | `/deputados/orgaos/:id/membros` | Câmara dos Deputados | Membros de um órgão | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 49 | Blocos | GET | `/deputados/blocos` | Câmara dos Deputados | Lista blocos partidários | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 50 | Buscar Bloco | GET | `/deputados/blocos/:id` | Câmara dos Deputados | Detalhes de um bloco | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 51 | Partidos do Bloco | GET | `/deputados/blocos/:id/partidos` | Câmara dos Deputados | Partidos de um bloco | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 52 | Frentes | GET | `/deputados/frentes` | Câmara dos Deputados | Lista frentes parlamentares | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 53 | Buscar Frente | GET | `/deputados/frentes/:id` | Câmara dos Deputados | Detalhes de uma frente | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 54 | Membros Frente | GET | `/deputados/frentes/:id/membros` | Câmara dos Deputados | Membros de uma frente | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 55 | Grupos | GET | `/deputados/grupos` | Câmara dos Deputados | Lista grupos parlamentares | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 56 | Buscar Grupo | GET | `/deputados/grupos/:id` | Câmara dos Deputados | Detalhes de um grupo | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 57 | Legislaturas | GET | `/deputados/legislaturas` | Câmara dos Deputados | Lista legislaturas | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 58 | Buscar Legislatura | GET | `/deputados/legislaturas/:id` | Câmara dos Deputados | Detalhes de uma legislatura | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 59 | Votações | GET | `/deputados/votacoes` | Câmara dos Deputados | Lista votações | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 60 | Buscar Votação | GET | `/deputados/votacoes/:id` | Câmara dos Deputados | Detalhes de uma votação | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |
| 61 | Votos | GET | `/deputados/votacoes/:id/votos` | Câmara dos Deputados | Votos individuais em uma votação | [`camara-dos-deputados.md`](./clientes/camara-dos-deputados.md) |

---

## Senado Federal

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 62 | Listar Senadores | GET | `/senado/senadores` | Senado Federal | Lista todos os senadores | [`senado-federal.md`](./clientes/senado-federal.md) |
| 63 | Detalhes Senador | GET | `/senado/senadores/:codigo/completo` | Senado Federal | Detalhes completos de um senador | [`senado-federal.md`](./clientes/senado-federal.md) |
| 64 | Cargos Senador | GET | `/senado/senadores/:codigo/cargos` | Senado Federal | Cargos exercidos pelo senador | [`senado-federal.md`](./clientes/senado-federal.md) |
| 65 | Comissões Senador | GET | `/senado/senadores/:codigo/comissoes` | Senado Federal | Comissões que o senador participa | [`senado-federal.md`](./clientes/senado-federal.md) |
| 66 | Mandatos Senador | GET | `/senado/senadores/:codigo/mandatos` | Senado Federal | Mandatos do senador | [`senado-federal.md`](./clientes/senado-federal.md) |
| 67 | Orçamento | GET | `/senado/orcamento` | Senado Federal | Dados orçamentários do Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 68 | Processos | GET | `/senado/processors` | Senado Federal | Lista processos legislativos | [`senado-federal.md`](./clientes/senado-federal.md) |
| 69 | Assuntos Processo | GET | `/senado/processo/assuntos` | Senado Federal | Assuntos dos processos | [`senado-federal.md`](./clientes/senado-federal.md) |
| 70 | Emendas Processo | GET | `/senado/processo/emendas` | Senado Federal | Emendas dos processos | [`senado-federal.md`](./clientes/senado-federal.md) |
| 71 | Detalhes Processo | GET | `/senado/processo/:id` | Senado Federal | Detalhes de um processo específico | [`senado-federal.md`](./clientes/senado-federal.md) |
| 72 | Votações | GET | `/senado/votacoes` | Senado Federal | Lista votações do plenário | [`senado-federal.md`](./clientes/senado-federal.md) |
| 73 | Votações Comissão | GET | `/senado/votacoes/comissao/:sigla` | Senado Federal | Votações em comissão por sigla | [`senado-federal.md`](./clientes/senado-federal.md) |
| 74 | Votações Parlamentar | GET | `/senado/votacoes/parlamentar/:codigo` | Senado Federal | Votações por parlamentar | [`senado-federal.md`](./clientes/senado-federal.md) |
| 75 | Tramitação Matéria | GET | `/senado/materia/tramitacao` | Senado Federal | Tramitação de matérias legislativas | [`senado-federal.md`](./clientes/senado-federal.md) |
| 76 | Agenda Dia | GET | `/senado/agenda/dia/:data` | Senado Federal | Agenda do dia no Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 77 | Agenda Mês | GET | `/senado/agenda/mes/:data` | Senado Federal | Agenda do mês no Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 78 | Encontro | GET | `/senado/encontro/:codigo` | Senado Federal | Detalhes de um encontro | [`senado-federal.md`](./clientes/senado-federal.md) |
| 79 | Comissões | GET | `/senado/comissoes` | Senado Federal | Lista todas comissões do Senado | [`senado-federal.md`](./clientes/senado-federal.md) |
| 80 | Detalhes Comissão | GET | `/senado/comissoes/:codigo` | Senado Federal | Detalhes de uma comissão | [`senado-federal.md`](./clientes/senado-federal.md) |
| 81 | Documento Emenda | GET | `/senado/emendas/:id/documento` | Senado Federal | Baixa documento de emenda | [`senado-federal.md`](./clientes/senado-federal.md) |

---

## Estados e IBGE

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 82 | Listar Estados | GET | `/ibge/estados` | IBGE | Lista todos os estados brasileiros | [`ibge.md`](./clientes/ibge.md) |
| 83 | Listar Municípios | GET | `/ibge/municipios/:uf` | IBGE | Lista municípios de uma UF | [`ibge.md`](./clientes/ibge.md) |
| 84 | Buscar População | POST | `/ibge/populacao` | IBGE | Estimativas populacionais dos municípios | [`ibge.md`](./clientes/ibge.md) |

---

## OpenCNPJ

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 85 | Buscar CNPJ | GET | `/opencnpj/:cnpj` | OpenCNPJ | Dados cadastrais de uma pessoa jurídica | [`opencnpj.md`](./clientes/opencnpj.md) |

---

## PNCP — Contratos

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 86 | Contratos por Município | POST | `/pncp/contratos/municipio/:codigoIBGE` | PNCP | Contratos por município IBGE | [`pncp.md`](./clientes/pncp.md) |
| 87 | Contratos por Órgão | POST | `/pncp/contratos/orgao/:cnpj` | PNCP | Contratos por órgão (CNPJ) | [`pncp.md`](./clientes/pncp.md) |
| 88 | Contratos por UF | POST | `/pncp/contratos/uf/:uf` | PNCP | Contratos por UF | [`pncp.md`](./clientes/pncp.md) |

---

## SICONFI — Sistema de Informações Contábeis e Fiscais

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 89 | Entes | POST | `/siconfi/entes` | SICONFI | Informações de entes da federação | [`siconfi.md`](./clientes/siconfi.md) |
| 90 | DCA | POST | `/siconfi/dca` | SICONFI | Dados Contábeis Anuais | [`siconfi.md`](./clientes/siconfi.md) |
| 91 | RGF | POST | `/siconfi/rgf` | SICONFI | Relatório de Gestão Fiscal | [`siconfi.md`](./clientes/siconfi.md) |
| 92 | RREO | POST | `/siconfi/rreo` | SICONFI | Relatório Resumido da Execução Orçamentária | [`siconfi.md`](./clientes/siconfi.md) |
| 93 | MSC Patrimonial | POST | `/siconfi/msc-patrimonial` | SICONFI | Matriz de Saldos Contábeis — Patrimonial | [`siconfi.md`](./clientes/siconfi.md) |
| 94 | MSC Orçamentária | POST | `/siconfi/msc-orcamentaria` | SICONFI | Matriz de Saldos Contábeis — Orçamentária | [`siconfi.md`](./clientes/siconfi.md) |
| 95 | MSC Controle | POST | `/siconfi/msc-controle` | SICONFI | Matriz de Saldos Contábeis — Controle | [`siconfi.md`](./clientes/siconfi.md) |
| 96 | Extrato Entregas | POST | `/siconfi/extrato-entregas` | SICONFI | Extrato de relatórios entregues | [`siconfi.md`](./clientes/siconfi.md) |
| 97 | Anexos Relatórios | GET | `/siconfi/anexos-relatorios` | SICONFI | Tabela de anexos dos relatórios | [`siconfi.md`](./clientes/siconfi.md) |

---

## Portal da Transparência

### Órgãos

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 98 | Órgãos SIAPE | GET | `/portal-transparencia/orgaos/siape` | Portal Transparência | Consulta órgãos por código SIAPE | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 99 | Órgãos SIAFI | GET | `/portal-transparencia/orgaos/siafi` | Portal Transparência | Consulta órgãos por código SIAFI | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Pessoas

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 100 | Pessoa Física | GET | `/portal-transparencia/pessoas/fisica` | Portal Transparência | Consulta pessoa física | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 101 | Pessoa Jurídica | GET | `/portal-transparencia/pessoas/juridica` | Portal Transparência | Consulta pessoa jurídica | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Cartões Corporativos

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 102 | Cartões | GET | `/portal-transparencia/cartoes` | Portal Transparência | Consulta gastos com cartões corporativos | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Servidores

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 103 | Servidores | GET | `/portal-transparencia/servidores` | Portal Transparência | Lista servidores públicos | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 104 | Servidor por ID | GET | `/portal-transparencia/servidores/:id` | Portal Transparência | Detalhes de um servidor | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 105 | Remuneração | GET | `/portal-transparencia/servidores/remuneracao` | Portal Transparência | Remuneração de servidores | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 106 | Servidores por Órgão | GET | `/portal-transparencia/servidores/por-orgao` | Portal Transparência | Servidores por órgão | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 107 | Funções e Cargos | GET | `/portal-transparencia/servidores/funcoes-e-cargos` | Portal Transparência | Funções e cargos comissionados | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 108 | PEPs | GET | `/portal-transparencia/servidores/peps` | Portal Transparência | Pessoas Expostas Politicamente | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Despesas

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 109 | Tipos de Transferência | GET | `/portal-transparencia/despesas/tipo-transferencia` | Portal Transparência | Tipos de transferência de recursos | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 110 | Recursos Recebidos | GET | `/portal-transparencia/despesas/recursos-recebidos` | Portal Transparência | Recursos recebidos por ente | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 111 | Despesas por Órgão | GET | `/portal-transparencia/despesas/por-orgao` | Portal Transparência | Despesas por órgão | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 112 | Funcional Programática | GET | `/portal-transparencia/despesas/por-funcional-programatica` | Portal Transparência | Despesas por classificação funcional programática | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 113 | Movimentação Líquida | GET | `/portal-transparencia/despesas/por-funcional-programatica/movimentacao-liquida` | Portal Transparência | Movimentação líquida da despesa | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 114 | Plano Orçamentário | GET | `/portal-transparencia/despesas/plano-orcamentario` | Portal Transparência | Plano orçamentário | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 115 | Itens de Empenho | GET | `/portal-transparencia/despesas/itens-de-empenho` | Portal Transparência | Itens de empenho | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 116 | Histórico Empenho | GET | `/portal-transparencia/despesas/itens-de-empenho/historico` | Portal Transparência | Histórico de itens de empenho | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 117 | Subfunções | GET | `/portal-transparencia/despesas/funcional-programatica/subfuncoes` | Portal Transparência | Subfunções da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 118 | Programas | GET | `/portal-transparencia/despesas/funcional-programatica/programs` | Portal Transparência | Programas da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 119 | Listar Funcional Programática | GET | `/portal-transparencia/despesas/funcional-programatica/listar` | Portal Transparência | Lista classificação funcional programática | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 120 | Funções | GET | `/portal-transparencia/despesas/funcional-programatica/funcoes` | Portal Transparência | Funções da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 121 | Ações | GET | `/portal-transparencia/despesas/funcional-programatica/acoes` | Portal Transparência | Ações da classificação | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 122 | Favorecidos Finais | GET | `/portal-transparencia/despesas/favorecidos-finais-por-documento` | Portal Transparência | Favorecidos finais por documento | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 123 | Empenhos Impactados | GET | `/portal-transparencia/despesas/empenhos-impactados` | Portal Transparência | Empenhos impactados | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 124 | Documentos | GET | `/portal-transparencia/despesas/documentos` | Portal Transparência | Documentos de despesa | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 125 | Documento por Código | GET | `/portal-transparencia/despesas/documentos/:codigo` | Portal Transparência | Documento de despesa por código | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 126 | Documentos Relacionados | GET | `/portal-transparencia/despesas/documentos-relacionados` | Portal Transparência | Documentos relacionados | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 127 | Documentos por Favorecido | GET | `/portal-transparencia/despesas/documentos-por-favorecido` | Portal Transparência | Documentos por favorecido | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

### Emendas

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 128 | Emendas | GET | `/portal-transparencia/emendas` | Portal Transparência | Consulta emendas parlamentares | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |
| 129 | Documentos Emenda | GET | `/portal-transparencia/emendas/documentos/:codigo` | Portal Transparência | Documentos de emenda parlamentar | [`portal-da-transparencia.md`](./clientes/portal-da-transparencia.md) |

---

## Convênios

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 130 | Listar Convênios | GET | `/convenios` | PostgreSQL (Convênios) | Lista convênios importados | [`portal-transparencia-convenios-importacao.md`](./clientes/portal-transparencia-convenios-importacao.md) |

---

## WebSocket

| # | Nome | Método | URL | Clients consultados | Descrição | Doc |
|---|------|--------|-----|---------------------|-----------|-----|
| 131 | WebSocket Hub | GET | `/ws` | — | Canal WebSocket para streaming de análises PNCP | — |

---

## Resumo

| Seção | Qtde Rotas | API/Doc principal |
|-------|-----------|-------------------|
| Importação | 1 | PostgreSQL (TSE + Convênios) |
| TSE (consulta) | 7 | PostgreSQL (TSE) |
| TSE Repositório | 14 | PostgreSQL (TSE) |
| Análise PNCP | 4 | PNCP + OpenCNPJ |
| TCU | 4 | TCU |
| Deputados | 31 | Câmara dos Deputados |
| Senado Federal | 20 | Senado Federal |
| Estados / IBGE | 3 | IBGE |
| OpenCNPJ | 1 | OpenCNPJ |
| PNCP Contratos | 3 | PNCP |
| SICONFI | 9 | SICONFI |
| Portal Transparência — Órgãos | 2 | Portal Transparência |
| Portal Transparência — Pessoas | 2 | Portal Transparência |
| Portal Transparência — Cartões | 1 | Portal Transparência |
| Portal Transparência — Servidores | 6 | Portal Transparência |
| Portal Transparência — Despesas | 19 | Portal Transparência |
| Portal Transparência — Emendas | 2 | Portal Transparência |
| Convênios | 1 | PostgreSQL (Convênios) |
| WebSocket | 1 | — |
| **Total** | **131** | **11 fontes de dados** |

---

> **Legenda:** Rotas definidas em `internal/app/routes.go` (Gin framework).
> **WebSocket:** rota `/ws` utilize WebSocket (`gorilla/websocket`) para streaming de dados de análises PNCP.
> **Cache:** Redis utilizado nas rotas de PNCP, Análise Órgão, Análise UF/Município e Deputados.
