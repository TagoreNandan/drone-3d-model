package mockdata

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Development only. Replace with proper origin validation before production.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsMessage struct {
	Orientation
	Ts string `json:"ts"`
}

// HandleWS upgrades the connection and streams orientation telemetry frames.
func HandleWS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Run a separate reader loop to detect client disconnects immediately,
	// rather than waiting for the next WriteJSON failure.
	go func() {
		defer cancel()
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				break
			}
		}
	}()

	generator := NewGenerator(DefaultConfig())
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	lastTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			delta := t.Sub(lastTime).Seconds()
			lastTime = t

			state := generator.Next(delta)

			msg := wsMessage{
				Orientation: state,
				Ts:          t.UTC().Format(time.RFC3339Nano),
			}

			if err := ws.WriteJSON(msg); err != nil {
				return nil
			}
		}
	}
}
