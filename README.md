# PODP — Projeto Observatório de Dados Públicos

Painel unificado para análise de dados públicos brasileiros. Consolida informações eleitorais (TSE), parlamentares (Câmara, Senado), gastos públicos (Portal da Transparência, TCU), contratações (PNCP) e dados fiscais (SICONFI/IBGE) em uma única plataforma com APIs REST e WebSocket.

O projeto importa dados eleitorais do TSE (CSVs 2006–2024) para um banco PostgreSQL relacional e consulta APIs públicas oficiais em tempo real, permitindo cruzamentos como ligação política entre licitações/contratos e dados eleitorais, além de análises detalhadas de estados e municípios.

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

## Análises

### Ligação Política

Cruza documentos de licitações/contratos com dados eleitorais do **TSE** (fornecedores/doadores), enriquece com **OpenCNPJ** (razão social, sócios, situação cadastral) e sanções do **TCU** (contas irregulares, inidôneos, inabilitados). Utiliza **Redis** como cache. Depende dos dados do TSE persistidos em PostgreSQL.

- **Rota:** `POST /busca/contexto`
- **Doc:** [`docs/mapeamento-de-rotas.md#ligação-política`](docs/mapeamento-de-rotas.md#ligação-política)

### Detalhe Município

Consulta dados consolidados de um município a partir do código IBGE, combinando informações do **SICONFI** (dados contábeis/fiscais) e **PNCP** (contratações públicas). Resultados entregues via WebSocket para processamento assíncrono.

- **Rota:** `GET /municipio/:codigoIBGE/detalhes/stream` (WebSocket)
- **Doc:** [`docs/mapeamento-de-rotas.md#município`](docs/mapeamento-de-rotas.md#município)

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
│   ├── ligacao-politica/            # Análise de ligação política
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
│   ├── dev-roadmap.md               # Roadmap de desenvolvimento
│   ├── db-tse.md                    # Arquitetura do banco TSE
│   ├── tse-importacao.md            # Processo de importação de CSVs
│   └── mapeamento-de-rotas.md       # Mapeamento completo de todas as 84 rotas
├── docker-compose.yml               # PostgreSQL + Redis + Swagger
├── Dockerfile                       # Build multi-stage
└── Makefile                         # Comandos de build/test/lint
```

---

## Documentação

| Documento | Descrição |
|-----------|-----------|
| [`docs/mapeamento-de-rotas.md`](docs/mapeamento-de-rotas.md) | Lista completa de todas as 84 rotas, organizadas por seção, com clients consultados e referências |
| [`docs/db-tse.md`](docs/db-tse.md) | Arquitetura do banco TSE: entidades, relacionamentos, FKs, índices, migrations |
| [`docs/tse-importacao.md`](docs/tse-importacao.md) | Processo de importação de CSVs do TSE: formato, pipeline, workers |
| [`docs/clientes/camara-dos-deputados.md`](docs/clientes/camara-dos-deputados.md) | API da Câmara dos Deputados |
| [`docs/clientes/senado-federal.md`](docs/clientes/senado-federal.md) | API do Senado Federal |
| [`docs/clientes/tcu.md`](docs/clientes/tcu.md) | API do TCU |
| [`docs/clientes/pncp.md`](docs/clientes/pncp.md) | API do PNCP |
| [`docs/clientes/portal-da-transparencia.md`](docs/clientes/portal-da-transparencia.md) | API do Portal da Transparência |
| [`docs/clientes/ibge.md`](docs/clientes/ibge.md) | API do IBGE |
| [`docs/clientes/opencnpj.md`](docs/clientes/opencnpj.md) | API do OpenCNPJ |
| [`docs/clientes/siconfi.md`](docs/clientes/siconfi.md) | API do SICONFI |
| [`docs/dev-roadmap.md`](docs/dev-roadmap.md) | Roadmap de desenvolvimento e status das integrações |
