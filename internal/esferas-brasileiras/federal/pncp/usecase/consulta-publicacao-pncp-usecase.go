package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/pncpbusca"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/utils"
)

type ConsultaPublicacaoPNCPUseCase struct {
	pncpClient     *pncp.PNCPClient
	opencnpjClient *opencnpj.OpenCNPJClient
	redis          *redis.RedisCache
	repo           repositorios.PNCPRepository
}

func NovoConsultaPublicacaoPNCPUseCase(pncp *pncp.PNCPClient, opencnpj *opencnpj.OpenCNPJClient, redis *redis.RedisCache, repo repositorios.PNCPRepository) *ConsultaPublicacaoPNCPUseCase {
	return &ConsultaPublicacaoPNCPUseCase{
		pncpClient:     pncp,
		opencnpjClient: opencnpj,
		redis:          redis,
		repo:           repo,
	}
}

func (u *ConsultaPublicacaoPNCPUseCase) BuscarPorUF(ctx context.Context, uf, dataInicial, dataFinal, _ string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	return u.buscar(ctx, "uf", uf, dataInicial, dataFinal, paginasErro)
}

func (u *ConsultaPublicacaoPNCPUseCase) BuscarPorMunicipio(ctx context.Context, codigoMunicipio, dataInicial, dataFinal, _ string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	return u.buscar(ctx, "municipio", codigoMunicipio, dataInicial, dataFinal, paginasErro)
}

func (u *ConsultaPublicacaoPNCPUseCase) buscar(ctx context.Context, tipo, valor, dataInicial, dataFinal string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: buscar")

	if result := u.checkCache(ctx, tipo, valor, dataInicial, dataFinal); result != nil {
		return result, nil
	}

	meses := utils.ExtrairMeses(dataInicial, dataFinal)
	if len(meses) == 0 {
		return nil, fmt.Errorf("periodo invalido: %s a %s", dataInicial, dataFinal)
	}

	contratos := u.coletarContratos(ctx, tipo, valor, meses)
	if len(contratos) == 0 {
		return nil, fmt.Errorf("erro ao consultar contratos para %s %s", tipo, valor)
	}

	log.Info("total de itens", "total", len(contratos))

	resultados, err := u.montarResultados(ctx, contratos, dataInicial, dataFinal)
	if err != nil {
		return nil, err
	}

	u.writeCache(ctx, tipo, valor, dataInicial, dataFinal, resultados)
	return resultados, nil
}

func (u *ConsultaPublicacaoPNCPUseCase) checkCache(ctx context.Context, tipo, valor, dataInicial, dataFinal string) []*pncp.AnaliseResultado {
	cacheParams := map[string]interface{}{
		"tipo": tipo, "valor": valor,
		"dataInicial": dataInicial, "dataFinal": dataFinal,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache(redis.ChaveLicitacoesTrimestre, raw)

	var cached []*pncp.AnaliseResultado
	if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
		logger.New("PNCP: UseCase: buscar").Warn("cache indisponivel", "erro", err)
	} else if ok {
		logger.New("PNCP: UseCase: buscar").Info("cache hit", "tipo", tipo, "valor", valor)
		return cached
	}
	return nil
}

