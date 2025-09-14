package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	rdb    *redis.Client
	ctx    context.Context
	limit  int           // max requests
	window time.Duration // time window
}

func NewRateLimiter(rdb *redis.Client, ctx context.Context, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		rdb:    rdb,
		ctx:    ctx,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(ip string) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", ip)

	count, err := rl.rdb.Incr(rl.ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rl.rdb.Expire(rl.ctx, key, rl.window)
	}

	if count > int64(rl.limit) {
		return false, nil
	}

	return true, nil
}
