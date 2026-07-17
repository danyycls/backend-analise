# Importação de Dados — TSE e Convênios

## Visão Geral

O ODT importa dados de fontes públicas que são disponibilizadas como **planilhas CSV**:

- **TSE (Tribunal Superior Eleitoral)**: dados eleitorais históricos (2006–2024), incluindo candidatos, partidos, eleições, prestação de contas, receitas e despesas de campanha.
- **Portal da Transparência**: convênios e acordos firmados pela Administração Pública Federal.

A importação é executada sob demanda via rota `POST /import` (Server-Sent Events), que percorre o diretório de CSVs, processa os arquivos em paralelo (workers) e persiste em lotes via `pgCOPY` + `MERGE`.

---

## Modelo do Banco

A estrutura completa das tabelas, relacionamentos (FKs), índices e migrations está documentada em **[docs/db-tse.md](./db-tse.md)**.

Resumo das entidades persistidas pela importação:

| Grupo | Entidades |
|-------|-----------|
| Dimensões TSE | `eleicao`, `unidade_eleitoral`, `partido`, `candidato`, `fornecedor`, `doador`, `prestacao_contas` |
| Transações TSE | `despesa_candidato`, `despesa_orgao_partidario`, `receita_candidato`, `receita_orgao_partidario`, `receita_doador_originario_candidato`, `receita_doador_originario_orgao_partidario`, `bem_candidato` |
| Convênios | `convenio` |
| Rastreabilidade | `arquivo_importado` (PK por `caminho_relativo`, com `hash_sha256`) |

O mapeamento CSV → tabelas também está em `docs/db-tse.md` (seção "CSV → Tabelas").

---

## Como preencher o banco (rota `/import`)

### Pré-requisitos

1. **PostgreSQL** configurado e acessível (banco default: `tse_data`).
2. **CSVs do TSE** baixados e organizados no diretório de importação (default: `dataCSV/`).
3. **Migrations** aplicadas (ver `internal/shared/migrations/schema/`).

### Estrutura esperada do diretório

```
dataCSV/
├── 2024/
│   ├── consulta_cand_2024/
│   │   ├── consulta_cand_2024_SP.csv
│   │   └── ...
│   ├── bem_candidato_2024/
│   ├── prestacao_de_contas_eleitorais_candidatos_2024/
│   └── prestacao_de_contas_eleitorais_orgaos_partidarios_2024/
├── 2022/
│   └── ...
├── portalTransparencia/
│   └── 20240101_Convenios.csv
└── ...
```

Os CSVs são particionados por UF (ex: `consulta_cand_2024_SP.csv`) e, quando o volume exige, sub-particionados (`_partaa`, `_partab`, etc.). O caminhamento é recursivo a partir da raiz (`IMPORTACAO_DIRETORIO_CSV`).

### Variáveis de ambiente

| Variável | Default | Descrição |
|----------|---------|-----------|
| `IMPORTACAO_DIRETORIO_CSV` | `dataCSV` | Diretório raiz contendo os CSVs do TSE |
| `IMPORT_BATCH_SIZE` | `10000` | Número de registros por lote no `pgCOPY` |
| `IMPORT_MAX_WORKERS` | `NumCPU * 2` | Número máximo de goroutines (workers) lendo arquivos em paralelo |
| `IMPORT_FILES_PER_BATCH` | `50` | Quantos arquivos são lidos por lote antes de persistir |
| `DB_HOST` | `localhost` | Host do PostgreSQL |
| `DB_PORT` | `5432` | Porta do PostgreSQL |
| `DB_USER` | `postgres` | Usuário do PostgreSQL |
| `DB_PASSWORD` | `postgres` | Senha do PostgreSQL |
| `DB_NAME` | `tse_data` | Nome do banco de dados |

### Como chamar

```
POST /import
Content-Type: text/event-stream
```

A rota retorna um **stream SSE** com eventos de progresso em tempo real. A requisição não possui body.

### Eventos SSE

| Evento | Gatilho | Payload |
|--------|---------|---------|
| `progression` | A cada 1 segundo | `{ total_arquivos, total_diretorios, diretorio_indice, arquivos_lendo, arquivos_lidos, arquivos_persistindo, arquivos_persistidos, arquivos_ignorados, arquivos_restantes }` |
| `concluido` | Término com sucesso | `{ sucesso: 1, total_registros, arquivos_persistidos }` |
| `erro` | Falha na importação | `{ sucesso: 0 }` |

### Idempotência

Antes de processar, o UseCase consulta a tabela `arquivo_importado` e **ignora arquivos já importados** (por `caminho_relativo`). Re-executar `POST /import` após uma falha processa apenas os arquivos faltantes.

---

## Arquitetura do Leitor CSV

### Diagrama de fluxo

