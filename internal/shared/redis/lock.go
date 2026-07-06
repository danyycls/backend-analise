package redis

import (
	"context"
	"fmt"
	"time"
)

const ttlLock = 30 * time.Second

func chaveLock(tipo, valor string, ano, mes int) string {
	return fmt.Sprintf("podp:lock:pncp:%s:%s:%04d:%02d", tipo, valor, ano, mes)
}

func (r *RedisCache) AdquirirLock(ctx context.Context, tipo, valor string, ano, mes int) (bool, error) {
	return r.client.SetNX(ctx, chaveLock(tipo, valor, ano, mes), "1", ttlLock).Result()
}

func (r *RedisCache) LiberarLock(ctx context.Context, tipo, valor string, ano, mes int) error {
	return r.client.Del(ctx, chaveLock(tipo, valor, ano, mes)).Err()
}
