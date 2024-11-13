package handler

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/handler/echo"
	"github.com/ktm-m/playground-go-websocket/handler/gin"
	"github.com/ktm-m/playground-go-websocket/internal/port/inbound"
	"github.com/ktm-m/playground-go-websocket/internal/port/outbound"
)

type Handler struct {
	EchoWebSocketHandler outbound.EchoWebSocketHandlerPort
	GinWebSocketHandler  outbound.GinWebSocketHandlerPort
}

func NewHandler(
	processMessageService inbound.ProcessMessagePort,
	upgrader *websocket.Upgrader,
	server *socketio.Server,
) *Handler {
	return &Handler{
		EchoWebSocketHandler: echo.NewEchoWebSocketHandler(processMessageService, upgrader, server),
		GinWebSocketHandler:  gin.NewGinWebSocketHandler(processMessageService, upgrader, server),
	}
}
