# APIs Integradas — Detalhamento por Cliente

Este documento detalha cada cliente HTTP externo integrado ao projeto, listando as APIs chamadas, os tipos de input/output e as interfaces Go expostas.

---

## Índice

1. [Câmara dos Deputados (Dados Abertos)](#1-camara-dos-deputados)
2. [Senado Federal (Dados Abertos)](#2-senado-federal)
3. [TCU — Certidões](#3-tcu)
4. [Portal da Transparência](#4-portal-da-transparencia)
5. [PNCP — Portal Nacional de Contratações Públicas](#5-pncp)
6. [OpenCNPJ](#6-opencnpj)
7. [IBGE — Localidades e População](#7-ibge)
8. [SICONFI — Tesouro Nacional](#8-siconfi)

---

## 1. Câmara dos Deputados

**Pacote:** `internal/shared/clients/deputados/`
**Base URL:** `https://dadosabertos.camara.leg.br/api/v2`
**Estruturas:** `DeputadosClient` (sem interface, usado diretamente)

### APIs Integradas

#### Deputados
| Método | Input | Output | Descrição |
|--------|-------|--------|-----------|
| `ListarInfoDeputadosAtivos(ctx, params)` | `map[string]string` | `[]Deputado` | Lista deputados ativos com filtros por partido, UF, legislatura |
| `BuscarDeputado(ctx, id)` | `int` (ID) | `*DeputadoDetalhe` | Detalhes completos de um deputado |
| `ListarDespesasPorDeputado(ctx, id, params)` | `int`, `map[string]string` | `[]DeputadoDespesa` | Despesas de um deputado por ano/mês |
| `ListarTodasDespesasPorDeputado(ctx, id, params)` | `int`, `map[string]string` | `[]DeputadoDespesa` | Todas as despesas (paginado automaticamente) |
| `ListarFrentesDeputado(ctx, id)` | `int` | `[]Frente` | Frentes parlamentares do deputado |
| `ListarHistorico(ctx, id)` | `int` | `[]DeputadoHistorico` | Histórico de mandatos e filiações |
| `ListarMandatosExternos(ctx, id)` | `int` | `[]DeputadoMandatoExterno` | Mandatos externos anteriores |
| `ListarOrgaos(ctx, id, params)` | `int`, `map[string]string` | `[]DeputadoOrgao` | Órgãos/comissões do deputado |

#### Legislaturas
| Método | Input | Output | Descrição |
|--------|-------|--------|-----------|
| `ListarLegislaturas(ctx)` | — | `[]Legislatura` | Todas as legislaturas |
| `BuscarLegislatura(ctx, id)` | `int` | `*Legislatura` | Detalhes de uma legislatura |
| `ListarLideres(ctx, idLeg)` | `int` | `[]Lider` | Líderes de bancada |
| `ListarMesa(ctx, idLeg)` | `int` | `[]MembroMesa` | Mesa diretora |

#### Blocos
| Método | Input | Output |
|--------|-------|--------|
| `ListarBlocos(ctx, params)` | `map[string]string` | `[]Bloco` |
| `BuscarBloco(ctx, id)` | `string` | `*Bloco` |
| `ListarPartidosDoBloco(ctx, idBloco)` | `string` | `[]Partido` |

#### Votações
| Método | Input | Output |
|--------|-------|--------|
| `ListarVotacoes(ctx, params)` | `map[string]string` | `[]Votacao` |
| `BuscarVotacao(ctx, id)` | `int` | `*VotacaoDetalhe` |
| `ListarOrientacoes(ctx, idVotacao)` | `int` | `[]Orientacao` |
| `ListarVotos(ctx, idVotacao)` | `int` | `[]Voto` |

#### Frentes
| Método | Input | Output |
|--------|-------|--------|
| `ListarFrentes(ctx, idLeg)` | `int` | `[]Frente` |
| `BuscarFrente(ctx, id)` | `int` | `*FrenteDetalhe` |
| `ListarMembrosFrente(ctx, idFrente)` | `int` | `[]MembroFrente` |

#### Grupos
| Método | Input | Output |
|--------|-------|--------|
| `ListarGrupos(ctx)` | — | `[]Grupo` |
| `BuscarGrupo(ctx, id)` | `int` | `*GrupoDetalhe` |
| `ListarHistoricoGrupo(ctx, id)` | `int` | `[]HistoricoGrupo` |
| `ListarMembrosGrupo(ctx, id)` | `int` | `[]MembroGrupo` |

#### Referências
| Método | Input | Output |
|--------|-------|--------|
| `ListarReferencias(ctx, tipo)` | `string` | `[]Referencia` |

---

## 2. Senado Federal

**Pacote:** `internal/shared/clients/senado/`
**Base URL:** `https://legis.senado.leg.br/dadosabertos`
**Estruturas:** `SenadoClient` (sem interface, usado diretamente)

### APIs Integradas

#### Senadores
| Método | Input | Output |
|--------|-------|--------|
| `ListarSenadores(ctx)` | — | `[]ParlamentarResumo` |
| `BuscarSenador(ctx, codigo)` | `string` | `*ParlamentarDetalhe` |
| `ListarCargos(ctx, codigo)` | `string` | `[]Cargo` |
| `ListarComissoes(ctx, codigo)` | `string` | `[]ComissaoMembro` |
| `ListarMandatos(ctx, codigo)` | `string` | `[]MandatoDetalhe` |

#### Orçamento
| Método | Input | Output |
|--------|-------|--------|
| `ListarOrcamento(ctx)` | — | `[]LoteEmendasOrcamento` |

#### Votações
| Método | Input | Output |
|--------|-------|--------|
| `ListarVotacoes(ctx, params)` | `map[string]string` | `[]VotacaoItem` |
| `ListarVotacoesComissao(ctx, sigla, params)` | `string`, `map` | `[]VotacaoComissao` |
| `ListarVotacoesComissaoParlamentar(ctx, codigo, params)` | `string`, `map` | `[]VotacaoComissao` |
| `ListarMateriaTramitacao(ctx, params)` | `map[string]string` | `[]MateriaItem` |

#### Processos
| Método | Input | Output |
|--------|-------|--------|
| `ListarProcessos(ctx, params)` | `map[string]string` | `[]ProcessoItem` |
| `ListarProcessoAssuntos(ctx)` | — | `[]ProcessoAssunto` |
| `ListarProcessoEmendas(ctx, params)` | `map[string]string` | `[]ProcessoEmenda` |
| `BuscarProcesso(ctx, id)` | `string` | `*ProcessoItem` |

#### Comissões
| Método | Input | Output |
|--------|-------|--------|
| `ListarTodasComissoes(ctx)` | — | `[]ComissaoResumo` |
| `BuscarComissao(ctx, codigo)` | `string` | `*ComissaoDetalhe` |

#### Plenário (Agenda)
| Método | Input | Output |
|--------|-------|--------|
| `ListarAgendaDia(ctx, data, params)` | `string`, `map` | `[]Reuniao` |
| `ListarAgendaMes(ctx, data, params)` | `string`, `map` | `[]Reuniao` |
| `BuscarEncontro(ctx, codigo, params)` | `string`, `map` | `*PlenarioEncontro` |

---

## 3. TCU

**Pacote:** `internal/shared/clients/tcu/`
**Base URL:** `https://certidoes.apps.gov.br/api/publico`
**Estruturas:** `Client` (interface), `TCUClient` (implementação)

### Interface

```go
type Client interface {
    BuscarContasIrregulares(ctx, filter) ([]ContasIrregulares, error)
    BuscarInabilitados(ctx, filter) ([]Sancoes, error)
    BuscarInidoneos(ctx, filter) ([]Sancoes, error)
    BuscarFinsEleitorais(ctx, filter) ([]FinsEleitorais, error)
}
```

### Input comum

`TCUQueryParams`:
```go
type TCUQueryParams struct {
    ParteNome string  // nome parcial da pessoa
    CPF       string  // CPF sem formatação
    CNPJ      string  // CNPJ sem formatação
    UF        string  // sigla da UF
    Municipio string  // nome do município
}
```

### Outputs

| Tipo | Campos principais |
|------|-------------------|
| `ContasIrregulares` | `numeroProcessoFormatado`, `nome`, `tipoRegistro`, `numeroRegistro`, `municipio`, `uf`, `dataTransitoEmJulgado`, `linkDeliberacoesProcesso`, `linkAcompanhamentoProcesso` |
| `FinsEleitorais` | Mesmos + `dataFinalFinsEleitorais` |
| `Sancoes` | `numeroProcessoFormatado`, `nome`, `tipoRegistro`, `numeroRegistro`, `municipio`, `uf`, `numeroAcordaoFormatado`, `dataAcordao`, `dataTransitoEmJulgado`, `dataFinalSancao` |

**Mock gerado:** `go:generate mockgen` em `interface.go` → `mock.go`

---

## 4. Portal da Transparência

**Pacote:** `internal/shared/clients/portaltransparencia/`
**Base URL:** `https://api.portaldatransparencia.gov.br`
**Autenticação:** `chave-api-dados` no header (API Key)
**Estruturas:** `PortalTransparenciaClient` (sem interface)

### APIs Integradas

#### Órgãos
| Método | Input | Output |
|--------|-------|--------|
| `ListarOrgaosSIAPE(ctx, filtro)` | `OrgaoQueryParams` | `[]Orgao` |
| `ListarOrgaosSIAFI(ctx, filtro)` | `OrgaoQueryParams` | `[]Orgao` |

**Input:**
```go
type OrgaoQueryParams struct {
    Pagina    int
    Codigo    string
    Descricao string
}
```

#### Pessoas
| Método | Input | Output |
|--------|-------|--------|
| `ListarPessoasFisicas(ctx, filtro)` | `PessoaFisicaQueryParams` | `*PessoaFisica` |
| `ListarPessoasJuridicas(ctx, filtro)` | `PessoaJuridicaQueryParams` | `*PessoaJuridica` |

#### Cartões Corporativos
| Método | Input | Output |
|--------|-------|--------|
| `ListarCartoes(ctx, filtro)` | `CartaoQueryParams` | `[]Cartao` |

**Input:**
```go
type CartaoQueryParams struct {
    Pagina, MesExtratoInicio, MesExtratoFim, DataTransacaoInicio,
    DataTransacaoFim, TipoCartao, CodigoOrgao, CPFPortador,
    CPFCNPJFavorecido, ValorDe, ValorAte
}
```

#### Servidores
| Método | Input | Output |
|--------|-------|--------|
| `ListarServidores(ctx, filtro)` | `ServidorQueryParams` | `[]CadastroServidor` |
| `BuscarServidorPorID(ctx, id)` | `int` | `*CadastroServidor` |
| `ListarRemuneracaoServidores(ctx, filtro)` | `ServidorRemuneracaoQueryParams` | `[]ServidorRemuneracao` |
| `ListarServidoresPorOrgao(ctx, filtro)` | `ServidorPorOrgaoQueryParams` | `[]ServidorPorOrgao` |
| `ListarFuncoesECargos(ctx, filtro)` | `FuncaoCargoQueryParams` | `[]FuncaoServidor` |
| `ListarPEPs(ctx, filtro)` | `PEPQueryParams` | `[]PEP` |

#### Despesas
| Método | Input | Output |
|--------|-------|--------|
| `ListarRecursosRecebidos(ctx, filtro)` | `DespesaRecursosRecebidosQueryParams` | `[]PessoaRecursosRecebidosUGMesDesnormalizada` |
| `ListarDespesasPorOrgao(ctx, filtro)` | `DespesaPorOrgaoQueryParams` | `[]DespesaAnualPorOrgao` |
| `ListarDespesasPorFuncionalProgramatica(ctx, filtro)` | `DespesaFuncionalProgramaticaQueryParams` | `[]DespesaAnualPorFuncaoESubfuncao` |
| `ListarDespesasMovimentacaoLiquida(ctx, filtro)` | `DespesaMovimentacaoLiquidaQueryParams` | `[]DespesaLiquidaAnualPorFuncaoESubfuncao` |
| `ListarDespesasPlanoOrcamentario(ctx, filtro)` | `DespesaPlanoOrcamentarioQueryParams` | `[]DespesasPorPlanoOrcamentario` |
| `ListarItensEmpenho(ctx, codDoc, pagina)` | `string`, `int` | `[]DetalhamentoDoGasto` |
| `ListarHistoricoItemEmpenho(ctx, codDoc, seq, pagina)` | `string`, `int`, `int` | `[]HistoricoSubItemEmpenho` |
| `ListarSubfuncoes(ctx, filtro)` | `ListarFuncionalProgramaticaQueryParams` | `[]Subfuncao` |
| `ListarProgramas(ctx, filtro)` | `ListarFuncionalProgramaticaQueryParams` | `[]CodigoDescricao` |
| `ListarFuncionalProgramatica(ctx, ano, pagina)` | `int`, `int` | `[]FuncionalProgramatica` |
| `ListarFuncoes(ctx, filtro)` | `ListarFuncionalProgramaticaQueryParams` | `[]Funcao` |
| `ListarAcoes(ctx, filtro)` | `ListarFuncionalProgramaticaQueryParams` | `[]CodigoDescricao` |
| `ListarFavorecidosFinaisPorDocumento(ctx, codDoc, pagina)` | `string`, `int` | `[]ConsultaFavorecidosFinaisPorDocumento` |
| `ListarEmpenhosImpactados(ctx, codDoc, fase, pagina)` | `string`, `string`, `int` | `[]EmpenhoImpactadoBasico` |
| `ListarDocumentos(ctx, filtro)` | `DespesaDocumentosQueryParams` | `[]interface{}` |
| `BuscarDocumentoPorCodigo(ctx, codigo)` | `string` | `*DespesasPorDocumento` |
| `ListarDocumentosRelacionados(ctx, codDoc, fase)` | `string`, `string` | `[]DocumentoRelacionado` |
| `ListarDocumentosPorFavorecido(ctx, filtro)` | `DespesaDocumentosPorFavorecidoQueryParams` | `[]interface{}` |
| `ListarTiposTransferencia(ctx)` | — | `[]CodigoDescricao` |

#### Emendas
| Método | Input | Output |
|--------|-------|--------|
| `ListarEmendas(ctx, filtro)` | `EmendaQueryParams` | `[]ConsultaEmendas` |
| `ListarDocumentosEmenda(ctx, codigo, pagina)` | `string`, `int` | `[]DocumentoRelacionadoEmenda` |

---

## 5. PNCP

**Pacote:** `internal/shared/clients/pncp/`
**Base URL:** `https://pncp.gov.br/pncp-consulta/v1`
**Estruturas:** `PNCPClient` (sem interface)

### APIs Integradas

| Método | Input | Output |
|--------|-------|--------|
| `BuscarContratos(ctx, cnpj, dataInicial, dataFinal, pagina, tamanho)` | `string, string, string, int, int` | `[]Contrato` |
| `BuscarContratacoesPorMunicipio(ctx, codMun, dataInicial, dataFinal, pagina, tamanho)` | `string, string, string, int, int` | `*PublicacaoResponse` |
| `BuscarContratacoesPorUF(ctx, uf, dataInicial, dataFinal, pagina, tamanho)` | `string, string, string, int, int` | `*PublicacaoResponse` |

**Inputs principais:**
- `AnalisePublicacaoRequest`: `{ tipo, uf, codigo_municipio_ibge, data_inicial, data_final }`
- `AmparoLegal`: `{ codigo, nome, descricao }`

**Output principal:**
- `PublicacaoResponse`: `{ data, totalRegistros, totalPaginas, numeroPagina, paginasRestantes, empty }`

O PNCP também é consumido via **WebSocket** para streaming (`SSE`) nos endpoints de análise de órgãos e publicações.

---

## 6. OpenCNPJ

**Pacote:** `internal/shared/clients/opencnpj/`
**Base URL:** `https://api.opencnpj.org/%s` (o `%s` é substituído pelo CNPJ)
**Estruturas:** `Client` (interface), `OpenCNPJClient` (implementação)

### Interface

```go
type Client interface {
    Buscar(ctx, cnpj string) (*types.OpenCNPJResponse, error)
}
```

### Output

```go
type OpenCNPJResponse struct {
    CNPJ              string
    RazaoSocial       string
    NomeFantasia      string
    SituacaoCadastral string
    CapitalSocial     string
    Socios            []Socio  // qsa
}
```

**Mock gerado:** `go:generate mockgen` em `interface.go` → `mock.go`

---

## 7. IBGE

**Pacote:** `internal/shared/clients/ibge/`
**Base URL Localidades:** `https://servicodados.ibge.gov.br/api/v1/localidades`
**Base URL Agregados:** `https://servicodados.ibge.gov.br/api/v3/agregados`
**Estruturas:** `IBGEClient` (sem interface)

### APIs Integradas

| Método | Input | Output |
|--------|-------|--------|
| `ListarEstados(ctx)` | — | `[]EstadoIBGE` |
| `ListarMunicipios(ctx, uf)` | `string` (UF) | `[]MunicipioIBGE` |
| `ListarMunicipiosCompleto(ctx)` | — | `[]MunicipioDetalhadoIBGE` |
| `BuscarPopulacao(ctx, municipioIDs)` | `[]int` | `map[int]int64` (ID → população) |

**Outputs:**
```go
type EstadoIBGE struct { ID int; Sigla string; Nome string }
type MunicipioIBGE struct { ID int; Nome string }
type MunicipioDetalhadoIBGE struct {
    ID int; Nome string; Microrregiao Microrregiao
}
```

---

## 8. SICONFI

**Pacote:** `internal/shared/clients/siconfi/`
**Base URL:** `https://apidatalake.tesouro.gov.br/ords/siconfi/tt`
**Rate limit:** 1 req/s
**Paginação:** 5000 itens/página
**Estruturas:** `SICONFIClient` (sem interface)

### APIs Integradas

| Método | Input | Output |
|--------|-------|--------|
| `ListarEntes(ctx)` | — | `[]Ente` |
| `BuscarDCA(ctx, anExercicio, idEnte, noAnexo?)` | `int64, int, ...string` | `[]DCAItem` |
| `BuscarRGF(ctx, params)` | `RGFParams` | `[]RGFItem` |
| `BuscarRREO(ctx, params)` | `RREOParams` | `[]RREOItem` |
| `BuscarMSCPatrimonial(ctx, params)` | `MSCParams` | `[]MSCItem` |
| `BuscarMSCOrcamentaria(ctx, params)` | `MSCParams` | `[]MSCItem` |
| `BuscarMSCControle(ctx, params)` | `MSCParams` | `[]MSCItem` |
| `BuscarExtratoEntregas(ctx, idEnte, anReferencia)` | `int, int64` | `[]ExtratoEntregasItem` |
| `ListarAnexosRelatorios(ctx)` | — | `[]AnexoRelatorio` |

### Envelope de resposta

Todas as respostas SICONFI seguem o formato paginado:
```go
type Response[T any] struct {
    Items   []T
    HasMore bool
    Count   int
    Limit   int
    Offset  int
    Links   []Link
}
```

### Inputs

```go
type RGFParams struct {
    AnExercicio         int64   // obrigatório
    InPeriodicidade     string  // Q (quadrimestral) ou S (semestral)
    NrPeriodo           int     // 1-3 (quadr.) ou 1-2 (semestral)
    CoTipoDemonstrativo string  // ex: "RGF"
    CoPoder             string  // E, L, J, M, D
    IdEnte              int     // código IBGE 7 dígitos
    NoAnexo             string  // opcional
    CoEsfera            string  // opcional: M, E, U, C
}

type RREOParams struct {
    AnExercicio         int64
    NrPeriodo           int     // 1-6 (bimestre)
    CoTipoDemonstrativo string  // ex: "RREO"
    IdEnte              int
    NoAnexo             string  // opcional
    CoEsfera            string  // opcional
}

type MSCParams struct {
    IdEnte        int
    AnReferencia  int64
    MeReferencia  int64     // 1-12 ou 13 (encerramento)
    CoTipoMatriz  string    // MSCC ou MSCE
    ClasseConta   int       // 1-8 (PCASP)
    IdTV          string    // beginning_balance, ending_balance, period_change
}
```

---

## Padrão de Mocks

Clientes que possuem mock gerado (via `go:generate mockgen`):

| Cliente | Arquivo de interface | Arquivo mock |
|---------|---------------------|--------------|
| OpenCNPJ | `internal/shared/clients/opencnpj/interface.go` | `mock.go` |
| TCU | `internal/shared/clients/tcu/interface.go` | `mock.go` |
| Redis | `internal/shared/redis/interface.go` | `mock.go` |
| Ligação Política Usecase | `internal/ligacao-politica/usecase/interface.go` | `mock.go` |
