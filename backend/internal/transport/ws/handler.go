package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yourorg/callcenter/internal/auth"
	"github.com/yourorg/callcenter/internal/config"
	"github.com/yourorg/callcenter/internal/hub"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // tighten in production with allowlist
	},
}

func AgentHandler(h *hub.Hub, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := verifyWSToken(cfg, r)
		if err != nil || claims.AgentID == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warn().Err(err).Msg("ws upgrade failed")
			return
		}

		log.Info().Str("agent_id", claims.AgentID.String()).Msg("agent ws connected")
		h.RegisterAgent(conn, claims.AgentID.String())
	}
}

func SupervisorHandler(h *hub.Hub, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := verifyWSToken(cfg, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if claims.Role != "supervisor" && claims.Role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warn().Err(err).Msg("ws upgrade failed")
			return
		}

		log.Info().Str("user_id", claims.UserID.String()).Msg("supervisor ws connected")
		h.RegisterSupervisor(conn, claims.UserID.String())
	}
}

// verifyWSToken extracts ?token= from the query string and verifies it.
// WebSocket clients cannot send Authorization headers during the upgrade.
func verifyWSToken(cfg *config.Config, r *http.Request) (*auth.Claims, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		// Also accept Bearer header (for non-browser clients)
		token = r.Header.Get("Authorization")
		if len(token) > 7 {
			token = token[7:]
		}
	}
	return auth.VerifyAccessToken(cfg, token)
}
