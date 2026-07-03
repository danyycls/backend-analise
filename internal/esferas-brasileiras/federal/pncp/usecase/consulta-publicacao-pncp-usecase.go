package usecase

import (
	"context"
	"encoding/json"

	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
)

type ConsultaPublicacaoPNCPUseCase struct {
	pncpClient     *pncp.PNCPClient
	opencnpjClient *opencnpj.OpenCNPJClient
	redis          *redis.RedisCache
	licitacaoCache *redis.LicitacaoCache
}

func NovoConsultaPublicacaoPNCPUseCase(pncp *pncp.PNCPClient, opencnpj *opencnpj.OpenCNPJClient, redis *redis.RedisCache, licitacaoCache *redis.LicitacaoCache) *ConsultaPublicacaoPNCPUseCase {
	return &ConsultaPublicacaoPNCPUseCase{
		pncpClient:     pncp,
		opencnpjClient: opencnpj,
		redis:          redis,
		licitacaoCache: licitacaoCache,
	}
}

func (u *ConsultaPublicacaoPNCPUseCase) BuscarPorUF(ctx context.Context, uf, dataInicial, dataFinal, codigoModalidade string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	return u.buscar(ctx, "uf", uf, dataInicial, dataFinal, codigoModalidade, paginasErro)
}

func (u *ConsultaPublicacaoPNCPUseCase) BuscarPorMunicipio(ctx context.Context, codigoMunicipio, dataInicial, dataFinal, codigoModalidade string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	return u.buscar(ctx, "municipio", codigoMunicipio, dataInicial, dataFinal, codigoModalidade, paginasErro)
}

func (u *ConsultaPublicacaoPNCPUseCase) buscar(ctx context.Context, tipo, valor, dataInicial, dataFinal, codigoModalidade string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: buscar")
	cacheParams := map[string]interface{}{
		"tipo":                        tipo,
		"valor":                       valor,
		"dataInicial":                 dataInicial,
		"dataFinal":                   dataFinal,
		"codigoModalidadeContratacao": codigoModalidade,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache(redis.ChaveLicitacoesTrimestre, raw)

	var cached []*pncp.AnaliseResultado
	if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	} else if ok {
		log.Info("cache hit", "tipo", tipo, "valor", valor)
		return cached, nil
	}

	contratosCache, encontrados, err := u.licitacaoCache.BuscarPorFiltros(ctx, tipo, valor, dataInicial, dataFinal)
	if err == nil && encontrados {
		log.Info("entity cache hit", "tipo", tipo, "valor", valor, "contratos", len(contratosCache))
		resultados, err := u.montarResultados(ctx, contratosCache, dataInicial, dataFinal)
		if err != nil {
			log.Warn("erro ao montar resultados do cache", "erro", err)
		} else {
			return resultados, nil
		}
	}

	items, paginas := u.buscarTodasPaginas(ctx, tipo, valor, dataInicial, dataFinal, codigoModalidade)
	if paginasErro != nil {
		*paginasErro = paginas
	}

	log.Info("total de itens", "total", len(items))

	resultados, err := u.montarResultados(ctx, items, dataInicial, dataFinal)
	if err != nil {
		return nil, err
	}

	if err := u.redis.Set(ctx, cacheKey, resultados); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}

	return resultados, nil
}

