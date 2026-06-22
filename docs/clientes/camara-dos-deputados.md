# Câmara dos Deputados

**Nome cliente:** Câmara dos Deputados

**Descrição:** Cliente Go para integração com a API de Dados Abertos da Câmara dos Deputados (RESTful v2), fornecendo acesso a deputados, despesas, frentes, blocos, grupos, legislaturas, votações e referências.

## Doc Client

**Documentação de integração client:** https://dadosabertos.camara.leg.br/swagger/api.html?tab=api
**Base URL:** https://dadosabertos.camara.leg.br/api/v2

## APIs Integradas

### Deputados

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarInfoDeputadosAtivos | `/deputados` | `params map[string]string` | `[]Deputado` | Lista de deputados, com filtros opcionais por nome, sigla de partido, sigla de UF, etc. |
| BuscarDeputado | `/deputados/{id}` | `id int` | `*DeputadoDetalhe` | Retorna informações detalhadas sobre um deputado, incluindo dados civis e último status. |
| ListarDespesasPorDeputado | `/deputados/{id}/despesas` | `idDeputado int`, `params map[string]string` | `[]DeputadoDespesa` | Lista as despesas cobertas pela Cota para Exercício da Atividade Parlamentar de um deputado (primeira página). |
| ListarTodasDespesasPorDeputado | `/deputados/{id}/despesas` | `idDeputado int`, `params map[string]string` | `[]DeputadoDespesa` | Lista todas as despesas de um deputado, percorrendo automaticamente a paginação via `links.next`. |
| ListarFrentesDeputado | `/deputados/{id}/frentes` | `idDeputado int` | `[]Frente` | Lista as frentes parlamentares das quais um deputado é integrante. |
| ListarHistorico | `/deputados/{id}/historico` | `idDeputado int` | `[]DeputadoHistorico` | Lista o histórico de mudanças de status do deputado. |
| ListarMandatosExternos | `/deputados/{id}/mandatosExternos` | `idDeputado int` | `[]DeputadoMandatoExterno` | Lista mandatos eletivos desempenhados pelo deputado fora da Câmara dos Deputados. |
| ListarOrgaos | `/deputados/{id}/orgaos` | `idDeputado int`, `params map[string]string` | `[]DeputadoOrgao` | Lista os órgãos dos quais um deputado é integrante. |

### Blocos

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarBlocos | `/blocos` | `params map[string]string` | `[]Bloco` | Lista os blocos partidários, com filtros opcionais. |
| BuscarBloco | `/blocos/{id}` | `id string` | `*Bloco` | Retorna informações detalhadas sobre um bloco partidário. |
| ListarPartidosDoBloco | `/blocos/{id}/partidos` | `idBloco string` | `[]Partido` | Lista os partidos que compõem um bloco partidário. |

### Frentes

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarFrentes | `/frentes` | `idLegislatura int` | `[]Frente` | Lista as frentes parlamentares de uma legislatura. |
| BuscarFrente | `/frentes/{id}` | `id int` | `*FrenteDetalhe` | Retorna informações detalhadas sobre uma frente parlamentar. |
| ListarMembrosFrente | `/frentes/{id}/membros` | `idFrente int` | `[]MembroFrente` | Lista os integrantes de uma frente parlamentar. |

### Grupos

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarGrupos | `/grupos` | — | `[]Grupo` | Lista os grupos parlamentares. |
| BuscarGrupo | `/grupos/{id}` | `id int` | `*GrupoDetalhe` | Retorna informações detalhadas sobre um grupo parlamentar. |
| ListarHistoricoGrupo | `/grupos/{id}/historico` | `id int` | `[]HistoricoGrupo` | Lista o histórico de instalações e trocas de presidente de um grupo. |
| ListarMembrosGrupo | `/grupos/{id}/membros` | `id int` | `[]MembroGrupo` | Lista os parlamentares que foram membros de um grupo parlamentar. |

### Legislaturas

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarLegislaturas | `/legislaturas` | — | `[]Legislatura` | Lista os períodos de trabalho da Câmara, com datas de início e fim. |
| BuscarLegislatura | `/legislaturas/{id}` | `id int` | `*Legislatura` | Retorna informações detalhadas sobre uma legislatura. |
| ListarLideres | `/legislaturas/{id}/lideres` | `idLegislatura int` | `[]Lider` | Lista os líderes de bancada de uma legislatura. |
| ListarMesa | `/legislaturas/{id}/mesa` | `idLegislatura int` | `[]MembroMesa` | Lista os integrantes da Mesa Diretora em uma legislatura. |

### Referências

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarReferencias | `/referencias/{tipo}` | `tipo string` | `[]Referencia` | Lista valores de domínio (tabelas de referência) de um tipo, como tipos de eventos, órgãos, etc. |

### Votações

| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarVotacoes | `/votacoes` | `params map[string]string` | `[]Votacao` | Lista votações realizadas, com filtros opcionais. |
| BuscarVotacao | `/votacoes/{id}` | `id int` | `*VotacaoDetalhe` | Retorna informações detalhadas sobre uma votação. |
| ListarOrientacoes | `/votacoes/{id}/orientacoes` | `idVotacao int` | `[]Orientacao` | Lista as orientações das bancadas em uma votação. |
| ListarVotos | `/votacoes/{id}/votos` | `idVotacao int` | `[]Voto` | Lista os votos individuais dos parlamentares em uma votação. |
