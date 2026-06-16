package usecase

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	pncpClient "github.com/danyele/laceu/internal/shared/clients/pncp"
	portalClient "github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
	"github.com/danyele/laceu/internal/shared/clients/siconfi"
	"github.com/danyele/laceu/internal/shared/logger"
	redis "github.com/danyele/laceu/internal/shared/redis"
	"github.com/danyele/laceu/internal/shared/types"
)

type EsferaMunicipalBuscarDetalhesUseCase struct {
	siconfiClient *siconfi.SICONFIClient
	portalClient  *portalClient.PortalTransparenciaClient
	pncpClient    *pncpClient.PNCPClient
	redis         *redis.RedisCache
}

func NovoEsferaMunicipalBuscarDetalhesUseCase(
	siconfiCli *siconfi.SICONFIClient,
	portalCli *portalClient.PortalTransparenciaClient,
	pncpCli *pncpClient.PNCPClient,
	redis *redis.RedisCache,
) *EsferaMunicipalBuscarDetalhesUseCase {
	return &EsferaMunicipalBuscarDetalhesUseCase{
		siconfiClient: siconfiCli,
		portalClient:  portalCli,
		pncpClient:    pncpCli,
		redis:         redis,
	}
}

func anoAlvo() int64 {
	return int64(time.Now().Year() - 1)
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarDespesaPessoal(ctx context.Context, codigoIBGE int, exercicio int64) *types.DespesaPessoalResumo {
	log := logger.New("Municipal: UseCase: BuscarDespesaPessoal")
	idEnte := codigoIBGE
	if exercicio == 0 {
		exercicio = anoAlvo()
	}

	tentativas := []struct {
		ano           int64
		periodicidade string
		periodo       int
	}{
		{exercicio, "Q", 3},
		{exercicio, "S", 2},
		{exercicio - 1, "Q", 3},
		{exercicio - 1, "S", 2},
		{exercicio - 2, "Q", 3},
		{exercicio - 2, "S", 2},
	}

	for _, t := range tentativas {
		params := siconfi.RGFParams{
			AnExercicio:         t.ano,
			InPeriodicidade:     t.periodicidade,
			NrPeriodo:           t.periodo,
			CoTipoDemonstrativo: "RGF",
			CoPoder:             "E",
			IdEnte:              idEnte,
			CoEsfera:            "M",
			NoAnexo:             "RGF-Anexo 01",
		}

		raw, _ := json.Marshal(params)
		cacheKey := redis.ChaveCache("municipal-rgf", raw)

		var cached []siconfi.RGFItem
		cacheHit := false
		if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
			log.Warn("cache indisponivel", "erro", err)
		} else if ok {
			cacheHit = true
		}

		var itens []siconfi.RGFItem
		var apiErr error
		if cacheHit {
			itens = cached
		} else {
			itens, apiErr = u.siconfiClient.BuscarRGF(ctx, params)
			if apiErr != nil {
				log.Error("erro ao buscar RGF", "ente", idEnte, "exercicio", t.ano, "periodicidade", t.periodicidade, "periodo", t.periodo, "erro", apiErr)
			} else {
				if setErr := u.redis.Set(ctx, cacheKey, itens); setErr != nil {
					log.Warn("cache indisponivel", "erro", setErr)
				}
			}
		}

		if apiErr != nil || len(itens) == 0 {
			continue
		}

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
				Periodo:       strconv.FormatInt(t.ano, 10),
			}
		}
	}

	return nil
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarRREO(ctx context.Context, codigoIBGE int, exercicio int64) ([]types.GastoPorFuncao, []types.ReceitaResumo) {
	log := logger.New("Municipal: UseCase: BuscarRREO")
	idEnte := codigoIBGE
	if exercicio == 0 {
		exercicio = anoAlvo()
	}

	var gastos []types.GastoPorFuncao
	var receitas []types.ReceitaResumo

	periods := []int{6, 5}

	for _, periodo := range periods {
		if len(gastos) == 0 {
			params := siconfi.RREOParams{
				AnExercicio:         exercicio,
				NrPeriodo:           periodo,
				CoTipoDemonstrativo: "RREO",
				IdEnte:              idEnte,
				NoAnexo:             "RREO-Anexo 02",
				CoEsfera:            "M",
			}

			raw, _ := json.Marshal(params)
			cacheKey := redis.ChaveCache("municipal-rreo", raw)

			var cached []siconfi.RREOItem
			cacheHit := false
			if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
				log.Warn("cache indisponivel", "erro", err)
			} else if ok {
				cacheHit = true
			}

			var itens []siconfi.RREOItem
			var apiErr error
			if cacheHit {
				itens = cached
			} else {
				itens, apiErr = u.siconfiClient.BuscarRREO(ctx, params)
				if apiErr == nil {
					if setErr := u.redis.Set(ctx, cacheKey, itens); setErr != nil {
						log.Warn("cache indisponivel", "erro", setErr)
					}
				}
			}

			if apiErr == nil && len(itens) > 0 {
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
				for _, g := range despesasPorFuncao {
					gastos = append(gastos, *g)
				}
				sort.Slice(gastos, func(i, j int) bool {
					return gastos[i].Empenhado > gastos[j].Empenhado
				})
			}
		}

		if len(receitas) == 0 {
			receitasParams := siconfi.RREOParams{
				AnExercicio:         exercicio,
				NrPeriodo:           periodo,
				CoTipoDemonstrativo: "RREO",
				IdEnte:              idEnte,
				NoAnexo:             "RREO-Anexo 03",
				CoEsfera:            "M",
			}

			raw, _ := json.Marshal(receitasParams)
			cacheKey := redis.ChaveCache("municipal-rreo", raw)

			var cached []siconfi.RREOItem
			cacheHit := false
			if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
				log.Warn("cache indisponivel", "erro", err)
			} else if ok {
				cacheHit = true
			}

			var itens []siconfi.RREOItem
			var apiErr error
			if cacheHit {
				itens = cached
			} else {
				itens, apiErr = u.siconfiClient.BuscarRREO(ctx, receitasParams)
				if apiErr == nil {
					if setErr := u.redis.Set(ctx, cacheKey, itens); setErr != nil {
						log.Warn("cache indisponivel", "erro", setErr)
					}
				}
			}

			if apiErr == nil && len(itens) > 0 {
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
			}
		}

		if len(gastos) > 0 && len(receitas) > 0 {
			break
		}
	}

	return gastos, receitas
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarRecursosFederais(ctx context.Context, codigoIBGE int) []types.RecursoFederalRecebido {
	log := logger.New("Municipal: UseCase: BuscarRecursosFederais")
	anoAlvo := time.Now().Year() - 1

	filtro := portalClient.DespesaRecursosRecebidosQueryParams{
		Pagina:       1,
		MesAnoInicio: strconv.Itoa(anoAlvo) + "-01",
		MesAnoFim:    strconv.Itoa(anoAlvo) + "-12",
		CodigoIBGE:   strconv.Itoa(codigoIBGE),
	}

	raw, _ := json.Marshal(filtro)
	cacheKey := redis.ChaveCache("municipal-recursos-federais", raw)

	var cached []types.RecursoFederalRecebido
	if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	} else if ok {
		return cached
	}

	itens, err := u.portalClient.ListarRecursosRecebidos(ctx, filtro)
	if err != nil {
		log.Error("erro ao buscar recursos recebidos para municipio", "codigo_ibge", codigoIBGE, "erro", err)
		return nil
	}

	result := make([]types.RecursoFederalRecebido, 0, len(itens))
	for _, item := range itens {
		result = append(result, types.RecursoFederalRecebido{
			NomePessoa:        item.NomePessoa,
			TipoPessoa:        item.TipoPessoa,
			NomeUG:            item.NomeUG,
			NomeOrgao:         item.NomeOrgao,
			NomeOrgaoSuperior: item.NomeOrgaoSuperior,
			Valor:             item.Valor,
			MesAno:            item.AnoMes,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Valor > result[j].Valor
	})

	if err := u.redis.Set(ctx, cacheKey, result); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}

	return result
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarContratos(ctx context.Context, codigoIBGE int) []types.ContratoPNCP {
	log := logger.New("Municipal: UseCase: BuscarContratos")
	codigoStr := strconv.Itoa(codigoIBGE)
	anoAlvo := time.Now().Year() - 1
	dataInicial := strconv.Itoa(anoAlvo) + "0101"
	dataFinal := strconv.Itoa(anoAlvo) + "1231"

	cacheParams := map[string]interface{}{
		"codigoIBGE":  codigoIBGE,
		"dataInicial": dataInicial,
		"dataFinal":   dataFinal,
		"pagina":      1,
		"tamanho":     20,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache("municipal-contratos-pncp", raw)

	var cached []types.ContratoPNCP
	if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	} else if ok {
		return cached
	}

	resp, err := u.pncpClient.BuscarContratacoesPorMunicipio(ctx, codigoStr, dataInicial, dataFinal, 1, 20)
	if err != nil {
		log.Error("erro ao buscar contratos PNCP para municipio", "codigo_ibge", codigoIBGE, "erro", err)
		return nil
	}

	if resp == nil || len(resp.Data) == 0 {
		return nil
	}

	result := make([]types.ContratoPNCP, 0, len(resp.Data))
	for _, c := range resp.Data {
		lic := contratoParaPNCP(c)
		result = append(result, lic)
	}

	if err := u.redis.Set(ctx, cacheKey, result); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}

	return result
}

