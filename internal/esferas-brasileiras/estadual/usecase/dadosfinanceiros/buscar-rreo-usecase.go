package dadosfinanceiros

import (
	"context"
	"sort"
	"strings"

	siconfiClient "github.com/danyele/podp/internal/shared/clients/siconfi"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarRREORequest struct {
	UF        string
	Exercicio int64
}

type EsferaEstadualBuscarRREOResponse struct {
	Gastos   []types.GastoPorFuncao
	Receitas []types.ReceitaResumo
}

type EsferaEstadualBuscarRREOUseCase struct {
	*BaseFinanceiroUseCase
}

func NovoEsferaEstadualBuscarRREOUseCase(base *BaseFinanceiroUseCase) *EsferaEstadualBuscarRREOUseCase {
	return &EsferaEstadualBuscarRREOUseCase{BaseFinanceiroUseCase: base}
}

func (u *EsferaEstadualBuscarRREOUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarRREORequest) (*EsferaEstadualBuscarRREOResponse, error) {
	gastos, receitas := u.buscarRREO(ctx, req.UF, req.Exercicio)
	return &EsferaEstadualBuscarRREOResponse{Gastos: gastos, Receitas: receitas}, nil
}

func (u *EsferaEstadualBuscarRREOUseCase) buscarRREO(ctx context.Context, uf string, exercicio int64) ([]types.GastoPorFuncao, []types.ReceitaResumo) {
	idEnte, err := u.estadoID(ctx, uf)
	if err != nil || idEnte == 0 {
		return nil, nil
	}

	if exercicio <= 0 {
		exercicio = u.anoAlvo()
	}
	var gastos []types.GastoPorFuncao
	var receitas []types.ReceitaResumo

	periods := []int{6, 5}
	for _, periodo := range periods {
		if len(gastos) == 0 {
			gastos = u.buscarGastosPorFuncao(ctx, idEnte, exercicio, periodo)
		}
		if len(receitas) == 0 {
			receitas = u.buscarReceitas(ctx, idEnte, exercicio, periodo)
		}
		if len(gastos) > 0 && len(receitas) > 0 {
			break
		}
	}

	return gastos, receitas
}

func (u *EsferaEstadualBuscarRREOUseCase) buscarGastosPorFuncao(ctx context.Context, idEnte int, exercicio int64, periodo int) []types.GastoPorFuncao {
	params := siconfiClient.RREOParams{
		AnExercicio:         exercicio,
		NrPeriodo:           periodo,
		CoTipoDemonstrativo: "RREO",
		IdEnte:              idEnte,
		NoAnexo:             "RREO-Anexo 02",
		CoEsfera:            "E",
	}

	itens, err := u.BaseFinanceiroUseCase.buscarRREO(ctx, params)
	if err != nil || len(itens) == 0 {
		return nil
	}

	despesasPorFuncao := make(map[string]*types.GastoPorFuncao)
	for _, item := range itens {
		funcao := item.Conta
		if funcao == "" {
			continue
		}
		if _, ok := despesasPorFuncao[funcao]; !ok {
			despesasPorFuncao[funcao] = &types.GastoPorFuncao{Funcao: funcao}
		}
		colUpper := strings.ToUpper(item.Coluna)
		switch {
		case strings.Contains(colUpper, "EMPENHAD"):
			despesasPorFuncao[funcao].Empenhado += item.Valor
		case strings.Contains(colUpper, "LIQUIDAD"):
			despesasPorFuncao[funcao].Liquidado += item.Valor
		case strings.Contains(colUpper, "PAG"):
			despesasPorFuncao[funcao].Pago += item.Valor
		}
	}

	var gastos []types.GastoPorFuncao
	for _, g := range despesasPorFuncao {
		gastos = append(gastos, *g)
	}
	sort.Slice(gastos, func(i, j int) bool {
		return gastos[i].Empenhado > gastos[j].Empenhado
	})
	return gastos
}

func (u *EsferaEstadualBuscarRREOUseCase) buscarReceitas(ctx context.Context, idEnte int, exercicio int64, periodo int) []types.ReceitaResumo {
	params := siconfiClient.RREOParams{
		AnExercicio:         exercicio,
		NrPeriodo:           periodo,
		CoTipoDemonstrativo: "RREO",
		IdEnte:              idEnte,
		NoAnexo:             "RREO-Anexo 03",
		CoEsfera:            "E",
	}

	itens, err := u.BaseFinanceiroUseCase.buscarRREO(ctx, params)
	if err != nil || len(itens) == 0 {
		return nil
	}

	var receitas []types.ReceitaResumo
	for _, item := range itens {
		receitas = append(receitas, types.ReceitaResumo{
			Conta:     item.Conta,
			Coluna:    item.Coluna,
			Valor:     item.Valor,
			Exercicio: exercicio,
		})
	}
	sort.Slice(receitas, func(i, j int) bool {
		return receitas[i].Valor > receitas[j].Valor
	})
	return receitas
}