func (u *ConsultaPublicacaoPNCPUseCase) montarResultados(ctx context.Context, contratos []pncp.Contrato, dataInicial, dataFinal string) ([]*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: montarResultados")

	fornecedoresMap := u.extrairFornecedores(contratos)
	log.Info("fornecedores unicos", "total", len(fornecedoresMap))

	enrichedFornecedores := u.enriquecerFornecedores(ctx, fornecedoresMap)
	contratos = u.montarContratosDTO(contratos, enrichedFornecedores)

	grupos := u.agruparPorOrgao(contratos)
	log.Info("grupos por orgao", "total", len(grupos))

	fornecedoresContados := make(map[string]struct{})
	for _, g := range grupos {
		for _, c := range g.Contratos {
			if c.Fornecedor != nil && c.Fornecedor.CNPJ != nil {
				fornecedoresContados[*c.Fornecedor.CNPJ] = struct{}{}
			} else if c.NIFornecedor != nil && *c.NIFornecedor != "" {
				fornecedoresContados[*c.NIFornecedor] = struct{}{}
			}
		}
	}

	resultados := make([]*pncp.AnaliseResultado, 0, len(grupos))
	for _, g := range grupos {
		totalContratos := len(g.Contratos)
		totalEmpresas := len(fornecedoresContados)
		var valorTotal float64
		for _, c := range g.Contratos {
			if c.ValorGlobal != nil {
				valorTotal += *c.ValorGlobal
			} else if c.ValorTotalEstimado != nil {
				valorTotal += *c.ValorTotalEstimado
			}
		}

		resultados = append(resultados, &pncp.AnaliseResultado{
			Orgao: &pncp.OrgaoInfo{
				CNPJ:        g.CNPJ,
				RazaoSocial: g.RazaoSocial,
			},
			Periodo: &pncp.Periodo{
				DataInicial: pncp.StrPtr(dataInicial),
				DataFinal:   pncp.StrPtr(dataFinal),
			},
			Resumo: &pncp.Resumo{
				TotalContratos:      &totalContratos,
				TotalEmpresas:       &totalEmpresas,
				ValorTotalContratos: &valorTotal,
			},
			Contratos: g.Contratos,
		})
	}

	return resultados, nil
}

func (u *ConsultaPublicacaoPNCPUseCase) buscarTodasPaginas(ctx context.Context, tipo, valor, dataInicial, dataFinal, codigoModalidade string) ([]pncp.Contrato, []int) {
	log := logger.New("PNCP: UseCase: buscarTodasPaginas")
	pagina := 1
	tamanho := 50
	items := make([]pncp.Contrato, 0)
	seen := make(map[string]struct{})
	totalPaginas := 1
	paginasComErro := make([]int, 0)

	for pagina <= totalPaginas {
		cacheParams := map[string]interface{}{
			"tipo":                        tipo,
			"valor":                       valor,
			"dataInicial":                 dataInicial,
			"dataFinal":                   dataFinal,
			"codigoModalidadeContratacao": codigoModalidade,
			"pagina":                      pagina,
			"tamanho":                     tamanho,
		}
		raw, _ := json.Marshal(cacheParams)
		cacheKey := redis.ChaveCache(redis.ChavePNCPublicacaoPagina, raw)

		var cached *pncp.PublicacaoResponse
		cacheHit := false
		if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
			log.Warn("cache indisponivel", "erro", err)
		} else if ok && cached != nil {
			cacheHit = true
			log.Info("cache hit pagina", "pagina", pagina)
		}

		var resp *pncp.PublicacaoResponse
		var err error

		if cacheHit {
			resp = cached
			totalPaginas = resp.TotalPaginas
		} else {
			if tipo == "uf" {
				resp, err = u.pncpClient.BuscarContratacoesPorUF(ctx, valor, dataInicial, dataFinal, codigoModalidade, pagina, tamanho)
			} else {
				resp, err = u.pncpClient.BuscarContratacoesPorMunicipio(ctx, valor, dataInicial, dataFinal, codigoModalidade, pagina, tamanho)
			}

			if err != nil {
				log.Warn("erro ao buscar pagina, ignorando", "pagina", pagina, "erro", err)
				paginasComErro = append(paginasComErro, pagina)
				break
			}

			totalPaginas = resp.TotalPaginas

			if setErr := u.redis.Set(ctx, cacheKey, resp); setErr != nil {
				log.Warn("cache indisponivel", "erro", setErr)
			}

			if idxErr := u.licitacaoCache.IndexarContratos(ctx, resp.Data); idxErr != nil {
				log.Warn("erro ao indexar contratos", "erro", idxErr)
			}
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

		if len(items) == 0 || pagina >= totalPaginas {
			break
		}
		pagina++
	}

	return items, paginasComErro
}

