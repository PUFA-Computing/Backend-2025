package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	// "os"
	"time"
)

var Rdb *redis.Client

// RedisEnabled indicates if Redis is available and connected
var RedisEnabled bool = false

func InitRedis() {
	// Use default local Redis connection
	options := &redis.Options{
		Addr: "localhost:6379", // Default Redis address
		DB: 0,
		DialTimeout: 5 * time.Second,
	}

	log.Println("Attempting to connect to Redis...")
	Rdb = redis.NewClient(options)
	
	ctx := context.Background()
	if err := Rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	
	RedisEnabled = true
	log.Println("Successfully connected to Redis")
}

func IsTokenRevoked(tokenString string) (bool, error) {
	if !RedisEnabled || Rdb == nil {
		// If Redis is not available, assume token is not revoked
		return false, nil
	}
	
	ctx := context.Background()
	exists, err := Rdb.SIsMember(ctx, "revoked_tokens", tokenString).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func RevokeToken(tokenString string) error {
	if !RedisEnabled || Rdb == nil {
		// If Redis is not available, just log and return success
		log.Println("WARNING: Redis not available, token revocation not persisted")
		return nil
	}
	
	ctx := context.Background()
	_, err := Rdb.SAdd(ctx, "revoked_tokens", tokenString).Result()
	if err != nil {
		return err
	}

	return nil
}

func CloseRedis() {
	if !RedisEnabled || Rdb == nil {
		return
	}
	
	err := Rdb.Close()
	if err != nil {
		log.Printf("Error closing Redis connection: %v", err)
		return
	}
	RedisEnabled = false
}
