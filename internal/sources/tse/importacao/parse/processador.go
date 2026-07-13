package parse

import (
	"context"

	"github.com/danyele/podp/internal/shared/logger"
	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
	tipos "github.com/danyele/podp/internal/sources/tse/importacao/types"
)

type ProcessadorLeitorCSV struct {
	TamanhoLote          int
	log                  *logger.Logger
	ultimoHash           string
	dados                *tipos.DadosImportacao
	cacheCandidatos      map[int64]*types.Candidato
	buscarCandidatoPorSQ func(context.Context, int64) (uuid.UUID, error)
	RegistrosIgnorados   int
}

func NovoProcessadorLeitorCSV(tamanhoLote int) *ProcessadorLeitorCSV {
	if tamanhoLote <= 0 {
		tamanhoLote = tipos.TamanhoLotePadraoImportacao
	}
	return &ProcessadorLeitorCSV{
		TamanhoLote: tamanhoLote,
		log:         logger.New("LeitorCSV: Processador"),
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
		log:         logger.New("LeitorCSV: Processador"),
		dados:       dados,
	}
}

func (p *ProcessadorLeitorCSV) Dados() *tipos.DadosImportacao {
	return p.dados
}

func (p *ProcessadorLeitorCSV) UltimoHash() string {
	return p.ultimoHash
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

func (p *ProcessadorLeitorCSV) processarCSV(caminho string, fn func(numeroLinha int, registro map[string]string) error) (int, error) {
	total := 0
	hash, err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		if err := fn(numeroLinha, registro); err != nil {
			return err
		}
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}
	p.ultimoHash = hash
	return total, nil
}
