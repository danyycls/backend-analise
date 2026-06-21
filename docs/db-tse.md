# Banco TSE — Arquitetura de Dados

## Visão Geral

Os dados do TSE (Tribunal Superior Eleitoral) são originalmente disponibilizados como planilhas CSV por ano eleitoral e por UF. O projeto importa esses CSVs para um banco PostgreSQL relacional, normalizando as entidades e mantendo a rastreabilidade por `arquivo_importado`.

Todos os dados históricos de **2006 a 2024** estão mapeados, sendo que a partir de **2018** também estão disponíveis os dados de prestação de contas (receitas e despesas de campanha).

- **Soft delete**: todas as tabelas possuem `deleted_at TIMESTAMPTZ`
- **Auditoria**: `created_at` e `updated_at` em todas as tabelas
- **IDs**: UUID gerados via `gen_random_uuid()` (pgcrypto)

---

## Entidades e Relacionamentos

### 1. `eleicao`
Processo eleitoral (ex: "Eleições Gerais 2022", "Eleições Municipais 2024").

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `codigo_tse` | INTEGER | UNIQUE | Código do TSE |
| `ano` | SMALLINT | | Ano da eleição |
| `codigo_tipo_eleicao` | INTEGER | | |
| `nome_tipo_eleicao` | VARCHAR(100) | | Ex: "Eleição Ordinária" |
| `descricao` | VARCHAR(255) | | Ex: "ELEIÇÃO FEDERAL 2022" |
| `data_eleicao` | DATE | | |

### 2. `unidade_eleitoral`
Unidades da federação / municípios no contexto eleitoral.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `sg_uf` | VARCHAR(2) | | UF |
| `codigo_tse` | VARCHAR(16) | | Código TSE da unidade |
| `nome` | VARCHAR(255) | | Nome (ex: "São Paulo") |
| UNIQUE | | `(sg_uf, codigo_tse)` | |

### 3. `partido`
Partidos políticos e coligações/federações.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `numero` | SMALLINT | UNIQUE | Número do partido |
| `sigla` | VARCHAR(20) | | Sigla (ex: "PT") |
| `nome` | VARCHAR(255) | | Nome completo |
| `federacao_codigo_tse` | BIGINT | | Código da federação partidária |
| `federacao_sigla` | VARCHAR(50) | | Sigla da federação |
| `federacao_nome` | VARCHAR(255) | | Nome da federação |
| `coligacao_codigo_tse` | BIGINT | | Código da coligação |
| `coligacao_nome` | VARCHAR(255) | | Nome da coligação |
| `coligacao_composicao` | TEXT | | Composición (ex: "PT/PCdoB/PV") |

### 4. `candidato`
Candidatos a cargos eletivos.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `sq_candidato` | BIGINT | UNIQUE | Sequencial TSE do candidato |
| `eleicao_id` | UUID | FK → `eleicao(id)` | Eleição |
| `sg_uf` | VARCHAR(2) | | UF da candidatura |
| `partido_id` | UUID | FK → `partido(id)` | Partido |
| `cargo_codigo` | INTEGER | | |
| `cargo_nome` | VARCHAR(100) | | Ex: "Deputado Federal" |
| `genero_descricao` | VARCHAR(100) | | |
| `cor_raca_descricao` | VARCHAR(100) | | |
| `estado_civil_nome` | VARCHAR(100) | | |
| `grau_instrucao_nome` | VARCHAR(150) | | |
| `ocupacao_codigo` | INTEGER | | |
| `ocupacao_nome` | VARCHAR(255) | | |
| `numero_candidato` | INTEGER | | Número de urna |
| `cpf` | VARCHAR(11) | | |
| `cpf_vice` | VARCHAR(11) | | CPF do vice |
| `nome_completo` | VARCHAR(255) | NOT NULL | |
| `nome_urna` | VARCHAR(255) | | |
| `nome_social` | VARCHAR(255) | | |
| `data_nascimento` | DATE | | |
| `situacao_totalizacao_descricao` | VARCHAR(255) | | Ex: "ELEITO" |

