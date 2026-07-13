# PODP — Projeto Observatório de Dados Públicos (Data Hub)

**P1 do ecossistema PODP.** Data hub central que consolida dados públicos brasileiros: fontes externas (Câmara, Senado, Portal da Transparência, TCU, PNCP, IBGE, SICONFI, OpenCNPJ), APIs REST e dados importados do TSE. Fornece dados brutos e processados para o **P2 (Motor de Análise)**.

> Este é o **P1 (backend-analise)**. O motor de análise de negócio (ligação política, anomalias, esferas brasileiras) está no **[P2 (projeto2-analise)](https://github.com/danyele/podp-analise)**.

---

## Arquitetura

```
Frontend → P1 (backend-analise :8080) → fontes externas (IBGE, TSE, Câmara, etc.)
         → P2 (projeto2-analise :8084) → P1 (via HTTP) + TCU/PortalTransparência (direto)
```

| Projeto | Papel | Porta |
|---------|-------|-------|
| **P1** (este) | Data hub: fontes, APIs REST, clients externos, dados TSE importados | `8080` |
| **P2** | Motor de análise: ligação política, anomalias, esferas brasileiras, WebSocket | `8084` |
| **Frontend** | Interface web | `8081` |

### O que está no P2

| Serviço | Rotas P2 |
|---------|----------|
| Ligação Política | `POST /busca/contexto` |
| Anomalias (worker) | `POST /worker/anomalia/iniciar`, `POST /worker/anomalia/parar/:jobId`, `GET /worker/anomalia/progression/:jobId` |
| Anomalias (consulta) | `GET /anomalias` |
| Feedback | `POST /feedback` |
| Estados (dados completos) | `GET /estado/:uf/basico`, `/dados-completos`, `/candidatos`, `/deputados`, `/senadores` |
| Municípios | `GET /municipio/:codigoIBGE/detalhes/stream` |
| Financeiro | `GET /estado/:uf/financeiro/stream` |
| WebSocket Hub | `GET /ws` (canal `anomalia_analise`) |

---

## Integrações

### APIs Externas (REST)

| Serviço | Descrição | Doc |
|---------|-----------|-----|
| **Câmara dos Deputados** | Deputados, legislaturas, votações, frentes parlamentares | [`docs/clientes/camara-dos-deputados.md`](docs/clientes/camara-dos-deputados.md) |
| **Senado Federal** | Senadores, comissões, votações, processos, orçamento, agenda | [`docs/clientes/senado-federal.md`](docs/clientes/senado-federal.md) |
| **TCU — Tribunal de Contas da União** | Contas irregulares, inabilitados, inidôneos, fins eleitorais | [`docs/clientes/tcu.md`](docs/clientes/tcu.md) |
| **Portal da Transparência** | Órgãos (SIAPE/SIAFI), pessoas, servidores, cartões, despesas, emendas | [`docs/clientes/portal-da-transparencia.md`](docs/clientes/portal-da-transparencia.md) |
| **PNCP — Portal Nacional de Contratações Públicas** | Contratos, publicações, análise de órgãos | [`docs/clientes/pncp.md`](docs/clientes/pncp.md) |
| **OpenCNPJ** | Dados cadastrais de CNPJ (razão social, sócios, situação) | [`docs/clientes/opencnpj.md`](docs/clientes/opencnpj.md) |
| **IBGE** | Estados, municípios, população, dados geográficos | [`docs/clientes/ibge.md`](docs/clientes/ibge.md) |
| **SICONFI — Sistema de Informações Contábeis e Fiscais** | DCA, RGF, RREO, MSC, entes | [`docs/clientes/siconfi.md`](docs/clientes/siconfi.md) |

### Dados Importados (CSV → PostgreSQL)

| Fonte | Descrição | Doc |
|-------|-----------|-----|
| **TSE — Tribunal Superior Eleitoral** | Dados eleitorais históricos (2006–2024): candidatos, partidos, cargos, doadores, fornecedores, prestação de contas. Importados de planilhas CSV para PostgreSQL relacional | [`docs/db-tse.md`](docs/db-tse.md) · [`docs/tse-importacao.md`](docs/tse-importacao.md) |

---

## Como Rodar o Projeto

```sh
# 1. Pré-requisitos: Go 1.25+, Docker Compose

# 2. Suba as dependências (PostgreSQL 15 + Redis 7)
docker compose up -d

# 3. Copie e configure o .env (preencha as variáveis com as credenciais)
cp .env.example .env

# 4. Execute a aplicação
go run .
```

A aplicação inicia na porta `8080` por padrão. O Swagger UI fica disponível em `http://localhost:8082` quando o docker-compose está em execução.

---

## Comandos Principais

```sh
make fix        # Formata código, organiza módulos (go mod tidy) e aplica autofix do lint
make test       # Executa todos os testes (unitários + integração, timeout 300s)
make test-unit  # Apenas testes unitários (60s, -short)
make lint       # Análise estática (golangci-lint)
```

Para ver todos os comandos disponíveis: `make help`.

---

## Arquitetura de Pastas

```
backend-analise/
├── main.go                          # Entrypoint
├── internal/
│   ├── app/                         # DI e registro de rotas
│   │   ├── app.go                   # Container de dependências
│   │   └── routes.go                # Todas as rotas HTTP (84 rotas)
│   ├── esferas-brasileiras/         # Módulos por esfera federativa
│   │   ├── tse/                     # Dados do TSE
│   │   │   ├── handler/             # Handlers HTTP
│   │   │   ├── usecase/             # Casos de uso
│   │   │   ├── repositorio/         # Acesso a banco
│   │   │   └── importacao/          # Importação de CSVs (worker pool, pgCOPY)
│   │   ├── federal/                 # Dados federais
│   │   │   ├── deputados/           # Câmara dos Deputados
│   │   │   ├── senadores/           # Senado Federal
│   │   │   ├── portaltransparencia/ # Portal da Transparência
│   │   │   ├── pncp/                # Contratações públicas
│   │   │   └── tcu/                 # TCU
│   │   ├── estadual/                # Dados estaduais
│   │   └── municipal/               # Dados municipais
│   └── shared/                      # Código compartilhado
│       ├── clients/                 # Clientes HTTP externos (8 APIs)
│       ├── database/                # Pool PostgreSQL
│       ├── migrations/              # Migrations SQL
│       ├── redis/                   # Cache Redis
│       ├── testkit/                 # Helpers de teste
│       ├── types/                   # Tipos compartilhados
│       └── websocket/               # WebSocket utilities
├── dataCSV/                         # CSVs do TSE por ano
├── api/                             # OpenAPI spec (openapi.yaml)
├── docs/                            # Documentação
│   ├── clientes/                    # Docs individuais de cada API
│   │   ├── camara-dos-deputados.md
│   │   ├── ibge.md
│   │   ├── opencnpj.md
│   │   ├── pncp.md
│   │   ├── portal-da-transparencia.md
│   │   ├── senado-federal.md
│   │   ├── siconfi.md
│   │   └── tcu.md
│   ├── dev-roadmap.md               # Roadmap de desenvolvimento
│   ├── db-tse.md                    # Arquitetura do banco TSE
│   ├── tse-importacao.md            # Processo de importação de CSVs
│   └── mapeamento-de-rotas.md       # Mapeamento completo de todas as rotas
├── docker-compose.yml               # PostgreSQL + Redis + Swagger
├── Dockerfile                       # Build multi-stage
└── Makefile                         # Comandos de build/test/lint
```

---

## Documentação

| Documento | Descrição |
|-----------|-----------|
| [`docs/mapeamento-de-rotas.md`](docs/mapeamento-de-rotas.md) | Lista completa de todas as rotas, organizadas por seção, com clients consultados e referências |
| [`docs/db-tse.md`](docs/db-tse.md) | Arquitetura do banco TSE: entidades, relacionamentos, FKs, índices, migrations |
| [`docs/tse-importacao.md`](docs/tse-importacao.md) | Processo de importação de CSVs do TSE: formato, pipeline, workers |
| [`docs/clientes/camara-dos-deputados.md`](docs/clientes/camara-dos-deputados.md) | API da Câmara dos Deputados |
| [`docs/clientes/senado-federal.md`](docs/clientes/senado-federal.md) | API do Senado Federal |
| [`docs/clientes/tcu.md`](docs/clientes/tcu.md) | API do TCU |
| [`docs/clientes/pncp.md`](docs/clientes/pncp.md) | API do PNCP |
| [`docs/clientes/portal-da-transparencia.md`](docs/clientes/portal-da-transparencia.md) | API do Portal da Transparência |
| [`docs/clientes/ibge.md`](docs/clientes/ibge.md) | API do IBGE |
| [`docs/clientes/opencnpj.md`](docs/clientes/opencnpj.md) | API do OpenCNPJ |
| [`docs/dev-roadmap.md`](docs/dev-roadmap.md) | Roadmap de desenvolvimento e status das integrações |
