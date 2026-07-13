package types

import (
	"github.com/danyele/podp/internal/shared/types"
	"github.com/google/uuid"
)

type DadosImportacao struct {
	Eleicoes                                map[int]*types.Eleicao
	UnidadesEleitorais                      map[string]*types.UnidadeEleitoral
	Partidos                                map[int16]*types.Partido
	Candidatos                              map[int64]*types.Candidato
	CandidatosPorID                         map[uuid.UUID]*types.Candidato
	Fornecedores                            map[string]*types.Fornecedor
	Doadores                                map[string]*types.Doador
	Prestacoes                              map[string]*types.PrestacaoContas
	PrestacoesPorTipoESQ                    map[string]*types.PrestacaoContas
	PrestacoesPorID                         map[uuid.UUID]*types.PrestacaoContas
	ReceitasCandidatoPorSQ                  map[int64]*types.ReceitaCandidato
	ReceitasOrgaoPorSQ                      map[int64]*types.ReceitaOrgaoPartidario
	DespesasCandidato                       []*types.DespesaCandidato
	DespesasOrgaoPartidario                 []*types.DespesaOrgaoPartidario
	ReceitasCandidato                       []*types.ReceitaCandidato
	ReceitasOrgaoPartidario                 []*types.ReceitaOrgaoPartidario
	ReceitasDoadorOriginarioCandidato       []*types.ReceitaDoadorOriginarioCandidato
	ReceitasDoadorOriginarioOrgaoPartidario []*types.ReceitaDoadorOriginarioOrgaoPartidario
	BensCandidato                           map[string]*types.BemCandidato
	Convenios                               []*types.Convenio
}

func NovoDadosImportacao() *DadosImportacao {
	return &DadosImportacao{
		Eleicoes:               make(map[int]*types.Eleicao),
		UnidadesEleitorais:     make(map[string]*types.UnidadeEleitoral),
		Partidos:               make(map[int16]*types.Partido),
		Candidatos:             make(map[int64]*types.Candidato),
		CandidatosPorID:        make(map[uuid.UUID]*types.Candidato),
		Fornecedores:           make(map[string]*types.Fornecedor),
		Doadores:               make(map[string]*types.Doador),
		Prestacoes:             make(map[string]*types.PrestacaoContas),
		PrestacoesPorTipoESQ:   make(map[string]*types.PrestacaoContas),
		PrestacoesPorID:        make(map[uuid.UUID]*types.PrestacaoContas),
		ReceitasCandidatoPorSQ: make(map[int64]*types.ReceitaCandidato),
		ReceitasOrgaoPorSQ:     make(map[int64]*types.ReceitaOrgaoPartidario),
		BensCandidato:          make(map[string]*types.BemCandidato),
	}
}
