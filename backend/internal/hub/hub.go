package hub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	goredis "github.com/redis/go-redis/v9"
	rkeys "github.com/yourorg/callcenter/internal/redis"
	"github.com/rs/zerolog/log"
)

type Hub struct {
	agents      map[string]*client // agentID → client
	supervisors map[string]*client // userID  → client

	register   chan *client
	unregister chan *client
	send       chan envelope

	rdb *goredis.Client
}

type envelope struct {
	target string // agentID, userID, or "" for supervisors broadcast
	role   clientRole
	event  Event
}

func New(rdb *goredis.Client) *Hub {
	return &Hub{
		agents:      make(map[string]*client),
		supervisors: make(map[string]*client),
		register:    make(chan *client, 64),
		unregister:  make(chan *client, 64),
		send:        make(chan envelope, 512),
		rdb:         rdb,
	}
}

func (h *Hub) Run(ctx context.Context) {
	go h.subscribeRedis(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-h.register:
			if c.role == roleAgent {
				h.agents[c.id] = c
				log.Debug().Str("agent_id", c.id).Msg("agent connected")
			} else {
				h.supervisors[c.id] = c
				log.Debug().Str("user_id", c.id).Msg("supervisor connected")
			}
		case c := <-h.unregister:
			if c.role == roleAgent {
				if existing, ok := h.agents[c.id]; ok && existing == c {
					delete(h.agents, c.id)
					close(c.send)
				}
			} else {
				if existing, ok := h.supervisors[c.id]; ok && existing == c {
					delete(h.supervisors, c.id)
					close(c.send)
				}
			}
		case env := <-h.send:
			msg, err := marshalEvent(env.event)
			if err != nil {
				continue
			}
			switch {
			case env.role == roleAgent && env.target != "":
				if c, ok := h.agents[env.target]; ok {
					select {
					case c.send <- msg:
					default:
						close(c.send)
						delete(h.agents, env.target)
					}
				}
			case env.role == roleSupervisor:
				for _, c := range h.supervisors {
					select {
					case c.send <- msg:
					default:
						close(c.send)
						delete(h.supervisors, c.id)
					}
				}
			}
		}
	}
}

func (h *Hub) RegisterAgent(conn *websocket.Conn, agentID string) {
	c := &client{hub: h, conn: conn, send: make(chan []byte, 256), id: agentID, role: roleAgent}
	h.register <- c
	go c.writePump()
	go c.readPump()
}

func (h *Hub) RegisterSupervisor(conn *websocket.Conn, userID string) {
	c := &client{hub: h, conn: conn, send: make(chan []byte, 256), id: userID, role: roleSupervisor}
	h.register <- c
	go c.writePump()
	go c.readPump()
}

func (h *Hub) SendToAgent(agentID string, event Event) {
	h.send <- envelope{target: agentID, role: roleAgent, event: event}
}

func (h *Hub) BroadcastSupervisors(event Event) {
	h.send <- envelope{role: roleSupervisor, event: event}
}

// PublishToAgent publishes via Redis so all instances fan out the event.
func (h *Hub) PublishToAgent(ctx context.Context, agentID string, event Event) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return h.rdb.Publish(ctx, rkeys.AgentPubSubChannel(agentID), msg).Err()
}

// PublishSupervisors publishes to the supervisors channel via Redis.
func (h *Hub) PublishSupervisors(ctx context.Context, event Event) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return h.rdb.Publish(ctx, rkeys.PubSubChannel("supervisors"), msg).Err()
}

func (h *Hub) subscribeRedis(ctx context.Context) {
	sub := h.rdb.PSubscribe(ctx, "ws:*")
	defer sub.Close()
	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Warn().Err(err).Msg("invalid redis pubsub event")
				continue
			}
			// Determine routing from channel name
			channel := msg.Channel
			switch {
			case len(channel) > len("ws:agent:") && channel[:9] == "ws:agent:":
				agentID := channel[9:]
				h.send <- envelope{target: agentID, role: roleAgent, event: event}
			case channel == "ws:supervisors":
				h.send <- envelope{role: roleSupervisor, event: event}
			}
		}
	}
}
