package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/utils"
)

type ConsultaContratoUFMunicipioPNCPUseCase struct {
	PncpUseCaseBase
}

func NovoConsultaContratoUFMunicipioPNCPUseCase(pncp *pncp.PNCPClient, opencnpj *opencnpj.OpenCNPJClient, redis *redis.RedisCache, repo repositorios.PNCPRepository) *ConsultaContratoUFMunicipioPNCPUseCase {
	return &ConsultaContratoUFMunicipioPNCPUseCase{
		PncpUseCaseBase: NewPncpUseCaseBase(pncp, opencnpj, redis, repo),
	}
}

func (u *ConsultaContratoUFMunicipioPNCPUseCase) BuscarPorUF(ctx context.Context, uf, dataInicial, dataFinal string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	return u.executar(ctx, "uf", uf, dataInicial, dataFinal, paginasErro)
}

func (u *ConsultaContratoUFMunicipioPNCPUseCase) BuscarPorMunicipio(ctx context.Context, codigoMunicipio, dataInicial, dataFinal string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	return u.executar(ctx, "municipio", codigoMunicipio, dataInicial, dataFinal, paginasErro)
}

func (u *ConsultaContratoUFMunicipioPNCPUseCase) executar(ctx context.Context, tipo, valor, dataInicial, dataFinal string, paginasErro *[]int) ([]*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: ContratoUFMunicipio")

	if result := checkCache(ctx, u.redis, tipo, valor, dataInicial, dataFinal); result != nil {
		return result, nil
	}

	meses := utils.ExtrairMeses(dataInicial, dataFinal)
	if len(meses) == 0 {
		return nil, fmt.Errorf("periodo invalido: %s a %s", dataInicial, dataFinal)
	}

	contratos, mesesPendentes := consultarContratosExistentes(ctx, u.repo, tipo, valor, meses)
	if len(contratos) == 0 && len(mesesPendentes) == 0 {
		return nil, fmt.Errorf("erro ao consultar repositorio para os contratos %s %s", tipo, valor)
	}

	if len(mesesPendentes) > 0 {
		pendentes := buscarMesesParalelo(ctx, tipo, valor, mesesPendentes, u.fetchContratosPorMes)
		contratos = append(contratos, pendentes...)
	}

	log.Info("total de itens", "total", len(contratos))

	resultados, err := montarResultadosAgrupados(ctx, u.repo, u.opencnpjClient, contratos, tipo, valor, dataInicial, dataFinal)
	if err != nil {
		return nil, err
	}

	for _, am := range mesesPendentes {
		for _, r := range resultados {
			contratosDoMes := filtrarContratosPorMes(r.Contratos, am.Ano, am.Mes)
			if len(contratosDoMes) > 0 {
				if err := persistirContratos(ctx, u.repo, tipo, valor, am.Ano, am.Mes, contratosDoMes); err != nil {
					log.Warn("erro ao persistir contratos", "ano", am.Ano, "mes", am.Mes, "erro", err)
				}
			}
		}
	}

	writeCache(ctx, u.redis, tipo, valor, dataInicial, dataFinal, resultados)
	return resultados, nil
}

func (u *ConsultaContratoUFMunicipioPNCPUseCase) fetchContratosPorMes(ctx context.Context, tipo, valor string, ano, mes int) []pncp.Contrato {
	log := logger.New("PNCP: UseCase: fetchContratosPorMes")

	var fetchPagina FetchPaginaContratos
	if tipo == "uf" {
		fetchPagina = func(ctx context.Context, valor, dataInicial, dataFinal string, pagina, tamanho int) (*pncp.ContratoResponse, error) {
			return u.pncpClient.BuscarContratosPorUF(ctx, valor, dataInicial, dataFinal, "", pagina, tamanho)
		}
	} else {
		fetchPagina = func(ctx context.Context, valor, dataInicial, dataFinal string, pagina, tamanho int) (*pncp.ContratoResponse, error) {
			return u.pncpClient.BuscarContratosPorMunicipio(ctx, valor, dataInicial, dataFinal, "", pagina, tamanho)
		}
	}

	contratos, err := buscarContratosPaginado(ctx, fetchPagina, valor, ano, mes, 50, 3*time.Second)
	if err != nil {
		log.Error("erro ao buscar contratos do PNCP", "tipo", tipo, "valor", valor, "ano", ano, "mes", mes, "erro", err)
		return nil
	}

	return contratos
}
