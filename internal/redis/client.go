package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	cache2HR   = 2 * time.Hour
	client     *redis.Client
	clientOnce sync.Once
	ctx        = context.Background()
	Nil        = redis.Nil
)

func InitializeClient() {
	clientOnce.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	})
}
func DeleteSession(token string) error {
	pattern := fmt.Sprintf("*:%s", token)
	var cursor uint64
	var err error

	for {
		var keys []string
		keys, cursor, err = client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return fmt.Errorf("error scanning keys for token: %w", err)
		}

		if len(keys) > 0 {
			err = client.Del(ctx, keys...).Err()
			if err != nil {
				return fmt.Errorf("error deleting keys for token: %w", err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
func BlacklistSession(token string) error {
	key := fmt.Sprintf("blacklist:%s", token)
	err := client.Set(ctx, key, true, cache2HR).Err()
	if err != nil {
		return err
	}
	return nil
}
func IsBlacklisted(token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	err := client.Get(ctx, key).Err()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
func GetAllKeys() []string {
	var cursor uint64
	var keys []string

	for {
		var newKeys []string
		newKeys, cursor, err := client.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return nil
		}

		keys = append(keys, newKeys...)

		if cursor == 0 {
			break
		}
	}

	return keys
}
