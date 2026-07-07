package parse

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
)

func uuidStr(id *uuid.UUID) string {
	if id == nil {
		return "<nil>"
	}
	return id.String()
}

// ValidarFKsEmMemoria percorre os dados acumulados e verifica que cada FK
// aponta para uma entidade existente nos mapas de dimensoes. Nao bloqueia
// a persistencia — apenas loga violacoes. Retorna o total encontrado.
func indexFromMap[K comparable, V any](m map[K]V, idFn func(V) uuid.UUID) map[uuid.UUID]bool {
	idx := make(map[uuid.UUID]bool, len(m))
	for _, v := range m {
		idx[idFn(v)] = true
	}
	return idx
}

func checkFK(idx map[uuid.UUID]bool, id *uuid.UUID) bool {
	return id == nil || idx[*id]
}

func checkFKDirect(idx map[uuid.UUID]bool, id uuid.UUID) bool {
	return idx[id]
}

func ValidarFKsEmMemoria(dados *tipos.DadosImportacao) int {
	if dados == nil {
		return 0
	}
	log := logger.New("LeitorCSV: ValidarFKsEmMemoria")
	violacoes := 0

	ei := indexFromMap(dados.Eleicoes, func(e *types.Eleicao) uuid.UUID { return e.ID })
	pi := indexFromMap(dados.Partidos, func(e *types.Partido) uuid.UUID { return e.ID })
	ui := indexFromMap(dados.UnidadesEleitorais, func(e *types.UnidadeEleitoral) uuid.UUID { return e.ID })
	fi := indexFromMap(dados.Fornecedores, func(e *types.Fornecedor) uuid.UUID { return e.ID })
	di := indexFromMap(dados.Doadores, func(e *types.Doador) uuid.UUID { return e.ID })
	ci := indexFromMap(dados.Candidatos, func(e *types.Candidato) uuid.UUID { return e.ID })
	rci := indexFromMap(dados.ReceitasCandidatoPorSQ, func(e *types.ReceitaCandidato) uuid.UUID { return e.ID })
	roi := indexFromMap(dados.ReceitasOrgaoPorSQ, func(e *types.ReceitaOrgaoPartidario) uuid.UUID { return e.ID })

	for _, c := range dados.Candidatos {
		if !checkFKDirect(ei, c.EleicaoID) {
			log.Warn("FK invalida: candidato eleicao_id nao encontrada", "sq_candidato", c.SQCandidato, "eleicao_id", c.EleicaoID)
			violacoes++
		}
		if !checkFK(pi, c.PartidoID) {
			log.Warn("FK invalida: candidato partido_id nao encontrada", "sq_candidato", c.SQCandidato, "partido_id", *c.PartidoID)
			violacoes++
		}
	}

	for _, p := range dados.Prestacoes {
		if !checkFKDirect(ei, p.EleicaoID) {
			log.Warn("FK invalida: prestacao_contas eleicao_id nao encontrada", "sq_prestador_contas", p.SQPrestadorContas, "eleicao_id", p.EleicaoID)
			violacoes++
		}
		if !checkFK(ci, p.CandidatoID) {
			log.Warn("FK invalida: prestacao_contas candidato_id nao encontrada", "sq_prestador_contas", p.SQPrestadorContas, "candidato_id", uuidStr(p.CandidatoID))
			violacoes++
		}
		if !checkFK(pi, p.PartidoID) {
			log.Warn("FK invalida: prestacao_contas partido_id nao encontrada", "sq_prestador_contas", p.SQPrestadorContas, "partido_id", uuidStr(p.PartidoID))
			violacoes++
		}
		if !checkFK(ui, p.UnidadeEleitoralID) {
			log.Warn("FK invalida: prestacao_contas unidade_eleitoral_id nao encontrada", "sq_prestador_contas", p.SQPrestadorContas, "unidade_eleitoral_id", uuidStr(p.UnidadeEleitoralID))
			violacoes++
		}
	}

	for _, d := range dados.DespesasCandidato {
		if !checkFKDirect(ci, d.CandidatoID) {
			log.Warn("FK invalida: despesa_candidato candidato_id nao encontrada", "sq_despesa", d.SQDespesa, "tipo_registro", d.TipoRegistro, "candidato_id", d.CandidatoID)
			violacoes++
		}
		if !checkFK(fi, d.FornecedorID) {
			log.Warn("FK invalida: despesa_candidato fornecedor_id nao encontrada", "sq_despesa", d.SQDespesa, "tipo_registro", d.TipoRegistro, "fornecedor_id", uuidStr(d.FornecedorID))
			violacoes++
		}
	}

	for _, d := range dados.DespesasOrgaoPartidario {
		if !checkFKDirect(pi, d.PartidoID) {
			log.Warn("FK invalida: despesa_orgao_partidario partido_id nao encontrada", "sq_despesa", d.SQDespesa, "tipo_registro", d.TipoRegistro, "partido_id", d.PartidoID)
			violacoes++
		}
		if !checkFK(fi, d.FornecedorID) {
			log.Warn("FK invalida: despesa_orgao_partidario fornecedor_id nao encontrada", "sq_despesa", d.SQDespesa, "tipo_registro", d.TipoRegistro, "fornecedor_id", uuidStr(d.FornecedorID))
			violacoes++
		}
	}

	for _, r := range dados.ReceitasCandidato {
		if !checkFKDirect(ci, r.CandidatoID) {
			log.Warn("FK invalida: receita_candidato candidato_id nao encontrada", "sq_receita", r.SQReceita, "candidato_id", r.CandidatoID)
			violacoes++
		}
		if !checkFK(di, r.DoadorID) {
			log.Warn("FK invalida: receita_candidato doador_id nao encontrada", "sq_receita", r.SQReceita, "doador_id", uuidStr(r.DoadorID))
			violacoes++
		}
	}

	for _, r := range dados.ReceitasOrgaoPartidario {
		if !checkFKDirect(pi, r.PartidoID) {
			log.Warn("FK invalida: receita_orgao_partidario partido_id nao encontrada", "sq_receita", r.SQReceita, "partido_id", r.PartidoID)
			violacoes++
		}
		if !checkFK(di, r.DoadorID) {
			log.Warn("FK invalida: receita_orgao_partidario doador_id nao encontrada", "sq_receita", r.SQReceita, "doador_id", uuidStr(r.DoadorID))
			violacoes++
		}
	}

	for _, b := range dados.BensCandidato {
		if !checkFKDirect(ci, b.CandidatoID) {
			log.Warn("FK invalida: bem_candidato candidato_id nao encontrada", "numero_ordem", b.NumeroOrdem, "candidato_id", b.CandidatoID)
			violacoes++
		}
	}

	for _, o := range dados.ReceitasDoadorOriginarioCandidato {
		if !checkFK(rci, o.ReceitaCandidatoID) {
			log.Warn("FK invalida: receita_doador_originario_candidato receita_candidato_id nao encontrada", "sq_receita", o.SQReceita, "receita_candidato_id", uuidStr(o.ReceitaCandidatoID))
			violacoes++
		}
	}

	for _, o := range dados.ReceitasDoadorOriginarioOrgaoPartidario {
		if !checkFK(roi, o.ReceitaOrgaoPartidarioID) {
			log.Warn("FK invalida: receita_doador_originario_orgao_partidario receita_orgao_partidario_id nao encontrada", "sq_receita", o.SQReceita, "receita_orgao_partidario_id", uuidStr(o.ReceitaOrgaoPartidarioID))
			violacoes++
		}
	}

	if violacoes > 0 {
		log.Error("violacoes de FK encontradas em memoria", "total", violacoes)
	}

	return violacoes
}

