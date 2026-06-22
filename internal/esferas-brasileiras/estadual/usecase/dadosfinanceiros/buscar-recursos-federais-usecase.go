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

type EsferaEstadualBuscarRecursosFederaisRequest struct {
	UF        string
	Exercicio int64
}

type EsferaEstadualBuscarRecursosFederaisResponse struct {
	Dados []types.RecursoFederalRecebido
}

type EsferaEstadualBuscarRecursosFederaisUseCase struct {
	portalCli  *portalClient.PortalTransparenciaClient
	redisCache *redis.RedisCache
}

func NovoEsferaEstadualBuscarRecursosFederaisUseCase(
	portalCli *portalClient.PortalTransparenciaClient,
	redisCache *redis.RedisCache,
) *EsferaEstadualBuscarRecursosFederaisUseCase {
	return &EsferaEstadualBuscarRecursosFederaisUseCase{
		portalCli:  portalCli,
		redisCache: redisCache,
	}
}

func (u *EsferaEstadualBuscarRecursosFederaisUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarRecursosFederaisRequest) (*EsferaEstadualBuscarRecursosFederaisResponse, error) {
	resultado := u.buscarRecursosFederais(ctx, req.UF, req.Exercicio)
	return &EsferaEstadualBuscarRecursosFederaisResponse{Dados: resultado}, nil
}

func (u *EsferaEstadualBuscarRecursosFederaisUseCase) buscarRecursosFederais(ctx context.Context, uf string, exercicio int64) []types.RecursoFederalRecebido {
	log := logger.New("Estadual: UseCase: BuscarRecursosFederais")
	anoAlvo := exercicio
	if anoAlvo <= 0 {
		anoAlvo = int64(time.Now().Year() - 1)
	}

	filtro := portalClient.DespesaRecursosRecebidosQueryParams{
		Pagina:       1,
		MesAnoInicio: strconv.Itoa(int(anoAlvo)) + "-01",
		MesAnoFim:    strconv.Itoa(int(anoAlvo)) + "-12",
		UF:           uf,
	}

	cached, hit := u.tentarCache(ctx, filtro)
	if hit {
		return cached
	}

	itens, err := u.portalCli.ListarRecursosRecebidos(ctx, filtro)
	if err != nil {
		log.Error("erro ao buscar recursos recebidos para estado", "uf", uf, "erro", err)
		return nil
	}

	result := u.montarResultado(itens)
	u.gravarCache(ctx, filtro, result)

	return result
}

func (u *EsferaEstadualBuscarRecursosFederaisUseCase) tentarCache(ctx context.Context, filtro portalClient.DespesaRecursosRecebidosQueryParams) ([]types.RecursoFederalRecebido, bool) {
	log := logger.New("Estadual: UseCase: BuscarRecursosFederais")
	raw, _ := json.Marshal(filtro)
	cacheKey := redis.ChaveCache("estadual-recursos-federais", raw)

	var cached []types.RecursoFederalRecebido
	if ok, err := u.redisCache.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
		return nil, false
	} else if ok {
		return cached, true
	}
	return nil, false
}

func (u *EsferaEstadualBuscarRecursosFederaisUseCase) montarResultado(itens []portalClient.PessoaRecursosRecebidosUGMesDesnormalizada) []types.RecursoFederalRecebido {
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
	return result
}

func (u *EsferaEstadualBuscarRecursosFederaisUseCase) gravarCache(ctx context.Context, filtro portalClient.DespesaRecursosRecebidosQueryParams, dados []types.RecursoFederalRecebido) {
	log := logger.New("Estadual: UseCase: BuscarRecursosFederais")
	raw, _ := json.Marshal(filtro)
	cacheKey := redis.ChaveCache("estadual-recursos-federais", raw)
	if err := u.redisCache.Set(ctx, cacheKey, dados); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}
}