### 5. `bem_candidato`
Bens declarados por candidatos.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `candidato_id` | UUID | FK → `candidato(id)` | |
| `tipo_bem_codigo` | INTEGER | | |
| `tipo_bem_nome` | VARCHAR(255) | | |
| `numero_ordem` | INTEGER | | |
| `descricao` | TEXT | NOT NULL | |
| `valor` | NUMERIC(18,2) | NOT NULL | |
| `data_ultima_atualizacao` | DATE | | |
| `hora_ultima_atualizacao` | TIME | | |
| UNIQUE | | `(candidato_id, numero_ordem)` | |

### 6. `fornecedor`
Fornecedores de campanha (despesas).

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `cpf_cnpj` | VARCHAR(14) | UNIQUE | Documento sem formatação |
| `nome` | VARCHAR(255) | NOT NULL | |
| `nome_rfb` | VARCHAR(255) | | Nome na Receita Federal |
| `tipo_fornecedor_codigo` | INTEGER | | |
| `tipo_fornecedor_descricao` | VARCHAR(100) | | |
| `cnae_codigo` | VARCHAR(20) | | |
| `cnae_descricao` | VARCHAR(255) | | |
| `esfera_partidaria_codigo` | VARCHAR(10) | | |
| `esfera_partidaria_descricao` | VARCHAR(100) | | |
| `sg_uf` | VARCHAR(2) | | |
| `municipio_nome` | VARCHAR(255) | | |
| `sq_candidato_relacionado` | BIGINT | | |
| `numero_candidato_relacionado` | INTEGER | | |
| `cargo_codigo_relacionado` | INTEGER | | |
| `cargo_descricao_relacionada` | VARCHAR(100) | | |
| `partido_numero_relacionado` | SMALLINT | | |
| `partido_sigla_relacionado` | VARCHAR(20) | | |
| `partido_nome_relacionado` | VARCHAR(255) | | |

### 7. `doador`
Doadores de campanha (receitas).

Mesma estrutura de `fornecedor` (cpf_cnpj, nome, nome_rfb, cnae, etc.).

### 8. `prestacao_contas`
Prestações de contas — ponto central que conecta eleições a candidatos ou órgãos partidários.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `sq_prestador_contas` | BIGINT | | Sequencial TSE |
| `eleicao_id` | UUID | FK → `eleicao(id)` | |
| `candidato_id` | UUID | FK → `candidato(id)` | NULL se órgão partidário |
| `partido_id` | UUID | FK → `partido(id)` | NULL se candidato |
| `sg_uf` | VARCHAR(2) | | |
| `unidade_eleitoral_id` | UUID | FK → `unidade_eleitoral(id)` | |
| `tipo_prestador` | VARCHAR(30) | CHECK `IN ('CANDIDATO', 'ORGAO_PARTIDARIO')` | |
| `tipo_prestacao` | VARCHAR(30) | | |
| `data_prestacao` | DATE | | |
| `turno` | SMALLINT | | 1 ou 2 |
| `cnpj_prestador_conta` | VARCHAR(14) | | |
| `esfera_partidaria_codigo` | VARCHAR(10) | | |
| `esfera_partidaria_descricao` | VARCHAR(100) | | |
| UNIQUE | | `(tipo_prestador, eleicao_id, sq_prestador_contas)` | |

**Regra de negócio**: se `tipo_prestador = 'CANDIDATO'`, então `candidato_id` é obrigatório e `partido_id` é NULL; se `tipo_prestador = 'ORGAO_PARTIDARIO'`, o inverso.

