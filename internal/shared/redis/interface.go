package redis

//go:generate mockgen -source=interface.go -destination=mock.go -package=redis

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, chave string, destino interface{}) (bool, error)
	Set(ctx context.Context, chave string, valor interface{}) error
	SetEx(ctx context.Context, chave string, valor interface{}, ttl time.Duration) error
	SAdd(ctx context.Context, chave string, membros ...string) (int64, error)
	SMembers(ctx context.Context, chave string) ([]string, error)
	SInter(ctx context.Context, chaves ...string) ([]string, error)
	SUnionStore(ctx context.Context, destino string, chaves ...string) (int64, error)
	Exists(ctx context.Context, chaves ...string) (int64, error)
	Del(ctx context.Context, chaves ...string) (int64, error)
}

var _ Cache = (*RedisCache)(nil)
