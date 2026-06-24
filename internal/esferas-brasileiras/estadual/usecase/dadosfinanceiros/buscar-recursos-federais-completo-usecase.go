package dadosfinanceiros

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"time"

	portalClient "github.com/danyele/podp/internal/shared/clients/portaltransparencia"
	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/types"
)

const maxPaginas = 20

type EsferaEstadualBuscarRecursosFederaisCompletoRequest struct {
	UF         string
	Exercicio  int64
	CodigoIBGE string
}

type EsferaEstadualBuscarRecursosFederaisCompletoResponse struct {
	Dados []types.RecursoFederalRecebido `json:"dados"`
	Total int                            `json:"total"`
}

type EsferaEstadualBuscarRecursosFederaisCompletoUseCase struct {
	portalCli  *portalClient.PortalTransparenciaClient
	redisCache *redis.RedisCache
}

func NovoEsferaEstadualBuscarRecursosFederaisCompletoUseCase(
	portalCli *portalClient.PortalTransparenciaClient,
	redisCache *redis.RedisCache,
) *EsferaEstadualBuscarRecursosFederaisCompletoUseCase {
	return &EsferaEstadualBuscarRecursosFederaisCompletoUseCase{
		portalCli:  portalCli,
		redisCache: redisCache,
	}
}

func (u *EsferaEstadualBuscarRecursosFederaisCompletoUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarRecursosFederaisCompletoRequest) (*EsferaEstadualBuscarRecursosFederaisCompletoResponse, error) {
	resultado := u.buscarRecursosFederaisCompleto(ctx, req.UF, req.Exercicio, req.CodigoIBGE)
	return &EsferaEstadualBuscarRecursosFederaisCompletoResponse{
		Dados: resultado,
		Total: len(resultado),
	}, nil
}

func (u *EsferaEstadualBuscarRecursosFederaisCompletoUseCase) buscarRecursosFederaisCompleto(ctx context.Context, uf string, exercicio int64, codigoIBGE string) []types.RecursoFederalRecebido {
	log := logger.New("Estadual: UseCase: BuscarRecursosFederaisCompleto")
	anoAlvo := exercicio
	if anoAlvo <= 0 {
		anoAlvo = int64(time.Now().Year() - 1)
	}

	filtroBase := portalClient.DespesaRecursosRecebidosQueryParams{
		MesAnoInicio: strconv.Itoa(int(anoAlvo)) + "-01",
		MesAnoFim:    strconv.Itoa(int(anoAlvo)) + "-12",
		UF:           uf,
		CodigoIBGE:   codigoIBGE,
	}

	var todas []types.RecursoFederalRecebido

	for pagina := 1; pagina <= maxPaginas; pagina++ {
		filtro := filtroBase
		filtro.Pagina = pagina

		cached, hit := u.tentarCache(ctx, filtro)
		if hit {
			if len(cached) == 0 {
				break
			}
			todas = append(todas, cached...)
			continue
		}

		itens, err := u.portalCli.ListarRecursosRecebidos(ctx, filtro)
		if err != nil {
			log.Error("erro ao buscar recursos recebidos completos", "uf", uf, "pagina", pagina, "erro", err)
			break
		}

		if len(itens) == 0 {
			u.gravarCache(ctx, filtro, nil)
			break
		}

		result := u.montarResultado(itens)
		u.gravarCache(ctx, filtro, result)
		todas = append(todas, result...)
	}

	sort.Slice(todas, func(i, j int) bool {
		return todas[i].Valor > todas[j].Valor
	})

	return todas
}

func (u *EsferaEstadualBuscarRecursosFederaisCompletoUseCase) tentarCache(ctx context.Context, filtro portalClient.DespesaRecursosRecebidosQueryParams) ([]types.RecursoFederalRecebido, bool) {
	log := logger.New("Estadual: UseCase: BuscarRecursosFederaisCompleto")
	raw, _ := json.Marshal(filtro)
	cacheKey := redis.ChaveCache("estadual-recursos-federais-completo", raw)

	var cached []types.RecursoFederalRecebido
	if ok, err := u.redisCache.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
		return nil, false
	} else if ok {
		return cached, true
	}
	return nil, false
}

func (u *EsferaEstadualBuscarRecursosFederaisCompletoUseCase) montarResultado(itens []portalClient.PessoaRecursosRecebidosUGMesDesnormalizada) []types.RecursoFederalRecebido {
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
	return result
}

func (u *EsferaEstadualBuscarRecursosFederaisCompletoUseCase) gravarCache(ctx context.Context, filtro portalClient.DespesaRecursosRecebidosQueryParams, dados []types.RecursoFederalRecebido) {
	log := logger.New("Estadual: UseCase: BuscarRecursosFederaisCompleto")
	raw, _ := json.Marshal(filtro)
	cacheKey := redis.ChaveCache("estadual-recursos-federais-completo", raw)
	if err := u.redisCache.Set(ctx, cacheKey, dados); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}
}