func (u *ConsultaPublicacaoPNCPUseCase) writeCache(ctx context.Context, tipo, valor, dataInicial, dataFinal string, resultados []*pncp.AnaliseResultado) {
	cacheParams := map[string]interface{}{
		"tipo": tipo, "valor": valor,
		"dataInicial": dataInicial, "dataFinal": dataFinal,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache(redis.ChaveLicitacoesTrimestre, raw)

	if err := u.redis.Set(ctx, cacheKey, resultados); err != nil {
		logger.New("PNCP: UseCase: buscar").Warn("cache indisponivel", "erro", err)
	}
}

func (u *ConsultaPublicacaoPNCPUseCase) coletarContratos(ctx context.Context, tipo, valor string, meses []utils.AnoMes) []pncp.Contrato {
	var contratos []pncp.Contrato
	var mesesPendentes []utils.AnoMes

	for _, am := range meses {
		jaRealizada, err := u.repo.BuscaJaRealizada(ctx, tipo, valor, am.Ano, am.Mes)
		if err != nil {
			logger.New("PNCP: UseCase: coletarContratos").Warn("erro ao verificar busca no PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
			mesesPendentes = append(mesesPendentes, am)
			continue
		}
		if jaRealizada {
			persistidos, err := u.repo.BuscarContratosPorFiltro(ctx, tipo, valor, am.Ano, am.Mes)
			if err != nil {
				logger.New("PNCP: UseCase: coletarContratos").Warn("erro ao buscar contratos do PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
				mesesPendentes = append(mesesPendentes, am)
				continue
			}
			for i := range persistidos {
				contratos = append(contratos, repositorios.PersistidoParaContrato(persistidos[i]))
			}
			continue
		}
		mesesPendentes = append(mesesPendentes, am)
	}

	if len(mesesPendentes) > 0 {
		pendentes := buscarMesesParalelo(ctx, tipo, valor, mesesPendentes, u.buscarMesComLock)
		contratos = append(contratos, pendentes...)
	}

	return contratos
}

func (u *ConsultaPublicacaoPNCPUseCase) buscarMesComLock(ctx context.Context, tipo, valor string, ano, mes int) []pncp.Contrato {
	return pncpbusca.BuscarMesComLock(ctx, u.redis, u.repo,
		func(c context.Context, t string, v string, a, m int) ([]pncp.Contrato, error) {
			return u.buscarTodasPaginasDoPNCP(c, t, v, a, m)
		},
		func(c context.Context, t string, v string, a, m int, contratos []pncp.Contrato) error {
			return persistirContratos(c, u.repo, t, v, a, m, contratos)
		},
		tipo, valor, ano, mes,
	)
}

func (u *ConsultaPublicacaoPNCPUseCase) buscarTodasPaginasDoPNCP(ctx context.Context, tipo, valor string, ano, mes int) ([]pncp.Contrato, error) {
	log := logger.New("PNCP: UseCase: buscarTodasPaginasDoPNCP")

	dataInicial, dataFinal := utils.FormatarPeriodoMes(ano, mes)
	pagina := 1
	tamanho := 50
	items := make([]pncp.Contrato, 0)
	seen := make(map[string]struct{})

	for {
		var resp *pncp.PublicacaoResponse
		var err error

		if tipo == "uf" {
			resp, err = u.pncpClient.BuscarContratacoesPorUF(ctx, valor, dataInicial, dataFinal, "", pagina, tamanho)
		} else {
			resp, err = u.pncpClient.BuscarContratacoesPorMunicipio(ctx, valor, dataInicial, dataFinal, "", pagina, tamanho)
		}
		if err != nil {
			return nil, fmt.Errorf("pagina %d: %w", pagina, err)
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
			return nil, ctx.Err()
		default:
		}
	}

	log.Info("contratos obtidos do PNCP publicacao", "tipo", tipo, "valor", valor, "ano", ano, "mes", mes, "total", len(items))
	return items, nil
}

func (u *ConsultaPublicacaoPNCPUseCase) montarResultados(ctx context.Context, contratos []pncp.Contrato, dataInicial, dataFinal string) ([]*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: montarResultados")

	fornecedoresMap := extrairFornecedores(contratos)
	log.Info("fornecedores unicos", "total", len(fornecedoresMap))

	enrichedFornecedores := enriquecerFornecedoresComRepo(ctx, u.repo, u.opencnpjClient, fornecedoresMap)
	contratos = montarContratosDTO(contratos, enrichedFornecedores)

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

type grupoOrgao struct {
	CNPJ        *string
	RazaoSocial *string
	Contratos   []pncp.Contrato
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
