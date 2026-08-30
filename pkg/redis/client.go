package redis

import (
    "context"
    "crypto/tls"
    "fmt"
    "time"

    "github.com/bingo/backend/config"
    "github.com/redis/go-redis/v9"
)

type Client struct {
    client *redis.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:      cfg.Redis.GetAddr(), // e.g. "uncommon-slug-85330.upstash.io:6379"
        Username:  "default",           // Upstash always uses "default"
        Password:  cfg.Redis.Password,  // your Upstash password
        DB:        cfg.Redis.DB,
        TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12}, // required for Upstash
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return &Client{client: rdb}, nil
}