```
POST /import (SSE)
      │
      ▼
┌──────────────┐
│   Handler     │  SSE: progression (1s) / concluido / erro
└──────┬───────┘
       ▼
┌──────────────┐
│   UseCase     │  1. Lista arquivos (Service)
│              │  2. Agrupa por diretório
│              │  3. Ordena por prioridade (consulta_cand → bem → ...)
│              │  4. Worker pool lê lotes de arquivos
└──┬───┬───┬──┘
   ▼   ▼   ▼  (goroutines — até IMPORT_MAX_WORKERS)
┌────┐┌────┐┌────┐
│ W1 ││ W2 ││ W3 │  Cada worker:
└─┬──┘└─┬──┘└─┬──┘  1. Abre CSV (ISO-8859-1 → UTF-8, delimitador ';')
  │     │     │     2. Parser por tipo → preenche DadosImportacao
  └─────┴─────┘     3. Merge no coletor do lote
       ▼
┌──────────────┐
│  pgCOPY +     │  14 passos em ordem de dependência:
│  MERGE        │  dimensões (com RETURNING + remapeamento de IDs) →
│              │  transações (insert em lote)
└──────┬───────┘
       ▼
   PostgreSQL (tse_data)
   + arquivo_importado (rastreabilidade)
```

### Camadas

| Camada | Arquivo | Responsabilidade |
|--------|---------|------------------|
| Handler | `importacao/handler/leitorCSV-handler.go` | Recebe `POST /import`, emite SSE de progresso |
| UseCase | `importacao/usecase/importar-csv-usecase.go` | Orquestra: agrupa por diretório, ordena, worker pool, persiste lotes, controla progresso |
| Service | `importacao/service/leitorCSV-service.go` | Caminha o diretório (`filepath.WalkDir`), identifica tipos, delega leitura ao parser |
| Parse | `importacao/parse/*.go` | Leitura CSV, parsers por entidade, merge, helpers de persistência |
| Repositórios | `importacao/repositorios/*.go` | `pgCOPY` com staging + `MERGE`, pool de conexões (`pgxpool`) |

### Ordem de processamento por diretório

Os diretórios são ordenados por prioridade para garantir que as dependências (ex: candidatos) sejam carregadas antes das transações que as referenciam (ex: bens, despesas, receitas):

| Prioridade | Diretório contém | Motivo |
|------------|-------------------|--------|
| 0 | `convenios` / `portalTransparencia` | Dados independentes, processados primeiro |
| 1 | `consulta_cand` | Cria candidatos/partidos/eleições/unidades |
| 2 | `bem_candidato` | Depende de candidatos |
| 3 | `candidatos` (prestação de contas) | Depende de candidatos |
| 4 | `orgaos_partidarios` / `orgao_partidario` | Depende de partidos |
| 99 | demais | Ordem alfabética |

Para os diretórios das prioridades 2–4, o UseCase carrega um **cache de candidatos** do banco em memória (`map[int64]*types.Candidato` indexado por `sq_candidato`) para acelerar a resolução de FKs sem consultar o banco a cada linha.

### Workers

- Pool de goroutines limitado por `IMPORT_MAX_WORKERS` (default: `NumCPU * 2`).
- A cada `IMPORT_FILES_PER_BATCH` arquivos (default: 50), os workers finalizam, os dados lidos são **merged** em um coletor acumulado e o lote é **persistido**.
- Após cada lote: `runtime.GC()` + `debug.FreeOSMemory()` para liberar memória.
- Erro em qualquer worker aborta o lote e propaga o erro (primeiro erro vence via `sync.Once`).

---

## Formato dos CSVs

Os arquivos do TSE seguem o padrão:

- **Encoding:** ISO-8859-1 (Latin-1), convertido para UTF-8 na leitura (`charmap.ISO8859_1.NewDecoder()`).
- **Delimitador:** ponto-e-vírgula (`;`).
- **Campos variáveis:** `FieldsPerRecord = -1` (linhas com mais/menos colunas que o header são toleradas).
- **Aspas:** `LazyQuotes = true` (tolera aspas mal-formadas, comuns em CSVs do TSE).
- **Buffer:** `bufio.NewReaderSize` de 64KB.
- **Mapas de linha:** reutilizados via `sync.Pool` para reduzir alocações.

Valores nulos do TSE (`#NULO`, `#NE`, `#NULO#`, `-1`) são normalizados para strings vazias ou zeros pelo parser (`parse-tipos.go`).

### Tipos de arquivo suportados

