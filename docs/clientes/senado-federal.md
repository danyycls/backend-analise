# Senado Federal

**Nome cliente:** Senado Federal

**Descrição:** Cliente Go para integração com a API de Dados Abertos do Senado Federal, fornecendo acesso a senadores, comissões, processos, votações, matérias, plenário e orçamento.

## Doc Client

**Documentação de integração client:** https://legis.senado.leg.br/dadosabertos/api-docs/swagger-ui/index.html#/
**Base URL:** https://legis.senado.leg.br/dadosabertos

## APIs Integradas

### Senadores

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarSenadores | `/senador/lista/atual` | — | `[]ParlamentarResumo` | Lista os senadores em exercício no mandato corrente. |
| BuscarSenador | `/senador/{codigo}` | `codigo string` | `*ParlamentarDetalhe` | Retorna informações detalhadas sobre um senador, incluindo dados de identificação, dados básicos, telefones e outras informações. |
| ListarCargos | `/senador/{codigo}/cargos` | `codigo string` | `[]Cargo` | Lista os cargos ocupados pelo senador em comissões e outros colegiados. |
| ListarComissoes | `/senador/{codigo}/comissoes` | `codigo string` | `[]ComissaoMembro` | Lista as comissões das quais o senador é membro. |
| ListarMandatos | `/senador/{codigo}/mandatos` | `codigo string` | `[]MandatoDetalhe` | Lista os mandatos do senador, com exercícios, suplentes e partidos. |

### Comissões

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarTodasComissoes | `/comissao/lista/colegiados` | — | `[]ComissaoResumo` | Lista todos os colegiados (comissões) do Senado. |
| BuscarComissao | `/comissao/{codigo}` | `codigo string` | `*ComissaoDetalhe` | Retorna informações detalhadas sobre uma comissao, incluindo código, sigla, nome e membros. |

### Processos

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarProcessos | `/processo` | `params map[string]string` | `[]ProcessoItem` | Lista processos, com filtros opcionais. |
| ListarProcessoAssuntos | `/processo/assuntos` | — | `[]ProcessoAssunto` | Lista os assuntos cadastrados para processos. |
| ListarProcessoEmendas | `/processo/emenda` | `params map[string]string` | `[]ProcessoEmenda` | Lista emendas relacionadas a processos, com filtros opcionais. |
| BuscarProcesso | `/processo/{id}` | `id string` | `*ProcessoItem` | Retorna informações detalhadas sobre um processo. |

### Votações e Matérias

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarVotacoes | `/votacao` | `params map[string]string` | `[]VotacaoItem` | Lista votações, com filtros opcionais. |
| ListarVotacoesComissao | `/votacaoComissao/comissao/{siglaComissao}` | `siglaComissao string`, `params map[string]string` | `[]VotacaoComissao` | Lista votações de uma comissão, com filtros opcionais. |
| ListarVotacoesComissaoParlamentar | `/votacaoComissao/parlamentar/{codigo}` | `codigo string`, `params map[string]string` | `[]VotacaoComissao` | Lista votações de comissão relacionadas a um parlamentar, com filtros opcionais. |
| ListarMateriaTramitacao | `/materia/lista/tramitacao` | `params map[string]string` | `[]MateriaItem` | Lista matérias em tramitação, com filtros opcionais. |

### Plenário

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarAgendaDia | `/plenario/agenda/dia/{data}` | `data string`, `params map[string]string` | `[]Reuniao` | Lista as reuniões do plenário em uma data específica. |
| ListarAgendaMes | `/plenario/agenda/mes/{data}` | `data string`, `params map[string]string` | `[]Reuniao` | Lista as reuniões do plenário em um mês específico. |
| BuscarEncontro | `/plenario/encontro/{codigo}` | `codigo string`, `params map[string]string` | `*PlenarioEncontro` | Retorna informações sobre um encontro do plenário. |

### Orçamento

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarOrcamento | `/orcamento/lista` | — | `[]LoteEmendasOrcamento` | Lista lotes de emendas orçamentárias. |

### Documentos

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| BaixarDocumentoEmenda | `https://legis.senado.leg.br/sdleg-getter/documento?dm={id}` | `idDocumento int` | `[]byte`, `string` | Baixa o conteúdo de um documento de emenda pelo ID. Retorna o bytes do documento e o Content-Type. |
