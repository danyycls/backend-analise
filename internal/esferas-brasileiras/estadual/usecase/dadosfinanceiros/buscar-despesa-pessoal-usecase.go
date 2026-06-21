package dadosfinanceiros

import (
	"context"
	"strconv"
	"strings"

	siconfiClient "github.com/danyele/podp/internal/shared/clients/siconfi"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarDespesaPessoalRequest struct {
	UF string
}

type EsferaEstadualBuscarDespesaPessoalResponse struct {
	Dados *types.DespesaPessoalResumo
}

type EsferaEstadualBuscarDespesaPessoalUseCase struct {
	*BaseFinanceiroUseCase
}

func NovoEsferaEstadualBuscarDespesaPessoalUseCase(base *BaseFinanceiroUseCase) *EsferaEstadualBuscarDespesaPessoalUseCase {
	return &EsferaEstadualBuscarDespesaPessoalUseCase{BaseFinanceiroUseCase: base}
}

func (u *EsferaEstadualBuscarDespesaPessoalUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarDespesaPessoalRequest) (*EsferaEstadualBuscarDespesaPessoalResponse, error) {
	resultado := u.buscarDespesaPessoal(ctx, req.UF)
	return &EsferaEstadualBuscarDespesaPessoalResponse{Dados: resultado}, nil
}

func (u *EsferaEstadualBuscarDespesaPessoalUseCase) buscarDespesaPessoal(ctx context.Context, uf string) *types.DespesaPessoalResumo {
	log := logger.New("Estadual: UseCase: BuscarDespesaPessoal")
	idEnte, err := u.estadoID(ctx, uf)
	if err != nil || idEnte == 0 {
		log.Error("erro ao obter id ente estado", "uf", uf, "erro", err)
		return nil
	}

	for _, t := range u.montarTentativas() {
		params := siconfiClient.RGFParams{
			AnExercicio:         t.ano,
			InPeriodicidade:     t.periodicidade,
			NrPeriodo:           t.periodo,
			CoTipoDemonstrativo: "RGF",
			CoPoder:             "E",
			IdEnte:              idEnte,
			CoEsfera:            "E",
			NoAnexo:             "RGF-Anexo 01",
		}

		itens, err := u.buscarRGF(ctx, params)
		if err != nil || len(itens) == 0 {
			continue
		}

		resultado := u.extrairDespesaPessoal(itens, t.ano)
		if resultado != nil {
			return resultado
		}
	}

	return nil
}

type tentativaRGF struct {
	ano           int64
	periodicidade string
	periodo       int
}

func (u *EsferaEstadualBuscarDespesaPessoalUseCase) montarTentativas() []tentativaRGF {
	alvo := u.anoAlvo()
	return []tentativaRGF{
		{alvo, "Q", 3},
		{alvo, "S", 2},
		{alvo - 1, "Q", 3},
		{alvo - 1, "S", 2},
		{alvo - 2, "Q", 3},
		{alvo - 2, "S", 2},
	}
}

func (u *EsferaEstadualBuscarDespesaPessoalUseCase) extrairDespesaPessoal(itens []siconfiClient.RGFItem, ano int64) *types.DespesaPessoalResumo {
	var totalDespesa float64
	var percentualRCL float64
	for _, item := range itens {
		colUpper := strings.ToUpper(item.Coluna)
		if strings.Contains(colUpper, "RCL") && strings.Contains(colUpper, "%") {
			percentualRCL = item.Valor
		}
		if strings.Contains(colUpper, "DESPESA") && strings.Contains(colUpper, "PESSOAL") {
			totalDespesa += item.Valor
		}
	}

	if totalDespesa > 0 || percentualRCL > 0 {
		return &types.DespesaPessoalResumo{
			ValorTotal:    totalDespesa,
			PercentualRCL: percentualRCL,
			Poder:         "Executivo",
			Periodo:       strconv.FormatInt(ano, 10),
		}
	}
	return nil
}
