# Portal da Transparência

**Nome cliente:** Portal da Transparência

**Descrição:** Cliente Go que fornece acesso à API pública do Portal da Transparência (https://api.portaldatransparencia.gov.br). Permite consultas a órgãos (SIAPE e SIAFI), pessoas físicas e jurídicas, despesas (recursos recebidos, despesas por órgão, por funcional‑programática, movimentação líquida, plano orçamentário, itens de empenho, histórico de itens, tipos de transferência, documentos e buscas por documento), cartões, emendas e servidores (dados cadastrais, remuneração, funções, cargos, PEPs, etc.).

## Doc Client

**Documentação de integração client:** Esta seção descreve, em português, os métodos públicos do cliente Go, seus parâmetros de entrada e tipos de retorno, além dos caminhos relativos à Base URL.

**Base URL:** https://api.portaldatransparencia.gov.br

## APIs Integradas

### Órgãos
| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarOrgaosSIAPE | /api-de-dados/orgaos-siape | OrgaoQueryParams | []Orgao | Lista órgãos da base SIAPE (swagger: “Lista órgãos SIAPE”). |
| ListarOrgaosSIAFI | /api-de-dados/orgaos-siafi | OrgaoQueryParams | []Orgao | Lista órgãos da base SIAFI (swagger: “Lista órgãos SIAFI”). |

### Pessoas
| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarPessoasFisicas | /api-de-dados/pessoa-fisica | PessoaFisicaQueryParams | *PessoaFisica | Consulta pessoa física com filtros avançados. |
| ListarPessoasJuridicas | /api-de-dados/pessoa-juridica | PessoaJuridicaQueryParams | *PessoaJuridica | Consulta pessoa jurídica com filtros avançados. |

### Cartões
| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarCartoes | /api-de-dados/cartoes | CartaoQueryParams | []Cartao | Recupera cartões de crédito/debito com filtragem por período, tipo, valores, etc. |

### Emendas
| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarEmendas | /api-de-dados/emendas | EmendaQueryParams | []ConsultaEmendas | Consulta emendas parlamentares com filtros por código, número, autor, tipo, ano, etc. |
| ListarDocumentosEmenda | /api-de-dados/emendas/documentos/{codigo} | pagina (int) – via parâmetro de consulta | []DocumentoRelacionadoEmenda | Obtém documentos associados a uma emenda específica. |

### Servidores
| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarServidores | /api-de-dados/servidores | ServidorQueryParams | []CadastroServidor | Lista servidores cadastrados com filtros por CPF, nome, órgão, situação, etc. |
| BuscarServidorPorID | /api-de-dados/servidores/{id} | — | *CadastroServidor | Busca servidor pelo identificador numérico. |
| ListarRemuneracaoServidores | /api-de-dados/servidores/remuneracao | ServidorRemuneracaoQueryParams | []ServidorRemuneracao | Consulta remuneração de servidores (mes/ano, CPF, ID). |
| ListarServidoresPorOrgao | /api-de-dados/servidores/por-orgao | ServidorPorOrgaoQueryParams | []ServidorPorOrgao | Lista servidores por órgão de lotação ou exercício. |
| ListarFuncoesECargos | /api-de-dados/servidores/funcoes-e-cargos | FuncaoCargoQueryParams | []FuncaoServidor | Consulta funções e cargos de servidores com paginação. |
| ListarPEPs | /api-de-dados/peps | PEPQueryParams | []PEP | Consulta Pessoas Expostas Politicamente (PEPs) com filtros. |

### Despesas
| Método | URL | Input | Output | Descrição |
|--------|-----|-------|--------|-----------|
| ListarRecursosRecebidos | /api-de-dados/despesas/recursos-recebidos | DespesaRecursosRecebidosQueryParams | []PessoaRecursosRecebidosUGMesDesnormalizada | Lista recursos recebidos por pessoa/UG/mês. |
| ListarDespesasPorOrgao | /api-de-dados/despesas/por-orgao | DespesaPorOrgaoQueryParams | []DespesaAnualPorOrgao | Despesas agregadas por órgão e ano. |
| ListarDespesasPorFuncionalProgramatica | /api-de-dados/despesas/por-funcional-programatica | DespesaFuncionalProgramaticaQueryParams | []DespesaAnualPorFuncaoESubfuncao | Despesas por função/subfunção/programa/ação. |
| ListarDespesasMovimentacaoLiquida | /api-de-dados/despesas/por-funcional-programatica/movimentacao-liquida | DespesaMovimentacaoLiquidaQueryParams | []DespesaLiquidaAnualPorFuncaoESubfuncao | Movimentação líquida de despesas por critérios programáticos. |
| ListarDespesasPlanoOrcamentario | /api-de-dados/despesas/plano-orcamentario | DespesaPlanoOrcamentarioQueryParams | []DespesasPorPlanoOrcamentario | Consulta despesas por plano orçamentário. |
| ListarItensEmpenho | /api-de-dados/despesas/itens-de-empenho | codigoDocumento (string), pagina (int) – via query | []DetalhamentoDoGasto | Detalha itens de empenho de um documento. |
| ListarHistoricoItemEmpenho | /api-de-dados/despesas/itens-de-empenho/historico | codigoDocumento (string), sequencial (int), pagina (int) – via query | []HistoricoSubItemEmpenho | Histórico de alterações de um item de empenho. |
| ListarSubfuncoes | /api-de-dados/despesas/funcional-programatica/subfuncoes | ListarFuncionalProgramaticaQueryParams | []Subfuncao | Lista subfunções disponíveis. |
| ListarProgramas | /api-de-dados/despesas/funcional-programatica/programs | ListarFuncionalProgramaticaQueryParams | []CodigoDescricao | Lista programas disponíveis. |
| ListarFuncionalProgramatica | /api-de-dados/despesas/funcional-programatica/listar | ano (int), pagina (int) – via query | []FuncionalProgramatica | Lista todas as combinações funcional‑programáticas para o ano. |
| ListarFuncoes | /api-de-dados/despesas/funcional-programatica/funcoes | ListarFuncionalProgramaticaQueryParams | []Funcao | Lista funções disponíveis. |
| ListarAcoes | /api-de-dados/despesas/funcional-programatica/acoes | ListarFuncionalProgramaticaQueryParams | []CodigoDescricao | Lista ações disponíveis. |
| ListarFavorecidosFinaisPorDocumento | /api-de-dados/despesas/favorecidos-finais-por-documento | codigoDocumento (string), pagina (int) – via query | []ConsultaFavorecidosFinaisPorDocumento | Favorecidos finais associados a um documento. |
| ListarEmpenhosImpactados | /api-de-dados/despesas/empenhos-impactados | codigoDocumento (string), fase (string), pagina (int) – via query | []EmpenhoImpactadoBasico | Empenhos impactados por um documento e fase. |
| ListarDocumentos | /api-de-dados/despesas/documentos | DespesaDocumentosQueryParams | []interface{} | Lista documentos de despesas com filtros de data, fase, unidade gestora, etc. |
| BuscarDocumentoPorCodigo | /api-de-dados/despesas/documentos/{codigo} | — | *DespesasPorDocumento | Busca documento de despesa por código. |
| ListarDocumentosRelacionados | /api-de-dados/despesas/documentos-relacionados | codigoDocumento (string), fase (string) – via query | []DocumentoRelacionado | Documentos relacionados ao documento especificado. |
| ListarDocumentosPorFavorecido | /api-de-dados/despesas/documentos-por-favorecido | DespesaDocumentosPorFavorecidoQueryParams | []interface{} | Documentos de despesa filtrados por favorecido. |
| ListarTiposTransferencia | /api-de-dados/despesas/tipo-transferencia | — | []CodigoDescricao | Tipos de transferência disponíveis. |

**Total de endpoints documentados:** 32
**Número de subtabelas (recursos):** 6
