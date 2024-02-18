package cache

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func Init() {
	var (
		REDIS_HOST     = helpers.GoDotEnvVariable("REDIS_HOST")
		REDIS_PORT     = helpers.GoDotEnvVariable("REDIS_PORT")
		REDIS_PASSWORD = helpers.GoDotEnvVariable("REDIS_PASSWORD")
	)
	if REDIS_HOST == "" {
		logger.Logger.Error("cache.Init: ERROR missing redis env vars.")
		return
	}
	rdb = redis.NewClient(&redis.Options{
		Password: REDIS_PASSWORD,
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		DB:       0, // use default db
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		logger.Logger.Error("cache.Init: Error connecting to redis:", err)
		rdb = nil
		return
	}
	logger.Logger.Log("cache.Init: Redis cache connected")
}

/*
Duration in minutes
Set 0 if you want the data to not expire

Ex: set(key, "test", 15) // Sets the expiry to 15 minutes
*/
func Set(key string, value interface{}, duration time.Duration) error {
	if rdb == nil {
		logger.Logger.Log("Redis client not connected")
		return nil
	}
	err := rdb.Set(ctx, key, value, duration*time.Minute).Err()
	if err != nil {
		logger.Logger.Error("Unable to set cache:", err)
	}
	return err
}

func Get(key string) (string, error) {
	if rdb == nil {
		logger.Logger.Log("Redis client not connected")
		return "", nil
	}
	return rdb.Get(ctx, key).Result()
}

// Redis accepts multiple keys to be deleted
func Del(keys ...string) {
	if rdb == nil {
		logger.Logger.Log("Redis client not connected")
	}
	rdb.Del(ctx, keys...)
}

func Expire(key string, duration time.Duration) {
	if rdb == nil {
		logger.Logger.Log("Redis client not connected")
		return
	}
	rdb.Expire(ctx, key, duration)
}

/*
This method first checks cache and fetchs data if cache miss and updates cache.
Set duration to `0` to not set expiration time.

Make sure to pass the address of the result var `&result` to get data returned from this func.
@Usage Fetch(key, duration, &target, cb)
*/
func Fetch(key string, duration time.Duration, target any, callback func() (interface{}, error)) error {
	data, _ := Get(key)
	if data != "" {
		logger.Logger.Log("cache.Fetch: Cache hit for:", key)
		if err := json.Unmarshal([]byte(data), target); err != nil {
			logger.Logger.Log("cache.Fetch: Fatal ERROR", err)
			return errors.New("unable to fetch data from cache")
		}
		return nil
	}
	logger.Logger.Log("cache.Fetch: Cache miss for:", key)
	result, err := callback()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		logger.Logger.Log("Unable to set cache", err)
	} else {
		Set(key, string(jsonData), duration)
	}
	json.Unmarshal(jsonData, target)
	return nil
}