// ObterUFDoNomeArquivo extrai a UF (sigla) do nome do arquivo CSV
// Ex: "prestacao_contas_SP.csv" -> "SP"
// Para arquivos nacionais (sem UF no nome), retorna "BR".
func ObterUFDoNomeArquivo(nome string) string {
	nomeSemExt := strings.TrimSuffix(nome, filepath.Ext(nome))
	partes := strings.Split(nomeSemExt, "_")
	// Remove part suffix (e.g., "part001") from split filenames
	if len(partes) > 1 && strings.HasPrefix(partes[len(partes)-1], "part") {
		partes = partes[:len(partes)-1]
	}
	if len(partes) == 0 {
		return "BR"
	}
	sigla := strings.ToUpper(partes[len(partes)-1])
	if !ufsBrasil[sigla] {
		return "BR"
	}
	return sigla
}

var ufsBrasil = map[string]bool{
	"AC": true, "AL": true, "AP": true, "AM": true, "BA": true, "CE": true,
	"DF": true, "ES": true, "GO": true, "MA": true, "MT": true, "MS": true,
	"MG": true, "PA": true, "PB": true, "PR": true, "PE": true, "PI": true,
	"RJ": true, "RN": true, "RS": true, "RO": true, "RR": true, "SC": true,
	"SP": true, "SE": true, "TO": true,
}