| Prefixo do arquivo | Tipo (interno) | Prioridade | Tabela(s) destino |
|--------------------|----------------|------------|-------------------|
| `*_convenios.csv` | `convenio_portal_transparencia` | 0 | `convenio` |
| `consulta_cand_` | `consulta_candidato` | 1 | `candidato`, `eleicao`, `unidade_eleitoral`, `partido` |
| `bem_candidato_` | `bem_candidato` | 2 | `bem_candidato` |
| `despesas_contratadas_candidatos_` | `despesa_contratada_candidato` | 3 | `despesa_candidato` (CONTRATADA), `fornecedor`, `prestacao_contas` |
| `receitas_candidatos_` | `receita_candidato` | 4 | `receita_candidato`, `doador`, `prestacao_contas` |
| `receitas_candidatos_doador_originario_` | `receita_candidato_doador_originario` | 5 | `receita_doador_originario_candidato` |
| `despesas_contratadas_orgaos_partidarios_` | `despesa_contratada_orgao_partidario` | 6 | `despesa_orgao_partidario`, `fornecedor`, `prestacao_contas` |
| `receitas_orgaos_partidarios_` | `receita_orgao_partidario` | 7 | `receita_orgao_partidario`, `doador`, `prestacao_contas` |
| `receitas_orgaos_partidarios_doador_originario_` | `receita_orgao_partidario_doador_originario` | 8 | `receita_doador_originario_orgao_partidario` |
| `despesas_pagas_candidatos_` | `despesa_paga_candidato` | 9 | `despesa_candidato` (PAGA) |
| `despesas_pagas_orgaos_partidarios_` | `despesa_paga_orgao_partidario` | 10 | `despesa_orgao_partidario` (PAGA) |

> **Nota:** Arquivos de convênios (`_convenios.csv`) são processados com prioridade máxima (0), antes de qualquer dado TSE.

Arquivos que não correspondam a nenhum prefixo são ignorados pela descoberta.

---

## Persistência (pgCOPY + MERGE)

### Padrão de escrita

Para cada entidade, o repositório aplica o padrão **staging → COPY → MERGE**:

1. Cria (ou recria) uma **temporary staging table** (`stg_<tabela>`) com `LIKE <tabela> INCLUDING DEFAULTS`.
2. `TRUNCATE` a staging.
3. `pgCOPY` (`tx.CopyFrom`) insere os registros em lote na staging.
4. `INSERT INTO <tabela> SELECT ... FROM stg_<tabela> ON CONFLICT <conflict_target> DO UPDATE SET ... RETURNING <cols>` faz o upsert e retorna os IDs resolvidos.
5. `TRUNCATE` a staging.

Para entidades de dimensão (com FKs que precisam ser resolvidas), o `RETURNING` retorna o ID persistido + a chave de negócio, permitindo **remapear os IDs temporários** nos relacionamentos dependêntes (ex: após inserir eleições, o `eleicao_id` de cada candidato é atualizado).

### Advisory locks

Cada dimensão é protegida por um **advisory lock transacional** (`pg_advisory_xact_lock(hashtext($1))`) para evitar concorrência entre lotes concorrentes sobre a mesma tabela.

### 14 passos em ordem de dependência

A persistência executa 14 passos sequenciais por lote, respeitando a ordem de dependência entre entidades (dimensões antes de transações):

| Passo | Entidade | Tipo | Dependências |
|-------|----------|------|--------------|
| 1 | `eleicao` | Dimensão | — |
| 2 | `unidade_eleitoral` | Dimensão | — |
| 3 | `partido` | Dimensão | — |
| 4 | `candidato` | Dimensão | eleicao, partido |
| 5 | `fornecedor` | Dimensão | — |
| 6 | `doador` | Dimensão | — |
| 7 | `prestacao_contas` | Dimensão | eleicao, candidato, partido, unidade_eleitoral |
| 8 | `despesa_candidato` | Transação | prestacao_contas, candidato, fornecedor |
| 9 | `despesa_orgao_partidario` | Transação | prestacao_contas, partido, fornecedor |
| 10 | `receita_candidato` | Transação | prestacao_contas, candidato, doador |
| 11 | `receita_orgao_partidario` | Transação | prestacao_contas, partido, doador |
| 12 | `receita_doador_originario_candidato` | Transação | receita_candidato |
| 13 | `receita_doador_originario_orgao_partidario` | Transação | receita_orgao_partidario |
| 14 | `bem_candidato` | Transação | candidato |

### Transação e rastreabilidade

- Cada lote de arquivos é persistido em uma **única transação** (`BEGIN` → 14 passos → `COMMIT`).
- Em caso de erro, a transação é desfeita (`ROLLBACK`) e o erro propaga, abortando o diretório corrente.
- Ao final do lote (após o commit), cada arquivo é registrado em `arquivo_importado` com `caminho_relativo`, `nome`, `tipo`, `uf` e `total_registros` — garantindo idempotência.

---

## Referências

- **[docs/db-tse.md](./db-tse.md)** — Modelo do banco (entidades, FKs, índices, migrations, mapeamento CSV → tabelas)
- **Código-fonte:** `internal/sources/tse/importacao/`
  - `handler/` — endpoint SSE
  - `usecase/` — orquestração e worker pool
  - `service/` — descoberta de arquivos
  - `parse/` — leitura CSV e parsers por entidade
  - `repositorios/` — pgCOPY, pool e repositório
  - `types/` — estruturas de dados e constantes
- **Migrations:** `internal/shared/migrations/schema/`
