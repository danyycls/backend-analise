package types

import (
	"github.com/danyele/laceu/internal/shared/types"
)

type DadosImportacao struct {
	Eleicoes                                map[int]*types.Eleicao
	UnidadesEleitorais                      map[string]*types.UnidadeEleitoral
	Partidos                                map[int16]*types.Partido
	Candidatos                              map[int64]*types.Candidato
	Fornecedores                            map[string]*types.Fornecedor
	Doadores                                map[string]*types.Doador
	Prestacoes                              map[string]*types.PrestacaoContas
	PrestacoesPorTipoESQ                    map[string]*types.PrestacaoContas
	ReceitasCandidatoPorSQ                  map[int64]*types.ReceitaCandidato
	ReceitasOrgaoPorSQ                      map[int64]*types.ReceitaOrgaoPartidario
	DespesasCandidato                       []*types.DespesaCandidato
	DespesasOrgaoPartidario                 []*types.DespesaOrgaoPartidario
	ReceitasCandidato                       []*types.ReceitaCandidato
	ReceitasOrgaoPartidario                 []*types.ReceitaOrgaoPartidario
	ReceitasDoadorOriginarioCandidato       []*types.ReceitaDoadorOriginarioCandidato
	ReceitasDoadorOriginarioOrgaoPartidario []*types.ReceitaDoadorOriginarioOrgaoPartidario
	BensCandidato                           []*types.BemCandidato
}

func NovoDadosImportacao() *DadosImportacao {
	return &DadosImportacao{
		Eleicoes:               make(map[int]*types.Eleicao),
		UnidadesEleitorais:     make(map[string]*types.UnidadeEleitoral),
		Partidos:               make(map[int16]*types.Partido),
		Candidatos:             make(map[int64]*types.Candidato),
		Fornecedores:           make(map[string]*types.Fornecedor),
		Doadores:               make(map[string]*types.Doador),
		Prestacoes:             make(map[string]*types.PrestacaoContas),
		PrestacoesPorTipoESQ:   make(map[string]*types.PrestacaoContas),
		ReceitasCandidatoPorSQ: make(map[int64]*types.ReceitaCandidato),
		ReceitasOrgaoPorSQ:     make(map[int64]*types.ReceitaOrgaoPartidario),
	}
}
