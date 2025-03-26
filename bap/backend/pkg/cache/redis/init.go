package redis

import (
    "context"
    "errors"
    "time"

    "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
    "github.com/redis/go-redis/v9"
)

// RedisKeys constants for different types of data
const (
    SearchResponseKey = "search-job-response:"
)

// RedisClient structure contains the Redis client instance
type RedisClient struct {
    Client    *redis.Client
    Context   context.Context
}

var (
    Client            *RedisClient
    ConnectionTimeout = 20 * time.Second
    backgroundContext = context.Background()
)

// NewRedisClient initializes Redis connection
func NewRedisClient() (*RedisClient, error) {
    client, err := connect()
    if err != nil {
        return nil, err
    }

    return &RedisClient{
        Client:  client,
        Context: backgroundContext,
    }, nil
}

func connect() (*redis.Client, error) {
    if config.Config.RedisHost == "" || config.Config.RedisPort == "" {
        return nil, errors.New("invalid redis configuration")
    }

    rdb := redis.NewClient(&redis.Options{
        Addr:     config.Config.RedisHost + ":" + config.Config.RedisPort,
        Password: config.Config.RedisPassword, // empty string if no password
        DB:       0,                          // default DB
    })

    ctx, cancel := context.WithTimeout(backgroundContext, ConnectionTimeout)
    defer cancel()

    // Check the connection
    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, err
    }

    return rdb, nil
}