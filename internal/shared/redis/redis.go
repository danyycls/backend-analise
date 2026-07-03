package redis

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const prefixoChave = "podp:cache:"
const ttlPadrao = 30 * 24 * time.Hour

const (
	ChaveLicitacoesTrimestre              = "licitacoes-trimestre"
	ChavePNCPublicacaoPagina              = "pncp-publicacao-pagina"
	ChavePNCBuscarContratos               = "pncp-buscar-contratos"
	ChaveOrgaoAnalise                     = "orgao-analise"
	ChavePublicacaoAnalise                = "publicacao-analise"
	ChaveLigacaoPolitica                  = "ligacao-politica"
	ChaveDeputados                        = "deputados"
	ChaveSenadoSenadores                  = "senado-senadores"
	ChaveEstadualRGF                      = "estadual-rgf"
	ChaveEstadualRREO                     = "estadual-rreo"
	ChaveEstadualRecursosFederais         = "estadual-recursos-federais"
	ChaveEstadualRecursosFederaisCompleto = "estadual-recursos-federais-completo"
)

type RedisCache struct {
	client *goredis.Client
}

func NovoRedisCache() *RedisCache {
	opt, err := goredis.ParseURL(os.Getenv("REDIS_ADDR"))
if err != nil {
    panic(err)
}

opt.MaxRetries = 1
opt.MinRetryBackoff = 500 * time.Millisecond

client := goredis.NewClient(opt)

	return &RedisCache{client: client}
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func ChaveCache(endpoint string, params []byte) string {
	h := md5.Sum(params)
	return prefixoChave + endpoint + ":" + hex.EncodeToString(h[:])
}

func (r *RedisCache) Get(ctx context.Context, chave string, destino interface{}) (bool, error) {
	data, err := r.client.Get(ctx, chave).Bytes()
	if errors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(data, destino); err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisCache) Set(ctx context.Context, chave string, valor interface{}) error {
	data, err := json.Marshal(valor)
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}
	return r.client.Set(ctx, chave, data, ttlPadrao).Err()
}

func (r *RedisCache) SetEx(ctx context.Context, chave string, valor interface{}, ttl time.Duration) error {
	data, err := json.Marshal(valor)
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}
	return r.client.Set(ctx, chave, data, ttl).Err()
}

func (r *RedisCache) SAdd(ctx context.Context, chave string, membros ...string) (int64, error) {
	args := make([]interface{}, len(membros))
	for i, m := range membros {
		args[i] = m
	}
	return r.client.SAdd(ctx, chave, args...).Result()
}

func (r *RedisCache) SMembers(ctx context.Context, chave string) ([]string, error) {
	return r.client.SMembers(ctx, chave).Result()
}

func (r *RedisCache) SInter(ctx context.Context, chaves ...string) ([]string, error) {
	return r.client.SInter(ctx, chaves...).Result()
}

func (r *RedisCache) SUnionStore(ctx context.Context, destino string, chaves ...string) (int64, error) {
	return r.client.SUnionStore(ctx, destino, chaves...).Result()
}

func (r *RedisCache) Exists(ctx context.Context, chaves ...string) (int64, error) {
	return r.client.Exists(ctx, chaves...).Result()
}

func (r *RedisCache) Del(ctx context.Context, chaves ...string) (int64, error) {
	return r.client.Del(ctx, chaves...).Result()
}
