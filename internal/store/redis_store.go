package store

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb      *redis.Client
	cacheTTL time.Duration
}

func NewRedisStore() (*RedisStore, error) {
	addr := getenv("REDIS_ADDR", "redis:6379")
	dbStr := getenv("REDIS_DB", "0")
	db, _ := strconv.Atoi(dbStr)
	ttlSec, _ := strconv.Atoi(getenv("CACHE_TTL_SECONDS", "3600"))

	rdb := redis.NewClient(&redis.Options{Addr: addr, DB: db})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return &RedisStore{rdb: rdb, cacheTTL: time.Duration(ttlSec) * time.Second}, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

const hashKey = "links" // 永続化: HSET links <code> <url>

func (s *RedisStore) Get(ctx context.Context, code string) (string, bool, error) {
	// 1) キャッシュ
	key := "cache:" + code
	if v, err := s.rdb.Get(ctx, key).Result(); err == nil && v != "" {
		return v, true, nil
	}
	// 2) 永続
	v, err := s.rdb.HGet(ctx, hashKey, code).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	// 3) キャッシュ書き込み
	_ = s.rdb.Set(ctx, key, v, s.cacheTTL).Err()
	return v, false, nil
}

func (s *RedisStore) Exists(ctx context.Context, code string) (bool, error) {
	return s.rdb.HExists(ctx, hashKey, code).Result()
}

func (s *RedisStore) Save(ctx context.Context, code, url string) error {
	if err := s.rdb.HSet(ctx, hashKey, code, url).Err(); err != nil {
		return err
	}
	_ = s.rdb.Del(ctx, "cache:"+code).Err()
	return nil
}