type grupoOrgao struct {
	CNPJ        *string
	RazaoSocial *string
	Contratos   []pncp.Contrato
}

func (u *ConsultaPublicacaoPNCPUseCase) extrairFornecedores(contratos []pncp.Contrato) map[string]string {
	fornecedores := make(map[string]string)
	for _, c := range contratos {
		if c.NIFornecedor == nil || *c.NIFornecedor == "" {
			continue
		}
		ni := utils.NormalizarCNPJ(*c.NIFornecedor)
		if _, ok := fornecedores[ni]; !ok {
			nome := ""
			if c.NomeRazaoSocialFornecedor != nil {
				nome = *c.NomeRazaoSocialFornecedor
			}
			fornecedores[ni] = nome
		}
	}
	return fornecedores
}

func (u *ConsultaPublicacaoPNCPUseCase) enriquecerFornecedores(ctx context.Context, fornecedoresMap map[string]string) map[string]*types.FornecedorOpenCNPJ {
	log := logger.New("PNCP: UseCase: enriquecerFornecedores")
	enriched := make(map[string]*types.FornecedorOpenCNPJ, len(fornecedoresMap))
	for cnpjF, nome := range fornecedoresMap {
		data, err := u.opencnpjClient.Buscar(ctx, cnpjF)
		if err != nil {
			log.Error("erro ao consultar OpenCNPJ", "cnpj", cnpjF, "erro", err)
			enriched[cnpjF] = &types.FornecedorOpenCNPJ{CNPJ: pncp.StrPtr(cnpjF), RazaoSocial: pncp.StrPtr(nome)}
			continue
		}
		enriched[cnpjF] = utils.BuildFornecedorDTO(data)
	}
	return enriched
}

func (u *ConsultaPublicacaoPNCPUseCase) montarContratosDTO(contratos []pncp.Contrato, enrichedFornecedores map[string]*types.FornecedorOpenCNPJ) []pncp.Contrato {
	contratosDTO := make([]pncp.Contrato, len(contratos))
	copy(contratosDTO, contratos)
	for i, c := range contratosDTO {
		if c.NIFornecedor == nil {
			continue
		}
		ni := utils.NormalizarCNPJ(*c.NIFornecedor)
		if enriched, ok := enrichedFornecedores[ni]; ok {
			contratosDTO[i].Fornecedor = enriched
		}
	}
	return contratosDTO
}

func (u *ConsultaPublicacaoPNCPUseCase) agruparPorOrgao(contratos []pncp.Contrato) []*grupoOrgao {
	grupos := make(map[string]*grupoOrgao)
	ordem := make([]string, 0)

	for _, c := range contratos {
		cnpj := ""
		razao := ""
		if c.OrgaoEntidade != nil && c.OrgaoEntidade.CNPJ != nil {
			cnpj = *c.OrgaoEntidade.CNPJ
		}
		if c.OrgaoEntidade != nil && c.OrgaoEntidade.RazaoSocial != nil {
			razao = *c.OrgaoEntidade.RazaoSocial
		}

		if cnpj == "" {
			cnpj = "sem_cnpj"
		}

		if _, ok := grupos[cnpj]; !ok {
			grupos[cnpj] = &grupoOrgao{
				CNPJ:        pncp.StrPtr(cnpj),
				RazaoSocial: pncp.StrPtr(razao),
				Contratos:   make([]pncp.Contrato, 0),
			}
			ordem = append(ordem, cnpj)
		}
		grupos[cnpj].Contratos = append(grupos[cnpj].Contratos, c)
	}

	resultado := make([]*grupoOrgao, 0, len(grupos))
	for _, key := range ordem {
		resultado = append(resultado, grupos[key])
	}
	return resultado
}
