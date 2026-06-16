package redis

//go:generate mockgen -source=interface.go -destination=mock.go -package=redis

import "context"

type Cache interface {
	Get(ctx context.Context, chave string, destino interface{}) (bool, error)
	Set(ctx context.Context, chave string, valor interface{}) error
}

var _ Cache = (*RedisCache)(nil)
