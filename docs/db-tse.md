# Banco de Dados — Arquitetura de Dados

## Visão Geral

O banco PostgreSQL do ODT armazena dados de múltiplas fontes públicas importadas via CSV e via API:

- **TSE (Tribunal Superior Eleitoral)**: dados eleitorais históricos (2006–2024), incluindo candidatos, partidos, eleições, prestação de contas, receitas e despesas de campanha.
- **Portal da Transparência**: convênios e acordos firmados pela Administração Pública Federal.
- **PNCP (Portal Nacional de Contratações Públicas)**: contratos, licitações, fornecedores e societários.

### Padrões do Banco

- **Soft delete**: tabelas TSE e convênios possuem `deleted_at TIMESTAMPTZ`
- **Auditoria**: `created_at` e `updated_at` em todas as tabelas
- **IDs**: UUID gerados via `gen_random_uuid()` (pgcrypto)
- **Hash de integridade**: `arquivo_importado` possui `hash_sha256` para rastreabilidade de conteúdo

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
| `hash_sha256` | VARCHAR(64) | Hash SHA-256 do conteúdo (idempotência) |
| `criado_em` | TIMESTAMPTZ | |

### 15. `convenio`
Convênios e acordos do Portal da Transparência (importados de CSV).

| Coluna | Tipo | Descrição |
|--------|------|-----------|
| `id` | UUID | PK |
| `numero_convenio` | VARCHAR(50) | Número do convênio |
| `uf` | VARCHAR(2) | UF do município convenente |
| `codigo_siafi_municipio` | VARCHAR(20) | Código SIAFI do município |
| `nome_municipio` | VARCHAR(255) | Nome do município |
| `situacao_convenio` | VARCHAR(100) | Situação (ativo, concluído, etc.) |
| `numero_original` | VARCHAR(100) | Número original do instrumento |
| `numero_processo` | VARCHAR(100) | Número do processo administrativo |
| `objeto_convenio` | TEXT | Descrição do objeto do convênio |
| `codigo_orgao_superior` | VARCHAR(20) | Código do órgão superior |
| `nome_orgao_superior` | VARCHAR(255) | Nome do órgão superior |
| `codigo_orgao_concedente` | VARCHAR(20) | Código do órgão concedente |
| `nome_orgao_concedente` | VARCHAR(255) | Nome do órgão concedente |
| `codigo_ug_concedente` | VARCHAR(20) | Código da unidade gestora |
| `nome_ug_concedente` | VARCHAR(255) | Nome da unidade gestora |
| `codigo_convenente` | VARCHAR(20) | CPF/CNPJ do convenente |
| `tipo_convenente` | VARCHAR(100) | Tipo de pessoa do convenente |
| `nome_convenente` | VARCHAR(255) | Nome do convenente |
| `tipo_ente_convenente` | VARCHAR(100) | Tipo de ente (município, estado, etc.) |
| `tipo_instrumento` | VARCHAR(100) | Tipo (convênio, contrato de repasse, etc.) |
| `valor_convenio` | NUMERIC(18,2) | Valor total do convênio |
| `valor_liberado` | NUMERIC(18,2) | Valor total liberado |
| `data_publicacao` | DATE | Data de publicação |
| `data_inicio_vigencia` | DATE | Início da vigência |
| `data_final_vigencia` | DATE | Fim da vigência |
| `valor_contrapartida` | NUMERIC(18,2) | Valor da contrapartida |
| `data_ultima_liberacao` | DATE | Data da última liberação |
| `valor_ultima_liberacao` | NUMERIC(18,2) | Valor da última liberação |
| `created_at` | TIMESTAMPTZ | Data de criação do registro |
| `updated_at` | TIMESTAMPTZ | Data de atualização |
| `deleted_at` | TIMESTAMPTZ | Soft delete |

### 16. `amparo_legal`
Tabela de apoio com amparos legais para licitações.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `codigo` | INTEGER | PK | |
| `nome` | TEXT | NOT NULL | |
| `descricao` | TEXT | | |

### 17. `licitacao_fornecedor`
Dados cadastrais de fornecedores de licitações.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `cnpj` | VARCHAR(14) | PK | |
| `razao_social` | TEXT | NOT NULL | |
| `dados_completos` | JSONB | | Dados completos do fornecedor |

### 18. `socio`
Sócios de pessoas jurídicas.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `cnpj_cpf_socio` | VARCHAR(14) | UNIQUE | CPF/CNPJ do sócio |
| `nome_socio` | TEXT | | Nome do sócio |

### 19. `fornecedor_socio`
Relação entre fornecedores e seus sócios.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `cnpj_fornecedor` | VARCHAR(14) | PK, FK → `licitacao_fornecedor(cnpj)` | |
| `socio_id` | UUID | PK, FK → `socio(id)` | |
| `data_entrada_sociedade` | TEXT | | |
| `identificador_socio` | TEXT | | |
| `nome_socio` | TEXT | | |
| `qualificacao_socio` | TEXT | | |
| `nome_representante` | TEXT | | |
| `qualificacao_representante` | TEXT | | |
| `representante_legal` | TEXT | | |
| `faixa_etaria` | TEXT | | |
| `pais_codigo` | TEXT | | |
| `pais_descricao` | TEXT | | |

