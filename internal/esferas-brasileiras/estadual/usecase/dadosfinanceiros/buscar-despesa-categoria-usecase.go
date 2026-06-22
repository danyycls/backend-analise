package dadosfinanceiros

import (
	"context"
	"sort"
	"strings"

	siconfiClient "github.com/danyele/podp/internal/shared/clients/siconfi"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarDespesaCategoriaRequest struct {
	UF        string
	Exercicio int64
}

type EsferaEstadualBuscarDespesaCategoriaResponse struct {
	Dados []types.ServidorMunicipio
}

type EsferaEstadualBuscarDespesaCategoriaUseCase struct {
	*BaseFinanceiroUseCase
}

func NovoEsferaEstadualBuscarDespesaCategoriaUseCase(base *BaseFinanceiroUseCase) *EsferaEstadualBuscarDespesaCategoriaUseCase {
	return &EsferaEstadualBuscarDespesaCategoriaUseCase{BaseFinanceiroUseCase: base}
}

func (u *EsferaEstadualBuscarDespesaCategoriaUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarDespesaCategoriaRequest) (*EsferaEstadualBuscarDespesaCategoriaResponse, error) {
	resultado := u.buscarDespesaCategoria(ctx, req.UF, req.Exercicio)
	return &EsferaEstadualBuscarDespesaCategoriaResponse{Dados: resultado}, nil
}

func (u *EsferaEstadualBuscarDespesaCategoriaUseCase) buscarDespesaCategoria(ctx context.Context, uf string, exercicio int64) []types.ServidorMunicipio {
	idEnte, err := u.estadoID(ctx, uf)
	if err != nil || idEnte == 0 {
		return nil
	}

	for _, t := range u.montarTentativas(exercicio) {
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

		resultado := u.agruparDespesaPorCategoria(itens)
		if len(resultado) > 0 {
			sort.Slice(resultado, func(i, j int) bool {
				return resultado[i].DespesaTotal > resultado[j].DespesaTotal
			})
			return resultado
		}
	}

	return nil
}

func (u *EsferaEstadualBuscarDespesaCategoriaUseCase) montarTentativas(exercicio int64) []tentativaRGF {
	alvo := exercicio
	if alvo <= 0 {
		alvo = u.anoAlvo()
	}
	return []tentativaRGF{
		{alvo, "Q", 3},
		{alvo, "S", 2},
		{alvo - 1, "Q", 3},
		{alvo - 1, "S", 2},
	}
}

func (u *EsferaEstadualBuscarDespesaCategoriaUseCase) agruparDespesaPorCategoria(itens []siconfiClient.RGFItem) []types.ServidorMunicipio {
	despesasPorCategoria := make(map[string]*types.ServidorMunicipio)
	for _, item := range itens {
		colUpper := strings.ToUpper(item.Coluna)
		if !strings.Contains(colUpper, "DESPESA") || !strings.Contains(colUpper, "PESSOAL") {
			continue
		}
		if strings.Contains(colUpper, "RCL") && strings.Contains(colUpper, "%") {
			continue
		}

		chave := item.Coluna
		for _, palavra := range []string{"ATIVO", "INATIVO", "PENSIONISTA", "TERCEIRIZADO", "NÃO COMPUTADA", "BRUTA", "LÍQUIDA"} {
			if strings.Contains(colUpper, palavra) {
				chave = palavra
				break
			}
		}

		if _, ok := despesasPorCategoria[chave]; !ok {
			despesasPorCategoria[chave] = &types.ServidorMunicipio{
				Categoria: chave,
			}
		}
		despesasPorCategoria[chave].DespesaTotal += item.Valor
	}

	var result []types.ServidorMunicipio
	for _, v := range despesasPorCategoria {
		if v.DespesaTotal > 0 {
			result = append(result, *v)
		}
	}
	return result
}
