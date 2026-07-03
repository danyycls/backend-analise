# Importacao de Convenios do Portal da Transparencia

## Origem dos Dados

Os dados sao obtidos atraves do download de arquivos CSV do Portal da Transparencia
do Governo Federal: https://portaldatransparencia.gov.br/download-de-dados

A secao "Convenios e Outros Acordos" disponibiliza arquivos diarios no formato
`YYYYMMDD_Convenios.csv`, contendo todos os convenios, contratos de repasse,
termos de parceria e demais instrumentos congêneres firmados pela Administracao
Publica Federal.

## Por que Persistir

Os dados de convenios do Portal da Transparencia representam a principal fonte
de transferencias voluntarias da Uniao para estados, municipios, organizacoes
da sociedade civil e entes privados. A persistencia local permite:

1. **Consultas analiticas** — filtrar, agregar e cruzar dados de convenios sem
   depender da disponibilidade da API ou do site do Portal da Transparencia.
2. **Cross-referencing com TSE** — relacionar conveniados a candidatos,
   doadores e fornecedores de campanha, permitindo analises de ligacao politica.
3. **Cross-referencing com TCU** — verificar se convenentes possuem contas
   irregulares, inidoneidade ou inabilitacao.
4. **Cross-referencing com PNCP** — cruzar contratos publicos com convenios
   para identificar overlap de recursos.
5. **Series historicas** — manter um historico estavel dos dados, mesmo que
   o Portal da Transparencia altere seu schema ou limite o periodo disponivel.

## Estrutura da Tabela

A tabela `convenio` armazena todas as colunas do CSV de forma denormalizada,
com indices para filtragem eficiente:

| Coluna | Tipo | Descricao |
|--------|------|-----------|
| id | UUID | Chave primaria |
| numero_convenio | VARCHAR(50) | Numero do convenio (unique) |
| uf | VARCHAR(2) | UF do municipio convenente |
| codigo_siafi_municipio | VARCHAR(20) | Codigo SIAFI do municipio |
| nome_municipio | VARCHAR(255) | Nome do municipio |
| situacao_convenio | VARCHAR(100) | Situacao (ativo, concluido, etc.) |
| numero_original | VARCHAR(100) | Numero original do instrumento |
| numero_processo | VARCHAR(100) | Numero do processo administrativo |
| objeto_convenio | TEXT | Descricao do objeto do convenio |
| codigo_orgao_superior | VARCHAR(20) | Codigo do orgao superior |
| nome_orgao_superior | VARCHAR(255) | Nome do orgao superior |
| codigo_orgao_concedente | VARCHAR(20) | Codigo do orgao concedente |
| nome_orgao_concedente | VARCHAR(255) | Nome do orgao concedente |
| codigo_ug_concedente | VARCHAR(20) | Codigo da unidade gestora |
| nome_ug_concedente | VARCHAR(255) | Nome da unidade gestora |
| codigo_convenente | VARCHAR(20) | CPF/CNPJ do convenente |
| tipo_convenente | VARCHAR(100) | Tipo de pessoa do convenente |
| nome_convenente | VARCHAR(255) | Nome do convenente |
| tipo_ente_convenente | VARCHAR(100) | Tipo de ente (municipio, estado, etc.) |
| tipo_instrumento | VARCHAR(100) | Tipo (convenio, contrato de repasse, etc.) |
| valor_convenio | NUMERIC(18,2) | Valor total do convenio |
| valor_liberado | NUMERIC(18,2) | Valor total liberado |
| data_publicacao | DATE | Data de publicacao |
| data_inicio_vigencia | DATE | Inicio da vigencia |
| data_final_vigencia | DATE | Fim da vigencia |
| valor_contrapartida | NUMERIC(18,2) | Valor da contrapartida |
| data_ultima_liberacao | DATE | Data da ultima liberacao |
| valor_ultima_liberacao | NUMERIC(18,2) | Valor da ultima liberacao |
| created_at | TIMESTAMPTZ | Data de criacao do registro |
| updated_at | TIMESTAMPTZ | Data de atualizacao |
| deleted_at | TIMESTAMPTZ | Soft delete |

## Indices

A tabela possui indices para as principais colunas de filtragem:

- `idx_convenio_uf` — filtro por UF
- `idx_convenio_nome_municipio` — filtro por nome do municipio
- `idx_convenio_nome_convenente` — filtro por nome do convenente
- `idx_convenio_tipo_instrumento` — filtro por tipo de instrumento
- `idx_convenio_situacao` — filtro por situacao
- `idx_convenio_objeto_trgm` — busca textual no objeto (pg_trgm)
- `idx_convenio_valores` — filtro por faixa de valores
- `idx_convenio_datas` — filtro por periodo

## Pipeline de Importacao

O arquivo CSV de convenios e detectado automaticamente pelo LeitorCSV e
processado com **prioridade maxima (0)**, antes de qualquer dado TSE.

Fluxo:
1. `localizarArquivos()` identifica arquivos com sufixo `_Convenios.csv`
2. `ProcessarArquivo()` delega ao `ProcessarConvenioPortal()`
3. `processarConvenioPortal()` le o CSV, converte valores (formato brasileiro
   de numero e data) e monta estruturas `types.Convenio`
4. `PersistirDadosImportacaoPgCopy()` persiste via pgCOPY com upsert pela
   chave `numero_convenio`

## Importancia

Convenios sao o principal mecanismo de transferencia de recursos da Uniao para
entes subnacionais e organizacoes. A persistencia desta base no banco local
permite:

- **Transparencia fiscal**: rastrear para onde vai o dinheiro publico federal
- **Analise politica**: cruzar dados de convenios com informacoes eleitorais (TSE)
- **Controle social**: permitir que cidadaos e orgaos de controle consultem
  convenientemente os dados agregados
- **Integracao com outros data sources**: PNCP, TCU, SICONFI, IBGE
