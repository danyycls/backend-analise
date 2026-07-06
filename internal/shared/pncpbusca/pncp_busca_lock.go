package pncpbusca

import (
	"context"
	"time"

	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/repositorios"
)

type FetchFn func(ctx context.Context, tipo, valor string, ano, mes int) ([]pncp.Contrato, error)
type PersistFn func(ctx context.Context, tipo, valor string, ano, mes int, contratos []pncp.Contrato) error

func BuscarMesComLock(
	ctx context.Context,
	redisCli *redis.RedisCache,
	repo repositorios.PNCPRepository,
	fetch FetchFn,
	persist PersistFn,
	tipo, valor string,
	ano, mes int,
) []pncp.Contrato {
	log := logger.New("PNCP: BuscarMesComLock")

	adquirido, err := redisCli.AdquirirLock(ctx, tipo, valor, ano, mes)
	if err != nil {
		log.Warn("erro ao adquirir lock", "ano", ano, "mes", mes, "erro", err)
		return nil
	}

	if adquirido {
		defer func() {
			if err := redisCli.LiberarLock(ctx, tipo, valor, ano, mes); err != nil {
				log.Warn("erro ao liberar lock", "ano", ano, "mes", mes, "erro", err)
			}
		}()

		contratos, err := fetch(ctx, tipo, valor, ano, mes)
		if err != nil {
			log.Error("erro ao buscar contratos do PNCP", "ano", ano, "mes", mes, "erro", err)
			return nil
		}

		if err := persist(ctx, tipo, valor, ano, mes, contratos); err != nil {
			log.Error("erro ao persistir contratos", "ano", ano, "mes", mes, "erro", err)
		}

		return contratos
	}

	for tentativa := 1; tentativa <= 5; tentativa++ {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(2 * time.Second):
		}

		jaRealizada, err := repo.BuscaJaRealizada(ctx, tipo, valor, ano, mes)
		if err != nil {
			log.Warn("erro ao verificar busca", "tentativa", tentativa, "erro", err)
			continue
		}
		if jaRealizada {
			persistidos, err := repo.BuscarContratosPorFiltro(ctx, tipo, valor, ano, mes)
			if err != nil {
				log.Warn("erro ao buscar contratos persistidos", "tentativa", tentativa, "erro", err)
				continue
			}
			contratos := make([]pncp.Contrato, len(persistidos))
			for i := range persistidos {
				contratos[i] = repositorios.PersistidoParaContrato(persistidos[i])
			}
			return contratos
		}

		adquirido, err := redisCli.AdquirirLock(ctx, tipo, valor, ano, mes)
		if err != nil {
			continue
		}
		if adquirido {
			defer func() {
				if err := redisCli.LiberarLock(ctx, tipo, valor, ano, mes); err != nil {
					log.Warn("erro ao liberar lock apos espera", "erro", err)
				}
			}()

			contratos, err := fetch(ctx, tipo, valor, ano, mes)
			if err != nil {
				log.Error("erro ao buscar contratos apos espera", "ano", ano, "mes", mes, "erro", err)
				return nil
			}

			if err := persist(ctx, tipo, valor, ano, mes, contratos); err != nil {
				log.Error("erro ao persistir contratos apos espera", "erro", err)
			}

			return contratos
		}
	}

	log.Warn("max tentativas excedido para lock", "tipo", tipo, "valor", valor, "ano", ano, "mes", mes)
	return nil
}
