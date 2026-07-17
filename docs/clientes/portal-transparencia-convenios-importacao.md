# Importação de Convênios do Portal da Transparência

## Origem dos Dados

Os dados são obtidos através do download de arquivos CSV do Portal da Transparência
do Governo Federal: https://portaldatransparencia.gov.br/download-de-dados

A seção "Convenios e Outros Acordos" disponibiliza arquivos diários no formato
`YYYYMMDD_Convenios.csv`, contendo todos os convênios, contratos de repasse,
termos de parceria e demais instrumentos congêneres firmados pela Administração
Pública Federal.

## Por que Persistir

Os dados de convênios do Portal da Transparência representam a principal fonte
de transferências voluntárias da União para estados, municípios, organizações
da sociedade civil e entes privados. A persistência local permite:

1. **Consultas analíticas** — filtrar, agregar e cruzar dados de convênios sem
   depender da disponibilidade da API ou do site do Portal da Transparência.
2. **Cross-referencing com TSE** — relacionar conveniados a candidatos,
   doadores e fornecedores de campanha, permitindo análises de ligação política.
3. **Cross-referencing com TCU** — verificar se convenentes possuem contas
   irregulares, inidoneidade ou inabilitação.
4. **Cross-referencing com PNCP** — cruzar contratos públicos com convênios
   para identificar overlap de recursos.
5. **Séries históricas** — manter um histórico estável dos dados, mesmo que
   o Portal da Transparência altere seu schema ou limite o período disponível.

## Estrutura da Tabela

A tabela `convenio` armazena todas as colunas do CSV de forma denormalizada,
com índices para filtragem eficiente:

| Coluna | Tipo | Descrição |
|--------|------|-----------|
| id | UUID | Chave primária |
| numero_convenio | VARCHAR(50) | Número do convênio |
| uf | VARCHAR(2) | UF do município convenente |
| codigo_siafi_municipio | VARCHAR(20) | Código SIAFI do município |
| nome_municipio | VARCHAR(255) | Nome do município |
| situacao_convenio | VARCHAR(100) | Situação (ativo, concluído, etc.) |
| numero_original | VARCHAR(100) | Número original do instrumento |
| numero_processo | VARCHAR(100) | Número do processo administrativo |
| objeto_convenio | TEXT | Descrição do objeto do convênio |
| codigo_orgao_superior | VARCHAR(20) | Código do órgão superior |
| nome_orgao_superior | VARCHAR(255) | Nome do órgão superior |
| codigo_orgao_concedente | VARCHAR(20) | Código do órgão concedente |
| nome_orgao_concedente | VARCHAR(255) | Nome do órgão concedente |
| codigo_ug_concedente | VARCHAR(20) | Código da unidade gestora |
| nome_ug_concedente | VARCHAR(255) | Nome da unidade gestora |
| codigo_convenente | VARCHAR(20) | CPF/CNPJ do convenente |
| tipo_convenente | VARCHAR(100) | Tipo de pessoa do convenente |
| nome_convenente | VARCHAR(255) | Nome do convenente |
| tipo_ente_convenente | VARCHAR(100) | Tipo de ente (município, estado, etc.) |
| tipo_instrumento | VARCHAR(100) | Tipo (convênio, contrato de repasse, etc.) |
| valor_convenio | NUMERIC(18,2) | Valor total do convênio |
| valor_liberado | NUMERIC(18,2) | Valor total liberado |
| data_publicacao | DATE | Data de publicação |
| data_inicio_vigencia | DATE | Início da vigência |
| data_final_vigencia | DATE | Fim da vigência |
| valor_contrapartida | NUMERIC(18,2) | Valor da contrapartida |
| data_ultima_liberacao | DATE | Data da última liberação |
| valor_ultima_liberacao | NUMERIC(18,2) | Valor da última liberação |
| created_at | TIMESTAMPTZ | Data de criação do registro |
| updated_at | TIMESTAMPTZ | Data de atualização |
| deleted_at | TIMESTAMPTZ | Soft delete |

> **Nota:** A constraint unique `uq_convenio_numero` foi removida pela migration 000006, permitindo duplicatas de `numero_convenio`.

## Índices

A tabela possui índices para as principais colunas de filtragem:

- `idx_convenio_uf` — filtro por UF
- `idx_convenio_nome_municipio` — filtro por nome do município
- `idx_convenio_nome_convenente` — filtro por nome do convenente
- `idx_convenio_tipo_instrumento` — filtro por tipo de instrumento
- `idx_convenio_situacao` — filtro por situação
- `idx_convenio_objeto_trgm` — busca textual no objeto (pg_trgm)
- `idx_convenio_valores` — filtro por faixa de valores
- `idx_convenio_datas` — filtro por período

## Pipeline de Importação

O arquivo CSV de convênios é detectado automaticamente pelo LeitorCSV e
processado com **prioridade máxima (0)**, antes de qualquer dado TSE.

Fluxo:
1. `localizarArquivos()` identifica arquivos com sufixo `_convenios.csv` (minúsculas)
2. `ProcessarArquivo()` delega ao `ProcessarConvenioPortal()`
3. `processarConvenioPortal()` lê o CSV, converte valores (formato brasileiro
   de número e data) e monta estruturas `types.Convenio`
4. `PersistirDadosImportacaoPgCopy()` persiste via pgCOPY com upsert pela
   chave `numero_convenio`

## Importância

Convênios são o principal mecanismo de transferência de recursos da União para
entes subnacionais e organizações. A persistência desta base no banco local
permite:

- **Transparência fiscal**: rastrear para onde vai o dinheiro público federal
- **Análise política**: cruzar dados de convênios com informações eleitorais (TSE)
- **Controle social**: permitir que cidadãos e órgãos de controle consultem
  convenientemente os dados agregados
- **Integração com outros data sources**: PNCP, TCU, SICONFI, IBGE
