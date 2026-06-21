# PODP — Projeto Observatório de Dados Públicos

Painel unificado para análise de dados públicos brasileiros. Consolida informações eleitorais (TSE), parlamentares (Câmara, Senado), gastos públicos (Portal da Transparência, TCU), contratações (PNCP) e dados fiscais (SICONFI/IBGE) em uma única plataforma com APIs REST e WebSocket.

---

## Integrações

| Serviço | Cliente | APIs consumidas | Doc |
|---------|---------|-----------------|-----|
| **Câmara dos Deputados** | `deputados` | Deputados, legislaturas, votações, frentes, blocos, grupos | [doc](docs/apis-clientes.md#1-camara-dos-deputados) |
| **Senado Federal** | `senado` | Senadores, comissões, votações, processos, orçamento, agenda | [doc](docs/apis-clientes.md#2-senado-federal) |
| **TCU** | `tcu` | Contas irregulares, inabilitados, inidôneos, fins eleitorais | [doc](docs/apis-clientes.md#3-tcu) |
| **Portal da Transparência** | `portaltransparencia` | Órgãos, pessoas, servidores, cartões, despesas, emendas | [doc](docs/apis-clientes.md#4-portal-da-transparencia) |
| **PNCP** | `pncp` | Contratos, publicações, SSE streaming | [doc](docs/apis-clientes.md#5-pncp) |
| **OpenCNPJ** | `opencnpj` | Dados cadastrais de CNPJ | [doc](docs/apis-clientes.md#6-opencnpj) |
| **IBGE** | `ibge` | Estados, municípios, população | [doc](docs/apis-clientes.md#7-ibge) |
| **SICONFI** | `siconfi` | DCA, RGF, RREO, MSC, entes | [doc](docs/apis-clientes.md#8-siconfi) |

---

## TSE — Dados Eleitorais

### Origem dos dados

Planilhas CSV disponibilizadas pelo TSE, reprocessadas para banco relacional PostgreSQL. Cada ano eleitoral possui diretório em `dataCSV/` com arquivos por UF.

### Tipos de planilhas persistidos

| Planilha (diretório) | Conteúdo | Anos |
|---|---|---|
| `consulta_cand_<ano>/` | Candidatos, eleições, partidos, cargos | 2006–2024 |
| `bem_candidato_<ano>/` | Bens declarados por candidatos | 2006–2024 |
| `prestacao_de_contas_eleitorais_candidatos_<ano>/` | Despesas, receitas e doadores de candidatos | 2018–2024 |
| `prestacao_de_contas_eleitorais_orgaos_partidarios_<ano>/` | Despesas, receitas e doadores de órgãos partidários | 2018–2024 |

### Anos mapeados

`2006` · `2008` · `2010` · `2012` · `2014` · `2016` · `2018` · `2020` · `2022` · `2024`

> A partir de 2018 estão disponíveis também os dados de prestação de contas (receitas/despesas).

### APIs de consulta aos dados do TSE

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/busca/cargos` | Lista todos os cargos disponíveis |
| `GET` | `/busca/partidos` | Lista todos os partidos |
| `POST` | `/busca/candidatos` | Busca candidatos com filtros (cargo, partido, UF, eleito) |
| `POST` | `/busca/doadores` | Busca doador por CPF/CNPJ |
| `POST` | `/busca/fornecedores` | Busca fornecedor por CPF/CNPJ |
| `POST` | `/busca/relacoes` | Busca relações (despesas/receitas) de um CNPJ |
| `POST` | `/entidade` | Consulta entidade por tipo + chave (candidato, fornecedor, doador) |
| `POST` | `/import` | Importa arquivo CSV |

Para detalhes da arquitetura do banco TSE, veja [docs/db-tse.md](docs/db-tse.md).

---

## Análises — Ligação Política

Serviço que cruza documentos de licitações/contratos com dados eleitorais e sanções do TCU.

### API

| Método | Rota | Descrição |
|--------|------|-----------|
| `POST` | `/busca/contexto` | Analisa ligação política de documentos |

### Request

```json
{
  "licitacoes": [
    {
      "numero_controle_pncp": "pncp-001",
      "cpf_cnpj": "11222333000181",
      "socios": [
        { "nome": "João", "documento": "11122233344" }
      ]
    }
  ]
}
```

### Response

```json
{
  "documentos_processados": 1,
  "resultados": [
    {
      "numero_controle_pncp": "pncp-001",
      "cpf_cnpj": "11222333000181",
      "documentos": [
        {
          "documento_input": "11222333000181",
          "documento_normalizado": "11222333000181",
          "nome": "Fornecedor Teste Ltda",
          "parcial": false,
          "origem": "licitacao",
          "vinculos": [
            {
              "tipo": "fornecedor",
              "descricao": "Fornecedor encontrado na base do TSE",
              "detalhes": {
                "fornecedor": { ... },
                "contas_irregulares": [],
                "inidoneos": [],
                "inabilitados": []
              }
            }
          ]
        }
      ]
    }
  ]
}
```

### Fluxo

1. Normaliza documentos (remove prefixo `000` de CPF)
2. Busca correspondência nas bases do TSE (fornecedor/doador)
3. Enriquece com dados do **OpenCNPJ** (razão social, situação cadastral, sócios)
4. Enriquece com sanções do **TCU** (contas irregulares, inidôneos, inabilitados)
5. Utiliza **Redis** como cache para evitar consultas repetidas

---

## Testes

### Tipos implementados

**Teste de integração com Testcontainers + Mocks:**
- Arquivo: `internal/ligacao-politica/handler/handler_integration_test.go`
- Sobe container PostgreSQL 15 real via `testcontainers-go`
- Usa `gomock` para mockar OpenCNPJ, TCU e Redis
- Casos: fornecedor encontrado, doador encontrado, CPF com prefixo 000, documento sem dados, múltiplas licitações, enriquecimento OpenCNPJ e TCU

**Test infrastructure (`internal/shared/testkit/`):**
- `StartPostgresContainer()` — sobe PostgreSQL em container
- `RunMigrations()` — aplica migrations SQL
- `MockDB` — banco mockado para testes unitários
- `NewGinEngine()`, `NewRequest()`, `ExecRequest()` — helpers HTTP
- `InsertFornecedor()`, `InsertDoador()`, `CleanTables()` — fixtures
- `RunTestCase[I]()` — runner genérico com `MockConfig`

### Comandos

```sh
make test              # todos os testes (300s)
make test-unit         # apenas unitários (60s, -short)
make test-integration  # apenas integração (300s)
make test-cover        # com cobertura
make test-race         # com race detector
```

---

## Lint

```sh
make lint       # Executa análise estática (golangci-lint run ./...)
make lint:fix   # Executa lint com correções automáticas (golangci-lint run --fix ./...)
```

O linter é configurado via `.golangci.yml` na raiz. O `make fix` combina `gofmt`, `go mod tidy` e `lint --fix`.

---

## Arquitetura de Pastas

```
backend-analise/
├── main.go                          # Entrypoint
├── internal/
│   ├── app/                         # DI e registro de rotas
│   │   ├── app.go                   # Container de dependências
│   │   └── routes.go                # Todas as rotas HTTP
│   ├── esferas-brasileiras/         # Módulos por esfera federativa
│   │   ├── tse/                     # Dados do TSE
│   │   │   ├── handler/             # Handlers HTTP
│   │   │   ├── usecase/             # Casos de uso
│   │   │   ├── repositorio/         # Acesso a banco
│   │   │   └── importacao/          # Importação de CSVs
│   │   ├── federal/                 # Dados federais
│   │   │   ├── deputados/           # Câmara
│   │   │   ├── senadores/           # Senado
│   │   │   ├── portaltransparencia/ # Portal da Transparência
│   │   │   ├── pncp/                # Contratações públicas
│   │   │   └── tcu/                 # TCU
│   │   ├── estadual/                # Dados estaduais
│   │   └── municipal/               # Dados municipais
│   ├── ligacao-politica/            # Análise de ligação política
│   │   ├── handler/
│   │   ├── usecase/
│   │   └── testutils/
│   └── shared/                      # Código compartilhado
│       ├── clients/                 # Clientes HTTP externos
│       ├── database/                # Pool PostgreSQL
│       ├── migrations/              # Migrations SQL
│       ├── redis/                   # Cache Redis
│       ├── testkit/                 # Helpers de teste
│       ├── types/                   # Tipos compartilhados
│       └── websocket/               # WebSocket utilities
├── dataCSV/                         # CSVs do TSE por ano
├── api/                             # OpenAPI spec
├── docs/                            # Documentação
│   ├── db-tse.md                    # Arquitetura do banco TSE
│   └── apis-clientes.md             # Detalhamento das APIs integradas
├── docker-compose.yml               # PostgreSQL + Redis + Swagger
├── Dockerfile                       # Build multi-stage
└── Makefile                         # Comandos de build/test/lint
```

---

## Documentação Adicional

- [Arquitetura do banco TSE](docs/db-tse.md) — entidades, relacionamentos, FKs, migrações
- [APIs integradas por cliente](docs/apis-clientes.md) — detalhamento de inputs/outputs de cada serviço externo
