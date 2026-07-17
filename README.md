# ODT: Observatório de Dados Públicos

ODT é um data hub que centraliza processos de extração, normalização e consulta a dados públicos brasileiros de múltiplas fontes.

É um projeto de estudo que propõe uma arquitetura de integração baseada em **Go** para os pipelines de extração e normalização, com **PostgreSQL** como camada de persistência — foco inicial nos datasets do TSE e do PNCP.

### ETL: Decisões técnicas

- **Demais portais** (exceto PNCP e TSE): a persistência não se mostrou necessária a curto prazo, pois as APIs suportam concorrência sem degradação significativa e os tempos de resposta se mantêm adequados mesmo em consultas extensas.

- **TSE:** a adoção de ETL foi necessária porque os dados são distribuídos em CSV, o que inviabiliza consultas relacionais e cruzamentos eficientes sem uma camada intermediária. Além disso, foram identificados IDs (SQ) inválidos — registros de prestação de contas referenciando candidatos e partidos inexistentes nas tabelas base, exigindo tratamento específico durante a importação.

- **PNCP:** embora o PNCP ofereça documentação e uma API REST para consulta, os testes em larga escala revelaram limitações que tornam a persistência local mais eficiente a longo prazo:
  - **Problemas observados:** throttling agressivo, falhas de conexão recorrentes (tanto no servidor web quanto no banco subjacente) e vazamento de mensagens de erro internas do banco pela API.
  - **Consulta por UF/Município:** o endpoint retorna licitações classificadas como DISPENSA por padrão. Apesar de juridicamente válido, esse comportamento é tecnicamente problemático, pois licitações que deveriam ser retornadas só aparecem quando a consulta é feita pelo CNPJ do órgão. Também foram observadas inconsistências nos resultados entre requisições idênticas nessa modalidade, comprometendo a confiabilidade dos dados. O mesmo não foi verificado na consulta por CNPJ do órgão.

  **Observação:** para sistemas que precisam de muitas requisições, consumir a API sob demanda não é eficiente — persistir os dados localmente e consultar contra o banco próprio é mais rápido e confiável a longo prazo.

---

### Fontes de Dados

| Tipo | Fonte | Descrição |
|------|-------|-----------|
| API REST | **Câmara dos Deputados** | Deputados, legislaturas, votações, frentes, proposições |
| API REST | **Senado Federal** | Senadores, comissões, votações, processos, orçamento, agenda |
| API REST | **TCU** | Contas irregulares, inabilitados, inidôneos, fins eleitorais |
| API REST | **Portal da Transparência** | Órgãos, servidores, despesas, cartões, emendas |
| API REST | **PNCP** | Contratos e contratações públicas |
| API REST | **OpenCNPJ** | Dados cadastrais de CNPJ |
| API REST | **IBGE** | Estados, municípios, população |
| API REST | **SICONFI** | DCA, RGF, RREO, MSC, entes |
| CSV Import | **TSE** | Dados eleitorais históricos (2006–2024) |
| CSV Import | **Portal da Transparência** | Convênios e acordos |

---

## Testando localmente

```sh
# 1. Pré-requisitos: Go 1.25+, Docker Compose

# 2. Copie e configure o .env (obtenha a API Key no Portal da Transparência)
cp .env.example .env

# 3. Suba as dependências (PostgreSQL 15 + Redis 7)
docker compose up -d

# 4. Execute a aplicação
go run .
```

A aplicação expõe a API na porta `8080`.

> O banco de dados do TSE é populado via rota `POST /import`. O throughput da importação é configurado por `IMPORT_MAX_WORKERS` no .env e varia conforme os recursos da máquina. Consulte a [doc de importação](docs/tse-importacao.md) para detalhes.

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

## Arquitetura

