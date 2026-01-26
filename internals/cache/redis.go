package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string, password string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisStore{client: rdb}
}

// AtomicBook attempts to lock a resource using a Lua Script.
// We make the key generic so it can be used for things other than just seats if needed.
func (r *RedisStore) AtomicBook(ctx context.Context, keyName string, value interface{}, expirySeconds int) error {
	// --- THE LUA SCRIPT ---
	script := `
		if redis.call("EXISTS", KEYS[1]) == 1 then
			return 0
		end
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
		return 1
	`

	// Execute
	result, err := r.client.Eval(ctx, script, []string{keyName}, value, expirySeconds).Int()

	if err != nil {
		return fmt.Errorf("redis execution failed: %w", err)
	}

	if result == 0 {
		return errors.New("resource locked by another process")
	}

	return nil
}
