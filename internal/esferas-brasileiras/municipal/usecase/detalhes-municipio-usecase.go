package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	pncpClient "github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/clients/siconfi"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/pncpbusca"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
)

type EsferaMunicipalBuscarDetalhesUseCase struct {
	siconfiClient   *siconfi.SICONFIClient
	pncpClient      *pncpClient.PNCPClient
	redis           *redis.RedisCache
	repo            repositorios.PNCPRepository
	apiIndisponivel atomic.Bool
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) SICONFIIndisponivel() bool {
	return u.apiIndisponivel.Load()
}

func NovoEsferaMunicipalBuscarDetalhesUseCase(
	siconfiCli *siconfi.SICONFIClient,
	pncpCli *pncpClient.PNCPClient,
	redis *redis.RedisCache,
	repo repositorios.PNCPRepository,
) *EsferaMunicipalBuscarDetalhesUseCase {
	return &EsferaMunicipalBuscarDetalhesUseCase{
		siconfiClient: siconfiCli,
		pncpClient:    pncpCli,
		redis:         redis,
		repo:          repo,
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarDividaConsolidada(ctx context.Context, idEnte int, exercicio int64) *types.DividaConsolidada {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarDividaConsolidada(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarDisponibilidadeCaixa(ctx context.Context, idEnte int, exercicio int64) *types.DisponibilidadeCaixa {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarDisponibilidadeCaixa(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarRestosAPagar(ctx context.Context, idEnte int, exercicio int64) *types.RestosAPagar {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarRestosAPagar(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarGastoSaude(ctx context.Context, idEnte int, exercicio int64) *types.GastoSaude {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarGastoSaude(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarGastoEducacao(ctx context.Context, idEnte int, exercicio int64) *types.GastoEducacao {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarGastoEducacao(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarFundeb(ctx context.Context, idEnte int, exercicio int64) *types.FundebResumo {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarFundeb(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarBalancoPatrimonial(ctx context.Context, idEnte int, exercicio int64) *types.BalancoPatrimonial {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarBalancoPatrimonial(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarDespesasPorGrupo(ctx context.Context, idEnte int, exercicio int64) []types.DespesaPorGrupoItem {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarDespesasPorGrupo(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarTransferencias(ctx context.Context, idEnte int, exercicio int64) []types.TransferenciaItem {
	if exercicio <= 0 {
		exercicio = int64(time.Now().Year() - 1)
	}
	return u.buscarTransferencias(ctx, idEnte, exercicio)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarContratos(ctx context.Context, codigoIBGE int, ano int) []pncpClient.Contrato {
	if ano <= 0 {
		ano = int(time.Now().Year() - 1)
	}
	return u.buscarContratos(ctx, codigoIBGE, ano)
}

type tentativaRGF struct {
	ano           int64
	periodicidade string
	periodo       int
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) tentativasRGF(exercicio int64) []tentativaRGF {
	alvo := exercicio
	if alvo <= 0 {
		alvo = int64(time.Now().Year() - 1)
	}
	return []tentativaRGF{
		{alvo, "Q", 3},
		{alvo, "S", 2},
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarRGF(ctx context.Context, anexo string, idEnte int, exercicio int64) ([]siconfi.RGFItem, error) {
	var apiIndisponivel bool
	for _, t := range u.tentativasRGF(exercicio) {
		params := siconfi.RGFParams{
			AnExercicio:         t.ano,
			InPeriodicidade:     t.periodicidade,
			NrPeriodo:           t.periodo,
			CoTipoDemonstrativo: "RGF",
			CoPoder:             "E",
			IdEnte:              idEnte,
			NoAnexo:             anexo,
			CoEsfera:            "M",
		}
		itens, err := u.siconfiClient.BuscarRGF(ctx, params)
		if err == nil && len(itens) > 0 {
			return itens, nil
		}
		if errors.Is(err, siconfi.ErrSICONFIIndisponivel) {
			apiIndisponivel = true
		}
	}
	if apiIndisponivel {
		u.apiIndisponivel.Store(true)
		return nil, siconfi.ErrSICONFIIndisponivel
	}
	return nil, nil
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarRREO(ctx context.Context, anexo string, idEnte int, exercicio int64) ([]siconfi.RREOItem, error) {
	alvo := exercicio
	if alvo <= 0 {
		alvo = int64(time.Now().Year() - 1)
	}
	var apiIndisponivel bool
	for _, periodo := range []int{6, 5} {
		params := siconfi.RREOParams{
			AnExercicio:         alvo,
			NrPeriodo:           periodo,
			CoTipoDemonstrativo: "RREO",
			IdEnte:              idEnte,
			NoAnexo:             anexo,
			CoEsfera:            "M",
		}
		itens, err := u.siconfiClient.BuscarRREO(ctx, params)
		if err == nil && len(itens) > 0 {
			return itens, nil
		}
		if errors.Is(err, siconfi.ErrSICONFIIndisponivel) {
			apiIndisponivel = true
		}
	}
	if alvo > 2013 {
		alvo--
		for _, periodo := range []int{6, 5} {
			params := siconfi.RREOParams{
				AnExercicio:         alvo,
				NrPeriodo:           periodo,
				CoTipoDemonstrativo: "RREO",
				IdEnte:              idEnte,
				NoAnexo:             anexo,
				CoEsfera:            "M",
			}
			itens, err := u.siconfiClient.BuscarRREO(ctx, params)
			if err == nil && len(itens) > 0 {
				return itens, nil
			}
			if errors.Is(err, siconfi.ErrSICONFIIndisponivel) {
				apiIndisponivel = true
			}
		}
	}
	if apiIndisponivel {
		u.apiIndisponivel.Store(true)
		return nil, siconfi.ErrSICONFIIndisponivel
	}
	return nil, nil
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarDividaConsolidada(ctx context.Context, idEnte int, exercicio int64) *types.DividaConsolidada {
	log := logger.New("Municipal: UseCase: buscarDividaConsolidada")
	itens, err := u.buscarRGF(ctx, "RGF-Anexo 02", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de divida consolidada encontrado", "ente", idEnte)
		return nil
	}

	var dcl float64
	var pctRCL float64
	var limite float64

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		conta := strings.ToUpper(item.Conta)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(rotulo, "DÍVIDA CONSOLIDADA") || strings.Contains(rotulo, "DCL") ||
			strings.Contains(conta, "DÍVIDA CONSOLIDADA") || strings.Contains(conta, "DCL") {
			if strings.Contains(coluna, "%") || strings.Contains(coluna, "RCL") {
				pctRCL = item.Valor
			} else {
				dcl = item.Valor
			}
		}

		if strings.Contains(rotulo, "LIMITE") || strings.Contains(conta, "LIMITE") {
			if strings.Contains(coluna, "%") {
				limite = item.Valor
			}
		}
	}

	if dcl == 0 && pctRCL == 0 {
		return nil
	}

	return &types.DividaConsolidada{
		ValorDCL:      dcl,
		PercentualRCL: pctRCL,
		LimiteLegal:   limite,
		Periodo:       strconv.FormatInt(exercicio, 10),
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarDisponibilidadeCaixa(ctx context.Context, idEnte int, exercicio int64) *types.DisponibilidadeCaixa {
	log := logger.New("Municipal: UseCase: buscarDisponibilidadeCaixa")
	itens, err := u.buscarRGF(ctx, "RGF-Anexo 05", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de disponibilidade de caixa encontrado", "ente", idEnte)
		return nil
	}

	var vinculada float64
	var naoVinculada float64

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		conta := strings.ToUpper(item.Conta)
		coluna := strings.ToUpper(item.Coluna)

		if (strings.Contains(rotulo, "LÍQUIDA") || strings.Contains(conta, "LÍQUIDA")) &&
			(strings.Contains(rotulo, "NÃO VINCULADA") || strings.Contains(rotulo, "NAO VINCULADA") ||
				strings.Contains(conta, "NÃO VINCULADA") || strings.Contains(conta, "NAO VINCULADA")) {
			if !strings.Contains(coluna, "%") {
				naoVinculada = item.Valor
			}
		}

		if (strings.Contains(rotulo, "VINCULADA") || strings.Contains(conta, "VINCULADA")) &&
			!strings.Contains(rotulo, "NÃO") && !strings.Contains(rotulo, "NAO") &&
			!strings.Contains(conta, "NÃO") && !strings.Contains(conta, "NAO") {
			if !strings.Contains(coluna, "%") {
				vinculada = item.Valor
			}
		}
	}

	if vinculada == 0 && naoVinculada == 0 {
		return nil
	}

	return &types.DisponibilidadeCaixa{
		Vinculada:    vinculada,
		NaoVinculada: naoVinculada,
		Periodo:      strconv.FormatInt(exercicio, 10),
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarRestosAPagar(ctx context.Context, idEnte int, exercicio int64) *types.RestosAPagar {
	log := logger.New("Municipal: UseCase: buscarRestosAPagar")
	itens, err := u.buscarRGF(ctx, "RGF-Anexo 06", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de restos a pagar encontrado", "ente", idEnte)
		return nil
	}

	var inscritos float64
	var pagos float64
	var cancelados float64

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		conta := strings.ToUpper(item.Conta)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(coluna, "%") {
			continue
		}

		if strings.Contains(rotulo, "INSCRITOS") || strings.Contains(conta, "INSCRITOS") {
			inscritos += item.Valor
		}
		if strings.Contains(rotulo, "PAGOS") || strings.Contains(conta, "PAGOS") {
			pagos += item.Valor
		}
		if strings.Contains(rotulo, "CANCELADOS") || strings.Contains(conta, "CANCELADOS") {
			cancelados += item.Valor
		}
	}

	if inscritos == 0 && pagos == 0 {
		return nil
	}

	return &types.RestosAPagar{
		Inscritos:  inscritos,
		Pagos:      pagos,
		Cancelados: cancelados,
		Periodo:    strconv.FormatInt(exercicio, 10),
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarGastoSaude(ctx context.Context, idEnte int, exercicio int64) *types.GastoSaude {
	log := logger.New("Municipal: UseCase: buscarGastoSaude")
	itens, err := u.buscarRREO(ctx, "RREO-Anexo 09", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de gasto saude encontrado", "ente", idEnte)
		return nil
	}

	var valorTotal float64
	var pctAplicado float64

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(coluna, "%") {
			if strings.Contains(coluna, "SAÚDE") || strings.Contains(rotulo, "APLICADO") ||
				strings.Contains(coluna, "APLICACAO") || strings.Contains(coluna, "APLICAÇÃO") {
				pctAplicado = item.Valor
			}
		} else if strings.Contains(rotulo, "TOTAL") || strings.Contains(rotulo, "DESPESA") {
			if strings.Contains(rotulo, "SAÚDE") || strings.Contains(rotulo, "SAUDE") {
				if strings.Contains(coluna, "EMPENHAD") || strings.Contains(coluna, "LIQUIDAD") || strings.Contains(coluna, "PAG") {
					valorTotal += item.Valor
				}
			}
			if rotulo == "TOTAL" || rotulo == "TOTAL DAS DESPESAS" || rotulo == "DESPESAS" {
				if strings.Contains(coluna, "EMPENHAD") || strings.Contains(coluna, "LIQUIDAD") || strings.Contains(coluna, "PAG") {
					valorTotal += item.Valor
				}
			}
		}
	}

	if valorTotal == 0 && pctAplicado == 0 {
		return nil
	}

	return &types.GastoSaude{
		ValorTotal:           valorTotal,
		PercentualAplicado:   pctAplicado,
		LimiteConstitucional: 15,
		Periodo:              strconv.FormatInt(exercicio, 10),
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarGastoEducacao(ctx context.Context, idEnte int, exercicio int64) *types.GastoEducacao {
	log := logger.New("Municipal: UseCase: buscarGastoEducacao")
	itens, err := u.buscarRREO(ctx, "RREO-Anexo 10", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de gasto educacao encontrado", "ente", idEnte)
		return nil
	}

	var valorTotal float64
	var pctAplicado float64

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(coluna, "%") {
			if strings.Contains(coluna, "EDUCAÇ") || strings.Contains(rotulo, "APLICADO") ||
				strings.Contains(coluna, "APLICACAO") || strings.Contains(coluna, "APLICAÇÃO") {
				pctAplicado = item.Valor
			}
		} else if strings.Contains(rotulo, "TOTAL") || strings.Contains(rotulo, "DESPESA") {
			if strings.Contains(rotulo, "EDUCA") || strings.Contains(rotulo, "ENSINO") {
				if strings.Contains(coluna, "EMPENHAD") || strings.Contains(coluna, "LIQUIDAD") || strings.Contains(coluna, "PAG") {
					valorTotal += item.Valor
				}
			}
			if rotulo == "TOTAL" || rotulo == "TOTAL DAS DESPESAS" || rotulo == "DESPESAS" {
				if strings.Contains(coluna, "EMPENHAD") || strings.Contains(coluna, "LIQUIDAD") || strings.Contains(coluna, "PAG") {
					valorTotal += item.Valor
				}
			}
		}
	}

	if valorTotal == 0 && pctAplicado == 0 {
		return nil
	}

	return &types.GastoEducacao{
		ValorTotal:           valorTotal,
		PercentualAplicado:   pctAplicado,
		LimiteConstitucional: 25,
		Periodo:              strconv.FormatInt(exercicio, 10),
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarFundeb(ctx context.Context, idEnte int, exercicio int64) *types.FundebResumo {
	log := logger.New("Municipal: UseCase: buscarFundeb")
	itens, err := u.buscarRREO(ctx, "RREO-Anexo 08", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de fundeb encontrado", "ente", idEnte)
		return nil
	}

	var receitaTotal float64
	var despesaTotal float64

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		conta := strings.ToUpper(item.Conta)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(coluna, "%") {
			continue
		}

		if strings.Contains(rotulo, "RECEITA") || strings.Contains(conta, "RECEITA") {
			if strings.Contains(rotulo, "FUNDEB") || strings.Contains(conta, "FUNDEB") {
				receitaTotal += item.Valor
			}
		}
		if strings.Contains(rotulo, "DESPESA") || strings.Contains(conta, "DESPESA") {
			if strings.Contains(rotulo, "FUNDEB") || strings.Contains(conta, "FUNDEB") {
				despesaTotal += item.Valor
			}
		}
	}

	if receitaTotal == 0 && despesaTotal == 0 {
		for _, item := range itens {
			coluna := strings.ToUpper(item.Coluna)
			if strings.Contains(coluna, "%") {
				continue
			}
			if strings.Contains(coluna, "EMPENHAD") || strings.Contains(coluna, "LIQUIDAD") || strings.Contains(coluna, "PAG") {
				despesaTotal += item.Valor
			} else if !strings.Contains(coluna, "EMPENHAD") {
				receitaTotal += item.Valor
			}
		}
	}

	if receitaTotal == 0 && despesaTotal == 0 {
		return nil
	}

	return &types.FundebResumo{
		ReceitaTotal: receitaTotal,
		DespesaTotal: despesaTotal,
		Periodo:      strconv.FormatInt(exercicio, 10),
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarDespesasPorGrupo(ctx context.Context, idEnte int, exercicio int64) []types.DespesaPorGrupoItem {
	log := logger.New("Municipal: UseCase: buscarDespesasPorGrupo")
	itens, err := u.buscarRREO(ctx, "RREO-Anexo 05", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de despesas por grupo encontrado", "ente", idEnte)
		return nil
	}

	grupos := make(map[string]*types.DespesaPorGrupoItem)
	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(coluna, "%") {
			continue
		}

		var chave string
		switch {
		case strings.Contains(rotulo, "CORRENTE"):
			chave = "Corrente"
		case strings.Contains(rotulo, "CAPITAL"):
			chave = "Capital"
		case strings.Contains(rotulo, "PESSOAL"):
			chave = "Pessoal"
		case strings.Contains(rotulo, "JUROS") || strings.Contains(rotulo, "ENCARGOS"):
			chave = "Juros e Encargos"
		case strings.Contains(rotulo, "INVESTIMENT") || strings.Contains(rotulo, "INVESTIMENTO"):
			chave = "Investimentos"
		case strings.Contains(rotulo, "INVERSÃO") || strings.Contains(rotulo, "INVERSAO"):
			chave = "Inversões Financeiras"
		case strings.Contains(rotulo, "AMORTIZA") || strings.Contains(rotulo, "AMORTIZACAO"):
			chave = "Amortização"
		default:
			continue
		}

		if _, ok := grupos[chave]; !ok {
			grupos[chave] = &types.DespesaPorGrupoItem{Grupo: chave}
		}
		switch {
		case strings.Contains(coluna, "EMPENHAD"):
			grupos[chave].Empenhado += item.Valor
		case strings.Contains(coluna, "LIQUIDAD"):
			grupos[chave].Liquidado += item.Valor
		case strings.Contains(coluna, "PAG"):
			grupos[chave].Pago += item.Valor
		}
	}

	if len(grupos) == 0 {
		return nil
	}

	resultado := make([]types.DespesaPorGrupoItem, 0, len(grupos))
	for _, g := range grupos {
		resultado = append(resultado, *g)
	}
	return resultado
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarTransferencias(ctx context.Context, idEnte int, exercicio int64) []types.TransferenciaItem {
	log := logger.New("Municipal: UseCase: buscarTransferencias")
	itens, err := u.buscarRREO(ctx, "RREO-Anexo 07", idEnte, exercicio)
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de transferencias encontrado", "ente", idEnte)
		return nil
	}

	var transferencias []types.TransferenciaItem
	orgaos := make(map[string]float64)

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(coluna, "%") {
			continue
		}

		nome := item.Conta
		if nome == "" {
			nome = item.Rotulo
		}
		if nome == "" {
			continue
		}
		if rotulo == "TOTAL" || rotulo == "TOTAL DAS TRANSFERÊNCIAS" || rotulo == "TOTAL DAS TRANSFERENCIAS" {
			continue
		}

		if strings.Contains(coluna, "RECEITA") || strings.Contains(coluna, "ARRECAD") || strings.Contains(coluna, "PREVIST") {
			orgaos[nome] += item.Valor
		} else if !strings.Contains(coluna, "EMPENHAD") && !strings.Contains(coluna, "LIQUIDAD") && !strings.Contains(coluna, "PAG") {
			orgaos[nome] += item.Valor
		}
	}

	for nome, valor := range orgaos {
		transferencias = append(transferencias, types.TransferenciaItem{
			Orgao: nome,
			Valor: valor,
		})
	}

	if len(transferencias) == 0 {
		return nil
	}
	return transferencias
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarBalancoPatrimonial(ctx context.Context, idEnte int, exercicio int64) *types.BalancoPatrimonial {
	log := logger.New("Municipal: UseCase: buscarBalancoPatrimonial")

	alvo := exercicio
	if alvo <= 0 {
		alvo = int64(time.Now().Year() - 1)
	}

	var itens []siconfi.DCAItem
	var err error
	for ano := alvo; ano >= alvo-2; ano-- {
		itens, err = u.siconfiClient.BuscarDCA(ctx, ano, idEnte, "DCA-Anexo I-AB")
		if err == nil && len(itens) > 0 {
			break
		}
	}
	if err != nil || len(itens) == 0 {
		log.Warn("nenhum dado de balanco patrimonial encontrado", "ente", idEnte)
		if errors.Is(err, siconfi.ErrSICONFIIndisponivel) {
			u.apiIndisponivel.Store(true)
		}
		return nil
	}

	var ativoCirc, ativoNaoCirc, passivoCirc, passivoNaoCirc, pl float64

	for _, item := range itens {
		rotulo := strings.ToUpper(item.Rotulo)
		coluna := strings.ToUpper(item.Coluna)

		if strings.Contains(coluna, "%") {
			continue
		}

		switch {
		case strings.Contains(rotulo, "TOTAL DO ATIVO CIRCULANTE") || strings.Contains(rotulo, "ATIVO CIRCULANTE"):
			ativoCirc = item.Valor
		case strings.Contains(rotulo, "TOTAL DO ATIVO NÃO CIRCULANTE") || strings.Contains(rotulo, "TOTAL DO ATIVO NAO CIRCULANTE") ||
			strings.Contains(rotulo, "ATIVO NÃO CIRCULANTE") || strings.Contains(rotulo, "ATIVO NAO CIRCULANTE"):
			ativoNaoCirc = item.Valor
		case strings.Contains(rotulo, "TOTAL DO PASSIVO CIRCULANTE") || strings.Contains(rotulo, "PASSIVO CIRCULANTE"):
			passivoCirc = item.Valor
		case strings.Contains(rotulo, "TOTAL DO PASSIVO NÃO CIRCULANTE") || strings.Contains(rotulo, "TOTAL DO PASSIVO NAO CIRCULANTE") ||
			strings.Contains(rotulo, "PASSIVO NÃO CIRCULANTE") || strings.Contains(rotulo, "PASSIVO NAO CIRCULANTE"):
			passivoNaoCirc = item.Valor
		case strings.Contains(rotulo, "PATRIMÔNIO LÍQUIDO") || strings.Contains(rotulo, "PATRIMONIO LIQUIDO") ||
			strings.Contains(rotulo, "PATRIMÔNIO LÍQUIDO") || strings.Contains(rotulo, "PATRIMONIO LÍQUIDO"):
			pl = item.Valor
		}
	}

	if ativoCirc == 0 && ativoNaoCirc == 0 && passivoCirc == 0 && passivoNaoCirc == 0 && pl == 0 {
		return nil
	}

	return &types.BalancoPatrimonial{
		AtivoCirculante:      ativoCirc,
		AtivoNaoCirculante:   ativoNaoCirc,
		PassivoCirculante:    passivoCirc,
		PassivoNaoCirculante: passivoNaoCirc,
		PatrimonioLiquido:    pl,
		Periodo:              strconv.FormatInt(alvo, 10),
	}
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarContratos(ctx context.Context, codigoIBGE int, ano int) []pncpClient.Contrato {
	log := logger.New("Municipal: UseCase: buscarContratos")
	if ano <= 0 {
		ano = int(time.Now().Year() - 1)
	}

	codigoStr := strconv.Itoa(codigoIBGE)
	meses := utils.ExtrairMesesDoAno(ano)
	total := make([]pncpClient.Contrato, 0)

	for _, am := range meses {
		jaRealizada, err := u.repo.BuscaJaRealizada(ctx, "municipio", codigoStr, am.Ano, am.Mes)
		if err != nil {
			log.Warn("erro ao verificar busca no PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
			continue
		}
		if jaRealizada {
			persistidos, err := u.repo.BuscarContratosPorFiltro(ctx, "municipio", codigoStr, am.Ano, am.Mes)
			if err != nil {
				log.Warn("erro ao buscar contratos do PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
				continue
			}
			for i := range persistidos {
				total = append(total, repositorios.PersistidoParaContrato(persistidos[i]))
			}
			continue
		}

		contratosMes := u.buscarMesComLock(ctx, "municipio", codigoStr, am.Ano, am.Mes)
		total = append(total, contratosMes...)
	}

	log.Info("contratos encontrados", "codigo_ibge", codigoIBGE, "ano", ano, "total", len(total))
	if len(total) == 0 {
		return nil
	}
	return total
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarMesComLock(ctx context.Context, tipo, valor string, ano, mes int) []pncpClient.Contrato {
	fetch := func(c context.Context, t string, v string, a, m int) ([]pncpClient.Contrato, error) {
		return u.buscarTodasPaginas(c, t, v, a, m), nil
	}
	persist := func(c context.Context, t string, v string, a, m int, contratos []pncpClient.Contrato) error {
		return persistirContratosMunicipio(c, u.repo, t, v, a, m, contratos)
	}
	return pncpbusca.BuscarMesComLock(
		ctx, u.redis, u.repo, fetch, persist,
		tipo, valor, ano, mes,
	)
}

func persistirContratosMunicipio(ctx context.Context, repo repositorios.PNCPRepository, tipo, valor string, ano, mes int, contratos []pncpClient.Contrato) error {
	if len(contratos) == 0 {
		controle := repositorios.BuscaControlePersistido{
			TipoBusca: tipo, ValorBusca: valor,
			Ano: ano, Mes: mes,
			DataInicial:              time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC),
			DataFinal:                time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1),
			TotalContratosEncontrados: 0,
		}
		return repo.RegistrarBusca(ctx, controle)
	}

	persistidos := make([]repositorios.ContratoPersistido, len(contratos))
	for i, c := range contratos {
		persistidos[i] = repositorios.ContratoParaPersistido(c)
	}
	if err := repo.SalvarContratos(ctx, persistidos); err != nil {
		return err
	}

	controle := repositorios.BuscaControlePersistido{
		TipoBusca: tipo, ValorBusca: valor,
		Ano: ano, Mes: mes,
		DataInicial:              time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC),
		DataFinal:                time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1),
		TotalContratosEncontrados: len(contratos),
	}
	return repo.RegistrarBusca(ctx, controle)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) buscarTodasPaginas(ctx context.Context, tipo, valor string, ano, mes int) []pncpClient.Contrato {
	log := logger.New("Municipal: UseCase: buscarTodasPaginas")

	dataInicial, dataFinal := utils.FormatarPeriodoMes(ano, mes)

	pagina := 1
	tamanho := 50
	items := make([]pncpClient.Contrato, 0)
	seen := make(map[string]struct{})

	for {
		resp, err := u.pncpClient.BuscarContratacoesPorMunicipio(ctx, valor, dataInicial, dataFinal, "", pagina, tamanho)
		if err != nil {
			log.Error("erro ao buscar contratacoes municipio", "pagina", pagina, "erro", err)
			break
		}

		if resp == nil || len(resp.Data) == 0 {
			break
		}

		for _, c := range resp.Data {
			key := ""
			if c.NumeroControlePNCP != nil {
				key = *c.NumeroControlePNCP
			}
			if key == "" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			items = append(items, c)
		}

		if pagina >= resp.TotalPaginas {
			break
		}
		pagina++

		select {
		case <-ctx.Done():
			return items
		default:
		}
	}

	log.Info("contratos obtidos PNCP municipio", "valor", valor, "ano", ano, "mes", mes, "total", len(items))
	return items
}
