package parse

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
)

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
				log.Info("fornecedor nao mapeado - definindo fornecedor_id NULL",
					"fornecedor_id", despesa.FornecedorID.String(), "despesa_id", despesa.ID)
				despesa.FornecedorID = nil
			}
		}
	}
	for _, despesa := range dados.DespesasOrgaoPartidario {
		if despesa.FornecedorID != nil {
			if novo, ok := mapeamento[*despesa.FornecedorID]; ok {
				despesa.FornecedorID = ponteiroUUID(novo)
			} else {
				log.Info("fornecedor nao mapeado - definindo fornecedor_id NULL",
					"fornecedor_id", despesa.FornecedorID.String(), "despesa_orgao_partidario_id", despesa.ID)
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
				log.Info("doador nao mapeado - definindo doador_id NULL",
					"doador_id", receita.DoadorID.String(), "receita_candidato_id", receita.ID)
				receita.DoadorID = nil
			}
		}
	}
	for _, receita := range dados.ReceitasOrgaoPartidario {
		if receita.DoadorID != nil {
			if novo, ok := mapeamento[*receita.DoadorID]; ok {
				receita.DoadorID = ponteiroUUID(novo)
			} else {
				log.Info("doador nao mapeado - definindo doador_id NULL",
					"doador_id", receita.DoadorID.String(), "receita_orgao_partidario_id", receita.ID)
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
// Sincronizacao de dependencies: apos reconciliar, atualiza IDs em memoria
// ---------------------------------------------------------------------------

func sincronizarDependenciasDeEleicao(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToEleicao := make(map[uuid.UUID]*types.Eleicao, len(dados.Eleicoes))
	for _, e := range dados.Eleicoes {
		idToEleicao[e.ID] = e
	}
	for _, candidato := range dados.Candidatos {
		if eleicao, ok := idToEleicao[candidato.EleicaoID]; ok {
			candidato.EleicaoID = eleicao.ID
		}
	}
	for _, prestacao := range dados.Prestacoes {
		if eleicao, ok := idToEleicao[prestacao.EleicaoID]; ok {
			prestacao.EleicaoID = eleicao.ID
		}
	}
}

func sincronizarDependenciasDeUnidadeEleitoral(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToUnidade := make(map[uuid.UUID]*types.UnidadeEleitoral, len(dados.UnidadesEleitorais))
	for _, u := range dados.UnidadesEleitorais {
		idToUnidade[u.ID] = u
	}
	for _, prestacao := range dados.Prestacoes {
		if prestacao.UnidadeEleitoralID != nil {
			if unidade, ok := idToUnidade[*prestacao.UnidadeEleitoralID]; ok {
				prestacao.UnidadeEleitoralID = ponteiroUUID(unidade.ID)
			}
		}
	}
}

func sincronizarDependenciasDePartido(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToPartido := make(map[uuid.UUID]*types.Partido, len(dados.Partidos))
	for _, p := range dados.Partidos {
		idToPartido[p.ID] = p
	}
	for _, candidato := range dados.Candidatos {
		if candidato.PartidoID != nil {
			if partido, ok := idToPartido[*candidato.PartidoID]; ok {
				candidato.PartidoID = ponteiroUUID(partido.ID)
			}
		}
	}
	for _, prestacao := range dados.Prestacoes {
		if prestacao.PartidoID != nil {
			if partido, ok := idToPartido[*prestacao.PartidoID]; ok {
				prestacao.PartidoID = ponteiroUUID(partido.ID)
			}
		}
	}
	for _, despesa := range dados.DespesasOrgaoPartidario {
		if partido, ok := idToPartido[despesa.PartidoID]; ok {
			despesa.PartidoID = partido.ID
		}
	}
	for _, receita := range dados.ReceitasOrgaoPartidario {
		if partido, ok := idToPartido[receita.PartidoID]; ok {
			receita.PartidoID = partido.ID
		}
	}
}

func sincronizarDependenciasDeCandidato(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToCandidato := make(map[uuid.UUID]*types.Candidato, len(dados.Candidatos))
	for _, c := range dados.Candidatos {
		idToCandidato[c.ID] = c
	}
	for _, bem := range dados.BensCandidato {
		if candidato, ok := idToCandidato[bem.CandidatoID]; ok {
			bem.CandidatoID = candidato.ID
		}
	}
	for _, prestacao := range dados.Prestacoes {
		if prestacao.CandidatoID != nil {
			if candidato, ok := idToCandidato[*prestacao.CandidatoID]; ok {
				prestacao.CandidatoID = ponteiroUUID(candidato.ID)
			}
		}
	}
	for _, receita := range dados.ReceitasCandidato {
		if candidato, ok := idToCandidato[receita.CandidatoID]; ok {
			receita.CandidatoID = candidato.ID
		}
	}
}

func sincronizarDependenciasDeFornecedor(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToFornecedor := make(map[uuid.UUID]*types.Fornecedor, len(dados.Fornecedores))
	for _, f := range dados.Fornecedores {
		idToFornecedor[f.ID] = f
	}
	for _, despesa := range dados.DespesasCandidato {
		if despesa.FornecedorID != nil {
			if fornecedor, ok := idToFornecedor[*despesa.FornecedorID]; ok {
				despesa.FornecedorID = ponteiroUUID(fornecedor.ID)
			}
		}
	}
	for _, despesa := range dados.DespesasOrgaoPartidario {
		if despesa.FornecedorID != nil {
			if fornecedor, ok := idToFornecedor[*despesa.FornecedorID]; ok {
				despesa.FornecedorID = ponteiroUUID(fornecedor.ID)
			}
		}
	}
}

func sincronizarDependenciasDeDoador(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToDoador := make(map[uuid.UUID]*types.Doador, len(dados.Doadores))
	for _, d := range dados.Doadores {
		idToDoador[d.ID] = d
	}
	for _, receita := range dados.ReceitasCandidato {
		if receita.DoadorID != nil {
			if doador, ok := idToDoador[*receita.DoadorID]; ok {
				receita.DoadorID = ponteiroUUID(doador.ID)
			}
		}
	}
	for _, receita := range dados.ReceitasOrgaoPartidario {
		if receita.DoadorID != nil {
			if doador, ok := idToDoador[*receita.DoadorID]; ok {
				receita.DoadorID = ponteiroUUID(doador.ID)
			}
		}
	}
}

func sincronizarDependenciasDePrestacao(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToPrestacao := make(map[uuid.UUID]*types.PrestacaoContas, len(dados.Prestacoes))
	for _, p := range dados.Prestacoes {
		idToPrestacao[p.ID] = p
	}
	for _, despesa := range dados.DespesasCandidato {
		if prestacao, ok := idToPrestacao[despesa.PrestacaoContasID]; ok {
			despesa.PrestacaoContasID = prestacao.ID
			if prestacao.CandidatoID != nil {
				despesa.CandidatoID = *prestacao.CandidatoID
			}
		}
	}
	for _, despesa := range dados.DespesasOrgaoPartidario {
		if prestacao, ok := idToPrestacao[despesa.PrestacaoContasID]; ok {
			despesa.PrestacaoContasID = prestacao.ID
			if prestacao.PartidoID != nil {
				despesa.PartidoID = *prestacao.PartidoID
			}
		}
	}
	for _, receita := range dados.ReceitasCandidato {
		if prestacao, ok := idToPrestacao[receita.PrestacaoContasID]; ok {
			receita.PrestacaoContasID = prestacao.ID
			if prestacao.CandidatoID != nil {
				receita.CandidatoID = *prestacao.CandidatoID
			}
		}
	}
	for _, receita := range dados.ReceitasOrgaoPartidario {
		if prestacao, ok := idToPrestacao[receita.PrestacaoContasID]; ok {
			receita.PrestacaoContasID = prestacao.ID
			if prestacao.PartidoID != nil {
				receita.PartidoID = *prestacao.PartidoID
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

func sincronizarDependenciasDeReceitaCandidato(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToReceita := make(map[uuid.UUID]*types.ReceitaCandidato, len(dados.ReceitasCandidatoPorSQ))
	for _, r := range dados.ReceitasCandidatoPorSQ {
		idToReceita[r.ID] = r
	}
	for _, origem := range dados.ReceitasDoadorOriginarioCandidato {
		if origem.ReceitaCandidatoID != nil {
			if receita, ok := idToReceita[*origem.ReceitaCandidatoID]; ok {
				origem.ReceitaCandidatoID = ponteiroUUID(receita.ID)
			}
		}
	}
}

func sincronizarDependenciasDeReceitaOrgaoPartidario(dados *tipos.DadosImportacao) {
	if dados == nil {
		return
	}
	idToReceita := make(map[uuid.UUID]*types.ReceitaOrgaoPartidario, len(dados.ReceitasOrgaoPorSQ))
	for _, r := range dados.ReceitasOrgaoPorSQ {
		idToReceita[r.ID] = r
	}
	for _, origem := range dados.ReceitasDoadorOriginarioOrgaoPartidario {
		if origem.ReceitaOrgaoPartidarioID != nil {
			if receita, ok := idToReceita[*origem.ReceitaOrgaoPartidarioID]; ok {
				origem.ReceitaOrgaoPartidarioID = ponteiroUUID(receita.ID)
			}
		}
	}
}

func obterEleicaoPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.Eleicao { //nolint:unused
	for _, item := range dados.Eleicoes {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterUnidadePorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.UnidadeEleitoral { //nolint:unused
	for _, item := range dados.UnidadesEleitorais {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterPartidoPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.Partido { //nolint:unused
	for _, item := range dados.Partidos {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterCandidatoPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.Candidato {
	for _, item := range dados.Candidatos {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterFornecedorPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.Fornecedor { //nolint:unused
	for _, item := range dados.Fornecedores {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterDoadorPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.Doador { //nolint:unused
	for _, item := range dados.Doadores {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterPrestacaoPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.PrestacaoContas {
	for _, item := range dados.Prestacoes {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterReceitaCandidatoPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.ReceitaCandidato { //nolint:unused
	for _, item := range dados.ReceitasCandidatoPorSQ {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func obterReceitaOrgaoPartidarioPorID(dados *tipos.DadosImportacao, id uuid.UUID) *types.ReceitaOrgaoPartidario { //nolint:unused
	for _, item := range dados.ReceitasOrgaoPorSQ {
		if item.ID == id {
			return item
		}
	}
	return nil
}

// valores extrai os valores de um mapa generico para um slice (util para iterar mapas com chave descartavel)
func valores[K comparable, V any](m map[K]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}