### 9. `despesa_candidato`
Despesas de campanha de candidatos.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `prestacao_contas_id` | UUID | FK → `prestacao_contas(id)` | |
| `candidato_id` | UUID | FK → `candidato(id)` | |
| `fornecedor_id` | UUID | FK → `fornecedor(id)` | |
| `sq_despesa` | BIGINT | | |
| `tipo_registro` | VARCHAR(20) | CHECK `IN ('CONTRATADA', 'PAGA')` | |
| `tipo_documento` | VARCHAR(100) | | |
| `numero_documento` | VARCHAR(100) | | |
| `origem_despesa_codigo` | INTEGER | | |
| `origem_despesa_descricao` | VARCHAR(255) | | |
| `fonte_despesa_codigo` | INTEGER | | |
| `fonte_despesa_descricao` | VARCHAR(255) | | |
| `natureza_despesa_codigo` | INTEGER | | |
| `natureza_despesa_descricao` | VARCHAR(255) | | |
| `especie_recurso_codigo` | INTEGER | | |
| `especie_recurso_descricao` | VARCHAR(255) | | |
| `sq_parcelamento_despesa` | BIGINT | | |
| `data_despesa` | DATE | | |
| `descricao` | TEXT | NOT NULL | |
| `valor` | NUMERIC(18,2) | NOT NULL | |
| UNIQUE | | `(sq_despesa, tipo_registro)` | |

### 10. `despesa_orgao_partidario`
Despesas de campanha de órgãos partidários.

Mesma estrutura de `despesa_candidato`, mas com `partido_id` no lugar de `candidato_id`.

### 11. `receita_candidato`
Receitas de campanha de candidatos.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `prestacao_contas_id` | UUID | FK → `prestacao_contas(id)` | |
| `candidato_id` | UUID | FK → `candidato(id)` | |
| `doador_id` | UUID | FK → `doador(id)` | |
| `sq_receita` | BIGINT | UNIQUE | |
| `fonte_receita_codigo` | INTEGER | | |
| `fonte_receita_descricao` | VARCHAR(255) | | |
| `origem_receita_codigo` | INTEGER | | |
| `origem_receita_descricao` | VARCHAR(255) | | |
| `natureza_receita_codigo` | INTEGER | | |
| `natureza_receita_descricao` | VARCHAR(255) | | |
| `especie_receita_codigo` | INTEGER | | |
| `especie_receita_descricao` | VARCHAR(255) | | |
| `numero_recibo_doacao` | VARCHAR(100) | | |
| `numero_documento_doacao` | VARCHAR(100) | | |
| `data_receita` | DATE | | |
| `descricao` | TEXT | NOT NULL | |
| `valor` | NUMERIC(18,2) | NOT NULL | |
| `natureza_recurso_estimavel` | TEXT | | |
| `genero` | VARCHAR(100) | | |
| `cor_raca` | VARCHAR(100) | | |

### 12. `receita_orgao_partidario`
Receitas de órgãos partidários.

Mesma estrutura de `receita_candidato`, mas com `partido_id` no lugar de `candidato_id`.

### 13. `receita_doador_originario_candidato` / `receita_doador_originario_orgao_partidario`
Doadores originários (quando a doação é recebida via outro doador — ex: pessoa física doando via partido).

| Coluna | Tipo | Restrição |
|--------|------|-----------|
| `id` | UUID | PK |
| `prestacao_contas_id` | UUID | FK → `prestacao_contas(id)` |
| `receita_candidato_id` / `receita_orgao_partidario_id` | UUID | FK |
| `sq_receita` | BIGINT | |
| `documento_doador` | VARCHAR(14) | |
| `nome_doador` | VARCHAR(255) | NOT NULL |
| `nome_doador_rfb` | VARCHAR(255) | |
| `tipo_doador` | VARCHAR(100) | |
| `cnae_codigo` | VARCHAR(20) | |
| `cnae_descricao` | VARCHAR(255) | |
| `data_receita` | DATE | |
| `descricao` | TEXT | NOT NULL |
| `valor` | NUMERIC(18,2) | NOT NULL |

