package redis

import (
	"auth-service/pkg/env"
	"auth-service/pkg/logs"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	PREFIX_ACTIVE = "session"
)

var (
	disableRedis        = false
	client              *redis.Client
	clientOnce          sync.Once
	ctx                 = context.Background()
	isConnected         = false
	Nil                 = redis.Nil
	retryAttempts       = 3
	retryDelay          = 5 * time.Second
	healthCheckInterval = 5 * time.Second
)

type CacheResponse struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func InitializeRedis() {
	const FUNC_NAME = "InitializeRedisClient"

	if disableRedis {
		return
	}
	redisAddr := env.REDIS_URL
	if redisAddr == "" {
		logs.Warn(FUNC_NAME, "REDIS_URL environment variable is not set", redisAddr)
	}

	clientOnce.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})

		go attemptConnection(redisAddr)
		go monitorRedisHealth(redisAddr)
	})
}

func GetClient() *redis.Client {
	if client == nil {
		InitializeRedis()
	}
	return client
}

func Close() error {
	return client.Close()
}

func IsClientConnected() bool {
	return isConnected
}

// Possible race conditions-may need revision in the future
func Set(key string, value any, ttl time.Duration) error {
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		var err error
		b, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}
	return client.Set(ctx, key, b, ttl).Err()
}

func Get[T any](key string, dest *T) error {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		logs.Error("", "error.getting.session", err)
		return err
	}

	if err := json.Unmarshal([]byte(result), dest); err != nil {
		logs.Error("", "error.unmarshalling.session", err)
		return err
	}
	return nil
}

func GetString(key string) (string, error) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func attemptConnection(redisAddr string) {
	const FUNC_NAME = "attemptConnection"

	for range retryAttempts {
		logs.Info(FUNC_NAME, "Attempting to ping redis client at", redisAddr)
		_, err := client.Ping(ctx).Result()
		if err != nil {
			logs.Error(FUNC_NAME, "Failed to connect to Redis", err)
			time.Sleep(retryDelay)
		} else {
			logs.Info("Redis status: Connected", redisAddr)
			isConnected = true
			break
		}
	}

	if !isConnected {
		logs.Warn("InitializeClient", "Could not connect to Redis", retryAttempts)
	}
}

func monitorRedisHealth(redisAddr string) {
	const FUNC_NAME = "monitorRedisHealth"

	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		if isConnected {
			_, err := client.Ping(ctx).Result()
			if err != nil {
				logs.Error(FUNC_NAME, "Redis connection lost", err)
				isConnected = false
				go attemptConnection(redisAddr)
			}
		}
	}
}