// LimparDadosAposPersistencia libera memoria dos registros transacionais
// mantendo as dimensoes (eleicoes, partidos, etc.) para deduplicacao entre arquivos
func LimparDadosAposPersistencia(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}

	dados.Convenios = nil
	dados.DespesasCandidato = nil
	dados.DespesasOrgaoPartidario = nil
	dados.ReceitasCandidato = nil
	dados.ReceitasOrgaoPartidario = nil
	dados.ReceitasDoadorOriginarioCandidato = nil
	dados.ReceitasDoadorOriginarioOrgaoPartidario = nil
	dados.BensCandidato = nil
}

// LimparTodosDados libera toda a memoria dos dados, incluindo dimensoes
func LimparTodosDados(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}

	dados.Convenios = nil
	dados.Eleicoes = nil
	dados.UnidadesEleitorais = nil
	dados.Partidos = nil
	dados.Candidatos = nil
	dados.CandidatosPorID = nil
	dados.Fornecedores = nil
	dados.Doadores = nil
	dados.Prestacoes = nil
	dados.PrestacoesPorTipoESQ = nil
	dados.ReceitasCandidatoPorSQ = nil
	dados.ReceitasOrgaoPorSQ = nil
	dados.DespesasCandidato = nil
	dados.DespesasOrgaoPartidario = nil
	dados.ReceitasCandidato = nil
	dados.ReceitasOrgaoPartidario = nil
	dados.ReceitasDoadorOriginarioCandidato = nil
	dados.ReceitasDoadorOriginarioOrgaoPartidario = nil
	dados.BensCandidato = nil
}

// ---------------------------------------------------------------------------
// Funcoes de remapeamento: quando um registro em memoria tem ID temporario
// (gerado pelo processador), substitui pelo ID real do banco em todos os
// registros dependentes
// ---------------------------------------------------------------------------

func remapearEleicaoIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	if len(mapeamento) == 0 {
		return
	}
	for _, candidato := range dados.Candidatos {
		if novo, ok := mapeamento[candidato.EleicaoID]; ok {
			candidato.EleicaoID = novo
		}
	}
	for _, prestacao := range dados.Prestacoes {
		if novo, ok := mapeamento[prestacao.EleicaoID]; ok {
			prestacao.EleicaoID = novo
		}
	}
}

func remapearUnidadeEleitoralIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	if len(mapeamento) == 0 {
		return
	}
	for _, prestacao := range dados.Prestacoes {
		if prestacao.UnidadeEleitoralID != nil {
			if novo, ok := mapeamento[*prestacao.UnidadeEleitoralID]; ok {
				prestacao.UnidadeEleitoralID = ponteiroUUID(novo)
			}
		}
	}
}

func remapearPartidoIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	if len(mapeamento) == 0 {
		return
	}
	for _, candidato := range dados.Candidatos {
		if candidato.PartidoID != nil {
			if novo, ok := mapeamento[*candidato.PartidoID]; ok {
				candidato.PartidoID = ponteiroUUID(novo)
			}
		}
	}
	for _, prestacao := range dados.Prestacoes {
		if prestacao.PartidoID != nil {
			if novo, ok := mapeamento[*prestacao.PartidoID]; ok {
				prestacao.PartidoID = ponteiroUUID(novo)
			}
		}
	}
	for _, despesa := range dados.DespesasOrgaoPartidario {
		if novo, ok := mapeamento[despesa.PartidoID]; ok {
			despesa.PartidoID = novo
		}
	}
	for _, receita := range dados.ReceitasOrgaoPartidario {
		if novo, ok := mapeamento[receita.PartidoID]; ok {
			receita.PartidoID = novo
		}
	}
}

func remapearCandidatoIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	if len(mapeamento) == 0 {
		return
	}
	for _, bem := range dados.BensCandidato {
		if novo, ok := mapeamento[bem.CandidatoID]; ok {
			bem.CandidatoID = novo
		}
	}
	for _, prestacao := range dados.Prestacoes {
		if prestacao.CandidatoID != nil {
			if novo, ok := mapeamento[*prestacao.CandidatoID]; ok {
				prestacao.CandidatoID = ponteiroUUID(novo)
			}
		}
	}
	for _, despesa := range dados.DespesasCandidato {
		if novo, ok := mapeamento[despesa.CandidatoID]; ok {
			despesa.CandidatoID = novo
		}
	}
	for _, receita := range dados.ReceitasCandidato {
		if novo, ok := mapeamento[receita.CandidatoID]; ok {
			receita.CandidatoID = novo
		}
	}
}

func remapearFornecedorIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	log := logger.New("LeitorCSV: Utils: remapearFornecedorIDs")
	if len(mapeamento) == 0 {
		return
	}
	for _, despesa := range dados.DespesasCandidato {
		if despesa.FornecedorID != nil {
			if novo, ok := mapeamento[*despesa.FornecedorID]; ok {
				despesa.FornecedorID = ponteiroUUID(novo)
			} else {
				log.Error("fornecedor nao mapeado - definindo fornecedor_id NULL",
					"fornecedor_id", despesa.FornecedorID.String(), "sq_despesa", despesa.SQDespesa, "tipo_registro", despesa.TipoRegistro, "despesa_id", despesa.ID,
					"tipo_documento", despesa.TipoDocumento, "numero_documento", despesa.NumeroDocumento, "valor", despesa.Valor)
				despesa.FornecedorID = nil
			}
		}
	}
	for _, despesa := range dados.DespesasOrgaoPartidario {
		if despesa.FornecedorID != nil {
			if novo, ok := mapeamento[*despesa.FornecedorID]; ok {
				despesa.FornecedorID = ponteiroUUID(novo)
			} else {
				log.Error("fornecedor nao mapeado - definindo fornecedor_id NULL",
					"fornecedor_id", despesa.FornecedorID.String(), "sq_despesa", despesa.SQDespesa, "tipo_registro", despesa.TipoRegistro, "despesa_orgao_partidario_id", despesa.ID,
					"tipo_documento", despesa.TipoDocumento, "numero_documento", despesa.NumeroDocumento, "valor", despesa.Valor)
				despesa.FornecedorID = nil
			}
		}
	}
}

func remapearDoadorIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	log := logger.New("LeitorCSV: Utils: remapearDoadorIDs")
	if len(mapeamento) == 0 {
		return
	}
	for _, receita := range dados.ReceitasCandidato {
		if receita.DoadorID != nil {
			if novo, ok := mapeamento[*receita.DoadorID]; ok {
				receita.DoadorID = ponteiroUUID(novo)
			} else {
				log.Error("doador nao mapeado - definindo doador_id NULL",
					"doador_id", receita.DoadorID.String(), "sq_receita", receita.SQReceita, "receita_candidato_id", receita.ID,
					"descricao", receita.Descricao, "valor", receita.Valor)
				receita.DoadorID = nil
			}
		}
	}
	for _, receita := range dados.ReceitasOrgaoPartidario {
		if receita.DoadorID != nil {
			if novo, ok := mapeamento[*receita.DoadorID]; ok {
				receita.DoadorID = ponteiroUUID(novo)
			} else {
				log.Error("doador nao mapeado - definindo doador_id NULL",
					"doador_id", receita.DoadorID.String(), "sq_receita", receita.SQReceita, "receita_orgao_partidario_id", receita.ID,
					"descricao", receita.Descricao, "valor", receita.Valor)
				receita.DoadorID = nil
			}
		}
	}
}

func remapearReceitaCandidatoIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	if len(mapeamento) == 0 {
		return
	}
	for _, origem := range dados.ReceitasDoadorOriginarioCandidato {
		if origem.ReceitaCandidatoID != nil {
			if novo, ok := mapeamento[*origem.ReceitaCandidatoID]; ok {
				origem.ReceitaCandidatoID = ponteiroUUID(novo)
			}
		}
	}
}

