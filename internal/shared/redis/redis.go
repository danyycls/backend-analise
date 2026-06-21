package redis

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const prefixoChave = "podp:cache:"
const ttlPadrao = 30 * 24 * time.Hour

type RedisCache struct {
	client *goredis.Client
}

func NovoRedisCache() *RedisCache {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")
	db := 0
	if s := os.Getenv("REDIS_DB"); s != "" {
		db, _ = strconv.Atoi(s)
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		MaxRetries:      1,
		MinRetryBackoff: 500 * time.Millisecond,
	})

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
