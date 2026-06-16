package parse

import (
	"context"

	"github.com/google/uuid"

	tipos "github.com/danyele/laceu/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/laceu/internal/shared/types"
)

type ProcessadorLeitorCSV struct {
	TamanhoLote          int
	dados                *tipos.DadosImportacao
	cacheCandidatos      map[int64]*types.Candidato
	buscarCandidatoPorSQ func(context.Context, int64) (uuid.UUID, error)
}

func NovoProcessadorLeitorCSV(tamanhoLote int) *ProcessadorLeitorCSV {
	if tamanhoLote <= 0 {
		tamanhoLote = tipos.TamanhoLotePadraoImportacao
	}
	return &ProcessadorLeitorCSV{
		TamanhoLote: tamanhoLote,
		dados:       tipos.NovoDadosImportacao(),
	}
}

func NovoProcessadorComDados(tamanhoLote int, dados *tipos.DadosImportacao) *ProcessadorLeitorCSV {
	if tamanhoLote <= 0 {
		tamanhoLote = tipos.TamanhoLotePadraoImportacao
	}
	if dados == nil {
		dados = tipos.NovoDadosImportacao()
	}
	return &ProcessadorLeitorCSV{
		TamanhoLote: tamanhoLote,
		dados:       dados,
	}
}

func (p *ProcessadorLeitorCSV) Dados() *tipos.DadosImportacao {
	return p.dados
}

func (p *ProcessadorLeitorCSV) ComResolverCandidato(fn func(context.Context, int64) (uuid.UUID, error)) {
	p.buscarCandidatoPorSQ = fn
}

func (p *ProcessadorLeitorCSV) ComCacheCandidatos(cache map[int64]*types.Candidato) {
	p.cacheCandidatos = cache
}

func NovoProcessadorComCacheCandidatos(tamanhoLote int, cache map[int64]*types.Candidato) *ProcessadorLeitorCSV {
	p := NovoProcessadorLeitorCSV(tamanhoLote)
	p.ComCacheCandidatos(cache)
	return p
}