```
ODT/
├── main.go                          # Entrypoint
├── cmd/
│   └── migrate/                     # Migrations
├── internal/
│   ├── app/                         # DI e registro de rotas
│   │   ├── app.go                   # Container de dependências
│   │   └── routes.go                # Todas as rotas HTTP (131 rotas)
│   ├── api/                         # Handlers e usecases por domínio
│   │   ├── deputados/               # Câmara dos Deputados
│   │   ├── ibge/                    # IBGE
│   │   ├── opencnpj/                # OpenCNPJ
│   │   ├── pncp/                    # PNCP
│   │   ├── senado/                  # Senado Federal
│   │   ├── siconfi/                 # SICONFI
│   │   ├── tcu/                     # TCU
│   │   ├── tse/                     # Dados TSE (consulta + repositório)
│   │   └── portaltransparencia/     # Portal da Transparência
│   ├── sources/                     # Clients HTTP externos
│   │   ├── deputados/client/
│   │   ├── ibge/client/
│   │   ├── opencnpj/client/
│   │   ├── pncp/client/
│   │   ├── portaltransparencia/client/
│   │   ├── senado/client/
│   │   ├── siconfi/client/
│   │   ├── tcu/client/
│   │   └── tse/importacao/          # Importação de CSVs (worker pool, pgCOPY)
│   ├── shared/                      # Código compartilhado
│   │   ├── database/                # Pool PostgreSQL
│   │   ├── migrations/              # Migrations SQL
│   │   ├── redis/                   # Cache Redis
│   │   ├── types/                   # Tipos compartilhados
│   │   └── websocket/               # WebSocket utilities
│   └── stream/                      # WebSocket hub (streaming PNCP)
├── dataCSV/                         # CSVs do TSE por ano
├── api/                             # OpenAPI spec (openapi.yaml)
├── docs/                            # Documentação
│   ├── clientes/                    # Docs individuais de cada API
│   ├── db-tse.md                    # Arquitetura do banco TSE
│   ├── tse-importacao.md            # Processo de importação de CSVs
│   ├── mapeamento-de-rotas.md       # Mapeamento completo de todas as rotas
│   └── dev-roadmap.md               # Roadmap de desenvolvimento
├── docker-compose.yml               # PostgreSQL + Redis
├── Dockerfile                       # Build multi-stage
└── Makefile                         # Comandos de build/test/lint
```

---

## Documentação

| Documento | Descrição |
|-----------|-----------|
| [`api/openapi.yaml`](api/openapi.yaml) | Especificação OpenAPI (Swagger) da API |
| [`docs/mapeamento-de-rotas.md`](docs/mapeamento-de-rotas.md) | Lista completa de todas as 131 rotas, organizadas por seção |
| [`docs/db-tse.md`](docs/db-tse.md) | Arquitetura do banco TSE: entidades, relacionamentos, FKs, índices, migrations |
| [`docs/tse-importacao.md`](docs/tse-importacao.md) | Processo de importação de CSVs do TSE e convênios |
| [`docs/clientes/camara-dos-deputados.md`](docs/clientes/camara-dos-deputados.md) | API da Câmara dos Deputados |
| [`docs/clientes/senado-federal.md`](docs/clientes/senado-federal.md) | API do Senado Federal |
| [`docs/clientes/tcu.md`](docs/clientes/tcu.md) | API do TCU |
| [`docs/clientes/pncp.md`](docs/clientes/pncp.md) | API do PNCP |
| [`docs/clientes/portal-da-transparencia.md`](docs/clientes/portal-da-transparencia.md) | API do Portal da Transparência |
| [`docs/clientes/ibge.md`](docs/clientes/ibge.md) | API do IBGE |
| [`docs/clientes/opencnpj.md`](docs/clientes/opencnpj.md) | API do OpenCNPJ |
| [`docs/clientes/siconfi.md`](docs/clientes/siconfi.md) | API do SICONFI |
| [`docs/dev-roadmap.md`](docs/dev-roadmap.md) | Roadmap de desenvolvimento e status das integrações |
