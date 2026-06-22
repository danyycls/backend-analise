package dadosfinanceiros

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync/atomic"
	"time"

	"github.com/danyele/podp/internal/shared/clients/ibge"
	siconfiClient "github.com/danyele/podp/internal/shared/clients/siconfi"
	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
)

type BaseFinanceiroUseCase struct {
	siconfiCli      *siconfiClient.SICONFIClient
	ibgeCli         *ibge.IBGEClient
	redisCache      *redis.RedisCache
	apiIndisponivel atomic.Bool
}

func (b *BaseFinanceiroUseCase) SICONFIIndisponivel() bool {
	return b.apiIndisponivel.Load()
}

func NovoBaseFinanceiroUseCase(
	siconfiCli *siconfiClient.SICONFIClient,
	ibgeCli *ibge.IBGEClient,
	redisCache *redis.RedisCache,
) *BaseFinanceiroUseCase {
	return &BaseFinanceiroUseCase{
		siconfiCli: siconfiCli,
		ibgeCli:    ibgeCli,
		redisCache: redisCache,
	}
}

func (b *BaseFinanceiroUseCase) estadoID(ctx context.Context, uf string) (int, error) {
	estados, err := b.ibgeCli.ListarEstados(ctx)
	if err != nil {
		return 0, err
	}
	for _, e := range estados {
		if strings.EqualFold(e.Sigla, uf) {
			return e.ID * 100000, nil
		}
	}
	return 0, nil
}

func (b *BaseFinanceiroUseCase) anoAlvo() int64 {
	return int64(time.Now().Year() - 1)
}

func (b *BaseFinanceiroUseCase) buscarRGF(ctx context.Context, params siconfiClient.RGFParams) ([]siconfiClient.RGFItem, error) {
	log := logger.New("Estadual: Financeiro: RGF")
	raw, _ := json.Marshal(params)
	cacheKey := redis.ChaveCache("estadual-rgf", raw)

	var cached []siconfiClient.RGFItem
	if ok, err := b.redisCache.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	} else if ok {
		return cached, nil
	}

	itens, err := b.siconfiCli.BuscarRGF(ctx, params)
	if err != nil {
		if errors.Is(err, siconfiClient.ErrSICONFIIndisponivel) {
			b.apiIndisponivel.Store(true)
		}
		return nil, err
	}

	if setErr := b.redisCache.Set(ctx, cacheKey, itens); setErr != nil {
		log.Warn("cache indisponivel", "erro", setErr)
	}

	return itens, nil
}

func (b *BaseFinanceiroUseCase) buscarRREO(ctx context.Context, params siconfiClient.RREOParams) ([]siconfiClient.RREOItem, error) {
	log := logger.New("Estadual: Financeiro: RREO")
	raw, _ := json.Marshal(params)
	cacheKey := redis.ChaveCache("estadual-rreo", raw)

	var cached []siconfiClient.RREOItem
	if ok, err := b.redisCache.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	} else if ok {
		return cached, nil
	}

	itens, err := b.siconfiCli.BuscarRREO(ctx, params)
	if err != nil {
		if errors.Is(err, siconfiClient.ErrSICONFIIndisponivel) {
			b.apiIndisponivel.Store(true)
		}
		return nil, err
	}

	if setErr := b.redisCache.Set(ctx, cacheKey, itens); setErr != nil {
		log.Warn("cache indisponivel", "erro", setErr)
	}

	return itens, nil
}
