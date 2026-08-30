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

// GetClient returns the underlying Redis client
func (c *Client) GetClient() *redis.Client {
    return c.client
}

// Close closes the Redis connection
func (c *Client) Close() error {
    return c.client.Close()
}

// --------------------
// Redis Key Generators
// --------------------

// Game State Keys
func GameStateKey(gameID string) string {
    return fmt.Sprintf("game:%s:state", gameID)
}

func GamePlayersKey(gameID string) string {
    return fmt.Sprintf("game:%s:players", gameID)
}

func GameDrawnNumbersKey(gameID string) string {
    return fmt.Sprintf("game:%s:drawn", gameID)
}

func GameTakenCardsKey(gameID string) string {
    return fmt.Sprintf("game:%s:cards:taken", gameID)
}

func GameCountdownKey(gameID string) string {
    return fmt.Sprintf("game:%s:countdown", gameID)
}

// GameDrawLeaseKey holds the single-owner lease for a game's draw loop
func GameDrawLeaseKey(gameID string) string {
    return fmt.Sprintf("game:%s:drawlease", gameID)
}

// LobbyActivityKey holds a short-lived marker that a real player recently browsed a tier's lobby
func LobbyActivityKey(tier string) string {
    return fmt.Sprintf("lobby:%s:activity", tier)
}

// Pub/Sub Channels
func GameChannel(gameID string) string {
    return fmt.Sprintf("game:%s:events", gameID)
}

// BonusCampaignChannel carries live "first N players" giveaway events
const BonusCampaignChannel = "admin:bonus_campaign:events"