### 14. `arquivo_importado`
Rastreabilidade de importação de CSVs.

| Coluna | Tipo | Descrição |
|--------|------|-----------|
| `caminho_relativo` | VARCHAR(500) | PK |
| `nome` | VARCHAR(255) | Nome do arquivo |
| `tipo` | VARCHAR(100) | Tipo (ex: "consulta_cand") |
| `uf` | VARCHAR(2) | UF |
| `total_registros` | INTEGER | |
| `criado_em` | TIMESTAMPTZ | |

---

## Diagrama de Relacionamentos (FKs)

```
eleicao ──┬── candidato
          ├── prestacao_contas
          └── (referenciada indiretamente)

candidato ──┬── bem_candidato
            ├── prestacao_contas (quando CANDIDATO)
            ├── despesa_candidato
            └── receita_candidato

partido ──┬── candidato
           ├── prestacao_contas (quando ORGAO_PARTIDARIO)
           ├── despesa_orgao_partidario
           └── receita_orgao_partidario

fornecedor ──┬── despesa_candidato
             └── despesa_orgao_partidario

doador ──┬── receita_candidato
         └── receita_orgao_partidario

prestacao_contas ──┬── despesa_candidato
                    ├── despesa_orgao_partidario
                    ├── receita_candidato
                    ├── receita_orgao_partidario
                    ├── receita_doador_originario_candidato
                    └── receita_doador_originario_orgao_partidario
```

---

## CSV → Tabelas (Mapeamento)

| Arquivo CSV (pasta) | Tabela(s) destino | Anos disponíveis |
|---|---|---|
| `consulta_cand_<ano>/` | `candidato`, `eleicao`, `unidade_eleitoral`, `partido` | 2006–2024 |
| `bem_candidato_<ano>/` | `bem_candidato` | 2006–2024 |
| `prestacao_de_contas_eleitorais_candidatos_<ano>/` | `prestacao_contas` (CANDIDATO), `despesa_candidato`, `receita_candidato`, `receita_doador_originario_candidato`, `fornecedor`, `doador` | 2018–2024 |
| `prestacao_de_contas_eleitorais_orgaos_partidarios_<ano>/` | `prestacao_contas` (ORGAO_PARTIDARIO), `despesa_orgao_partidario`, `receita_orgao_partidario`, `receita_doador_originario_orgao_partidario`, `fornecedor`, `doador` | 2018–2024 |

Todos os CSVs são particionados por UF (ex: `consulta_cand_2024_SP.csv`), exceto quando o volume exige sub-partição (`_partaa`, `_partab`, etc.).

---

## Índices Principais

- `candidato`: `sq_candidato`, `cpf`, `(eleicao_id, sg_uf)`, `partido_id`
- `fornecedor`: `cpf_cnpj`, `sg_uf`
- `doador`: `cpf_cnpj`, `sg_uf`
- `despesa_candidato`: `sq_despesa`, `candidato_id`, `fornecedor_id`, `prestacao_contas_id`, `data_despesa`
- `despesa_orgao_partidario`: `sq_despesa`, `partido_id`, `fornecedor_id`, `prestacao_contas_id`, `data_despesa`
- `receita_candidato`: `sq_receita`, `candidato_id`, `doador_id`, `prestacao_contas_id`, `data_receita`
- `receita_orgao_partidario`: `sq_receita`, `partido_id`, `doador_id`, `prestacao_contas_id`, `data_receita`

---

## Migrations

As migrations estão em `internal/shared/migrations/schema/`:

| Arquivo | Descrição |
|---|---|
| `000001_esquema_inicial.up.sql` | Cria todas as tabelas, índices e constraints |
| `000001_esquema_inicial.down.sql` | Remove todas as tabelas |
| `000002_remove_uq_partido_sigla.up.sql` | Remove unique constraint de `partido.sigla` |
| `000002_remove_uq_partido_sigla.down.sql` | Reverte a remoção |