func (u *EsferaMunicipalBuscarDetalhesUseCase) BuscarServidores(ctx context.Context, codigoIBGE int) []types.ServidorMunicipio {
	log := logger.New("Municipal: UseCase: BuscarServidores")
	idEnte := codigoIBGE
	exercicio := anoAlvo()

	tentativas := []struct {
		ano           int64
		periodicidade string
		periodo       int
	}{
		{exercicio, "Q", 3},
		{exercicio, "S", 2},
		{exercicio - 1, "Q", 3},
		{exercicio - 1, "S", 2},
	}

	for _, t := range tentativas {
		params := siconfi.RGFParams{
			AnExercicio:         t.ano,
			InPeriodicidade:     t.periodicidade,
			NrPeriodo:           t.periodo,
			CoTipoDemonstrativo: "RGF",
			CoPoder:             "E",
			IdEnte:              idEnte,
			CoEsfera:            "M",
			NoAnexo:             "RGF-Anexo 01",
		}

		raw, _ := json.Marshal(params)
		cacheKey := redis.ChaveCache("municipal-rgf", raw)

		var cached []siconfi.RGFItem
		cacheHit := false
		if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
			log.Warn("cache indisponivel", "erro", err)
		} else if ok {
			cacheHit = true
		}

		var itens []siconfi.RGFItem
		var apiErr error
		if cacheHit {
			itens = cached
		} else {
			itens, apiErr = u.siconfiClient.BuscarRGF(ctx, params)
			if apiErr == nil {
				if setErr := u.redis.Set(ctx, cacheKey, itens); setErr != nil {
					log.Warn("cache indisponivel", "erro", setErr)
				}
			}
		}

		if apiErr != nil || len(itens) == 0 {
			continue
		}

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

		if len(despesasPorCategoria) == 0 {
			continue
		}

		result := make([]types.ServidorMunicipio, 0, len(despesasPorCategoria))
		for _, v := range despesasPorCategoria {
			if v.DespesaTotal > 0 {
				result = append(result, *v)
			}
		}
		if len(result) > 0 {
			sort.Slice(result, func(i, j int) bool {
				return result[i].DespesaTotal > result[j].DespesaTotal
			})
			return result
		}
	}

	return nil
}

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func ptrFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func contratoParaPNCP(c pncpClient.Contrato) types.ContratoPNCP {
	modalidadeNome := ""
	if c.ModalidadeNome != nil {
		modalidadeNome = *c.ModalidadeNome
	}

	tipoContratoNome := ""
	if c.TipoContrato != nil && c.TipoContrato.Nome != nil {
		tipoContratoNome = *c.TipoContrato.Nome
	}

	ampLegal := ""
	if c.AmparoLegal != nil && c.AmparoLegal.Descricao != nil {
		ampLegal = *c.AmparoLegal.Descricao
	}

	return types.ContratoPNCP{
		Orgao:                ptrStr(c.NomeOrgao),
		Objeto:               ptrStr(c.ObjetoCompra),
		Valor:                valorContrato(c),
		NomeRazaoSocial:      ptrStr(c.NomeRazaoSocialFornecedor),
		DataVigenciaInicio:   ptrStr(c.DataInicioVigencia),
		DataVigenciaFim:      ptrStr(c.DataTerminoVigencia),
		DataPublicacao:       ptrStr(c.DataPublicacao),
		NumeroContrato:       ptrStr(c.NumeroContrato),
		NumeroControlePNCP:   ptrStr(c.NumeroControlePNCP),
		ModalidadeNome:       modalidadeNome,
		NumeroLicitacao:      ptrStr(c.NumeroLicitação),
		CodigoContrato:       ptrStr(c.CodigoContrato),
		OrigemLicitacao:      ptrStr(c.OrigemLicitação),
		TipoContratoNome:     tipoContratoNome,
		ValorGlobal:          ptrFloat(c.ValorGlobal),
		ValorParcela:         ptrFloat(c.ValorParcela),
		ValorTotalEstimado:   ptrFloat(c.ValorTotalEstimado),
		ValorTotalHomologado: ptrFloat(c.ValorTotalHomologado),
		AnoContrato:          ptrInt(c.AnoContrato),
		DataAssinatura:       ptrStr(c.DataAssinatura),
		AmpLegalDescricao:    ampLegal,
		Produto:              ptrStr(c.Produto),
		SubtipoContrato:      ptrStr(c.SubtipoContrato),
	}
}

func valorContrato(c pncpClient.Contrato) float64 {
	if c.ValorInicial != nil {
		return *c.ValorInicial
	}
	if c.ValorGlobal != nil {
		return *c.ValorGlobal
	}
	return 0
}
