package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"sync"
)

type MuxWebSocketHelper interface {
	UpgradeConnection(upgrader *websocket.Upgrader, c interface{}) (*websocket.Conn, error)
	AddClient(mu *sync.Mutex, conn *websocket.Conn, clients map[*websocket.Conn]struct{})
	RemoveClient(mu *sync.Mutex, conn *websocket.Conn, clients map[*websocket.Conn]struct{})
	BroadCastMessage(mu *sync.Mutex, msg []byte, clients map[*websocket.Conn]struct{})
}

type muxWebSocketHelper struct{}

func NewMuxWebSocketHelper() MuxWebSocketHelper {
	return &muxWebSocketHelper{}
}

func (h *muxWebSocketHelper) UpgradeConnection(upgrader *websocket.Upgrader, c interface{}) (*websocket.Conn, error) {
	switch c.(type) {
	case *gin.Context:
		conn, err := upgrader.Upgrade(c.(*gin.Context).Writer, c.(*gin.Context).Request, nil)
		if err != nil {
			return nil, err
		}

		return conn, nil
	case echo.Context:
		conn, err := upgrader.Upgrade(c.(echo.Context).Response(), c.(echo.Context).Request(), nil)
		if err != nil {
			return nil, err
		}

		return conn, nil
	default:
		return nil, errors.New("invalid context")
	}
}

func (h *muxWebSocketHelper) AddClient(mu *sync.Mutex, conn *websocket.Conn, clients map[*websocket.Conn]struct{}) {
	mu.Lock()
	defer mu.Unlock()

	clients[conn] = struct{}{}
}

func (h *muxWebSocketHelper) RemoveClient(mu *sync.Mutex, conn *websocket.Conn, clients map[*websocket.Conn]struct{}) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, conn)
}

func (h *muxWebSocketHelper) BroadCastMessage(mu *sync.Mutex, msg []byte, clients map[*websocket.Conn]struct{}) {
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("[HELPER] failed to write message:", err)
			continue
		}
	}
}