func remapearReceitaOrgaoPartidarioIDs(dados *tipos.DadosImportacao, mapeamento map[uuid.UUID]uuid.UUID) {
	if len(mapeamento) == 0 {
		return
	}
	for _, origem := range dados.ReceitasDoadorOriginarioOrgaoPartidario {
		if origem.ReceitaOrgaoPartidarioID != nil {
			if novo, ok := mapeamento[*origem.ReceitaOrgaoPartidarioID]; ok {
				origem.ReceitaOrgaoPartidarioID = ponteiroUUID(novo)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Sincronizacao: validacao de divergencia entre IDs nas transacoes e
// nas prestacoes. Usada como diagnostico — nao altera IDs.
// ---------------------------------------------------------------------------
func sincronizarDependenciasDePrestacao(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	log := logger.New("LeitorCSV: Utils: sincronizarDependenciasDePrestacao")
	idToPrestacao := make(map[uuid.UUID]*types.PrestacaoContas, len(dados.Prestacoes))
	for _, p := range dados.Prestacoes {
		idToPrestacao[p.ID] = p
	}
	for _, despesa := range dados.DespesasCandidato {
		if prestacao, ok := idToPrestacao[despesa.PrestacaoContasID]; ok {
			despesa.PrestacaoContasID = prestacao.ID
			if prestacao.CandidatoID != nil {
				if despesa.CandidatoID != *prestacao.CandidatoID {
					log.Warn("candidato_id divergente — mantendo valor original do CSV",
						"sq_despesa", despesa.SQDespesa,
						"despesa_candidato_id", despesa.CandidatoID,
						"prestacao_candidato_id", *prestacao.CandidatoID)
				}
			}
		}
	}
	for _, despesa := range dados.DespesasOrgaoPartidario {
		if prestacao, ok := idToPrestacao[despesa.PrestacaoContasID]; ok {
			despesa.PrestacaoContasID = prestacao.ID
			if prestacao.PartidoID != nil {
				if despesa.PartidoID != *prestacao.PartidoID {
					log.Warn("partido_id divergente — mantendo valor original do CSV",
						"sq_despesa", despesa.SQDespesa,
						"despesa_partido_id", despesa.PartidoID,
						"prestacao_partido_id", *prestacao.PartidoID)
				}
			}
		}
	}
	for _, receita := range dados.ReceitasCandidato {
		if prestacao, ok := idToPrestacao[receita.PrestacaoContasID]; ok {
			receita.PrestacaoContasID = prestacao.ID
			if prestacao.CandidatoID != nil {
				if receita.CandidatoID != *prestacao.CandidatoID {
					log.Warn("candidato_id divergente — mantendo valor original do CSV",
						"sq_receita", receita.SQReceita,
						"receita_candidato_id", receita.CandidatoID,
						"prestacao_candidato_id", *prestacao.CandidatoID)
				}
			}
		}
	}
	for _, receita := range dados.ReceitasOrgaoPartidario {
		if prestacao, ok := idToPrestacao[receita.PrestacaoContasID]; ok {
			receita.PrestacaoContasID = prestacao.ID
			if prestacao.PartidoID != nil {
				if receita.PartidoID != *prestacao.PartidoID {
					log.Warn("partido_id divergente — mantendo valor original do CSV",
						"sq_receita", receita.SQReceita,
						"receita_partido_id", receita.PartidoID,
						"prestacao_partido_id", *prestacao.PartidoID)
				}
			}
		}
	}
	for _, origem := range dados.ReceitasDoadorOriginarioCandidato {
		if prestacao, ok := idToPrestacao[origem.PrestacaoContasID]; ok {
			origem.PrestacaoContasID = prestacao.ID
		}
	}
	for _, origem := range dados.ReceitasDoadorOriginarioOrgaoPartidario {
		if prestacao, ok := idToPrestacao[origem.PrestacaoContasID]; ok {
			origem.PrestacaoContasID = prestacao.ID
		}
	}
}

func obterCandidatoPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.Candidato {
	return dados.CandidatosPorID[id]
}

func obterPrestacaoPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.PrestacaoContas {
	return dados.PrestacoesPorID[id]
}

// valores extrai os valores de um mapa generico para um slice (util para iterar mapas com chave descartavel)
func valores[K comparable, V any](m map[K]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}
