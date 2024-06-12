package redis

import (
	"auth-service/pkg/utils"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	PREFIX_ACTIVE      = "session:active"
	PREFIX_BLACKLISTED = "session:blacklisted"
)

var (
	cachePeriodBlacklist = 2 * time.Hour
	cachePeriod          = 15 * time.Minute
	client               *redis.Client
	clientOnce           sync.Once
	ctx                  = context.Background()
	Nil                  = redis.Nil
)

func InitializeClient() {
	redisAddr := os.Getenv("REDIS_ADDRESS")

	if redisAddr == "" {
		utils.LogWarn("InitializeClient", "REDIS_ADDRESS environment variable is not set", redisAddr)
	}

	clientOnce.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		utils.LogError("InitializeClient", "Failed to connect to Redis at %s: %v", err)
	} else {
		utils.LogInfo("InitializeClient", "Successfully connected to Redis at %s", redisAddr)
	}
}
func GetClient() *redis.Client {
	if client == nil {
		InitializeClient()
	}
	return client
}
func StartNewSession(token string) error {
	key := fmt.Sprintf("%s:%s", PREFIX_ACTIVE, token)
	if err := client.Set(ctx, key, token, cachePeriod).Err(); err != nil {
		utils.LogError("InitializeClient", "Failed to new session", err)
		return err
	}

	return nil
}
func InvalidateActiveSession(token string) error {
	key := fmt.Sprintf("%s:%s", PREFIX_ACTIVE, token)
	var cursor uint64
	var err error

	for {
		var keys []string
		keys, cursor, err = client.Scan(ctx, cursor, key, 10).Result()
		if err != nil {
			utils.LogError("InvalidateActiveSession", "error scanning keys for token", err)
			return err
		}

		if len(keys) > 0 {
			err = client.Del(ctx, keys...).Err()
			if err != nil {
				utils.LogError("InvalidateActiveSession", "error deleting keys for token", err)
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
func BlacklistSession(token string) error {
	key := fmt.Sprintf("%s:%s", PREFIX_BLACKLISTED, token)
	if err := client.Set(ctx, key, true, cachePeriodBlacklist).Err(); err != nil {
		utils.LogError("BlacklistSession", "error blacklisting session", err)
		return err
	}

	return nil
}

//	func IsBlacklisted(token string) (bool, error) {
//		key := fmt.Sprintf("blacklist:%s", token)
//		err := client.Get(ctx, key).Err()
//		if err == redis.Nil {
//			return false, nil
//		} else if err != nil {
//			return false, err
//		}
//		return true, nil
//	}
// func GetAllKeys() []string {
// 	var cursor uint64
// 	var keys []string

// 	for {
// 		var newKeys []string
// 		newKeys, cursor, err := client.Scan(ctx, cursor, "*", 10).Result()
// 		if err != nil {
// 			return nil
// 		}

// 		keys = append(keys, newKeys...)

// 		if cursor == 0 {
// 			break
// 		}
// 	}

// 	return keys
// }
