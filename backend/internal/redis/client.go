package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/yourorg/callcenter/internal/config"
)

func New(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("redis.ParseURL: %w", err)
	}
	return redis.NewClient(opts), nil
}

// Key helpers — centralising names prevents typos across packages.

func AgentKey(agentID string) string {
	return "agent:" + agentID
}

func AgentHeartbeatKey(agentID string) string {
	return "agent:" + agentID + ":heartbeat"
}

func AgentSetKey(status string) string {
	return "agents:" + status
}

func QueueWaitKey(queueID string) string {
	return "queue:" + queueID + ":waiting"
}

func CallKey(callID string) string {
	return "call:" + callID
}

func PubSubChannel(audience string) string {
	return "ws:" + audience
}

func AgentPubSubChannel(agentID string) string {
	return "ws:agent:" + agentID
}