### 20. `licitacao_contrato`
Contratos de licitações do PNCP.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `numero_controle_pncp` | VARCHAR(255) | UNIQUE, NOT NULL | Número de controle PNCP |
| `cnpj_orgao` | VARCHAR(14) | NOT NULL | CNPJ do órgão |
| `ug_uf_sigla` | VARCHAR(2) | | UF da UG |
| `ug_codigo_ibge` | VARCHAR(7) | | Código IBGE da UG |
| `data_publicacao_pncp` | DATE | | Data de publicação |
| `data_assinatura` | DATE | | Data de assinatura |
| `data_inicio_vigencia` | DATE | | Início da vigência |
| `data_termino_vigencia` | DATE | | Término da vigência |
| `valor_global` | NUMERIC(18,2) | | Valor global |
| `valor_inicial` | NUMERIC(18,2) | | Valor inicial |
| `valor_total_estimado` | NUMERIC(18,2) | | Valor total estimado |
| `valor_total_homologado` | NUMERIC(18,2) | | Valor total homologado |
| `ni_fornecedor` | VARCHAR(14) | | CNPJ/CPF do fornecedor |
| `codigo_amparo_legal` | INTEGER | FK → `amparo_legal(codigo)` | |
| `numero_contrato` | VARCHAR(50) | | Número do contrato |
| `codigo_contrato` | VARCHAR(50) | | Código do contrato |
| `codigo_tipo_contrato` | INTEGER | | |
| `tipo_contrato_nome` | TEXT | | |
| `codigo_ug` | VARCHAR(20) | | |
| `nome_ug` | TEXT | | |
| `ug_municipio_nome` | TEXT | | |
| `ug_uf_nome` | TEXT | | |
| `modalidade_nome` | TEXT | | |
| `codigo_orgao` | VARCHAR(20) | | |
| `nome_orgao` | TEXT | | |
| `nome_orgao_sub` | TEXT | | |
| `objeto_contrato` | TEXT | | Objeto do contrato |
| `numero_licitacao` | VARCHAR(50) | | |
| `origem_licitacao` | TEXT | | |
| `produto` | TEXT | | |
| `subtipo_contrato` | TEXT | | |
| `ano_contrato` | INTEGER | | |
| `nome_razao_social_fornecedor` | TEXT | | |
| `dados_completos` | JSONB | NOT NULL DEFAULT '{}' | Dados completos originais |
| `created_at` | TIMESTAMPTZ | NOT NULL | |
| `updated_at` | TIMESTAMPTZ | NOT NULL | |

### 21. `licitacao_busca_controle`
Controle de buscas realizadas no PNCP.

| Coluna | Tipo | Restrição | Descrição |
|--------|------|-----------|-----------|
| `id` | UUID | PK | |
| `tipo_busca` | VARCHAR(20) | NOT NULL | Tipo (municipio, uf, orgao) |
| `valor_busca` | VARCHAR(20) | NOT NULL | Valor buscado |
| `ano` | INTEGER | NOT NULL | Ano de referência |
| `mes` | INTEGER | NOT NULL | Mês de referência |
| `data_inicial` | DATE | NOT NULL | |
| `data_final` | DATE | NOT NULL | |
| `total_contratos_encontrados` | INTEGER | NOT NULL DEFAULT 0 | |
| `ultima_atualizacao` | TIMESTAMPTZ | NOT NULL | |
| UNIQUE | | `(tipo_busca, valor_busca, ano, mes)` | |

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

### TSE

| Arquivo CSV (pasta) | Tabela(s) destino | Anos disponíveis |
|---|---|---|
| `consulta_cand_<ano>/` | `candidato`, `eleicao`, `unidade_eleitoral`, `partido` | 2006–2024 |
| `bem_candidato_<ano>/` | `bem_candidato` | 2006–2024 |
| `prestacao_de_contas_eleitorais_candidatos_<ano>/` | `prestacao_contas` (CANDIDATO), `despesa_candidato`, `receita_candidato`, `receita_doador_originario_candidato`, `fornecedor`, `doador` | 2018–2024 |
| `prestacao_de_contas_eleitorais_orgaos_partidarios_<ano>/` | `prestacao_contas` (ORGAO_PARTIDARIO), `despesa_orgao_partidario`, `receita_orgao_partidario`, `receita_doador_originario_orgao_partidario`, `fornecedor`, `doador` | 2018–2024 |

### Portal da Transparência

| Arquivo CSV (pasta) | Tabela(s) destino |
|---|---|
| `*_convenios.csv` | `convenio` |

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
| `000001_esquema_inicial.up.sql` | Cria tabelas TSE: eleicao, unidade_eleitoral, partido, candidato, bem_candidato, fornecedor, doador, prestacao_contas, despesa_candidato, despesa_orgao_partidario, receita_candidato, receita_orgao_partidario, receita_doador_originario_candidato, receita_doador_originario_orgao_partidario, arquivo_importado |
| `000001_esquema_inicial.down.sql` | Remove todas as tabelas TSE |
| `000002_remove_uq_partido_sigla.up.sql` | Remove unique constraint de `partido.sigla` |
| `000002_remove_uq_partido_sigla.down.sql` | Reverte a remoção |
| `000003_convenios_portal_transparencia.up.sql` | Cria tabela `convenio` com índices (pg_trgm para busca textual) |
| `000003_convenios_portal_transparencia.down.sql` | Remove tabela `convenio` |
| `000004_pncp_licitacoes.up.sql` | Cria tabelas PNCP: `amparo_legal`, `licitacao_fornecedor`, `socio`, `fornecedor_socio`, `licitacao_contrato`, `licitacao_busca_controle` |
| `000004_pncp_licitacoes.down.sql` | Remove tabelas PNCP |
| `000005_add_hash_arquivo.up.sql` | Adiciona coluna `hash_sha256` em `arquivo_importado` + índice |
| `000005_add_hash_arquivo.down.sql` | Remove coluna `hash_sha256` |
| `000006_remove_uq_transacionais.up.sql` | Remove unique constraints de tabelas transacionais (prestacao_contas, despesa_*, receita_*, convenio) |
| `000006_remove_uq_transacionais.down.sql` | Reverte a remoção |
