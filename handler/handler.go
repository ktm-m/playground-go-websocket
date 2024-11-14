package handler

import (
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/handler/echo"
	"github.com/ktm-m/playground-go-websocket/handler/gin"
	"github.com/ktm-m/playground-go-websocket/helper"
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
	muxWebSocketHelper helper.MuxWebSocketHelper,
) *Handler {
	return &Handler{
		EchoWebSocketHandler: echo.NewHandler(processMessageService, upgrader, muxWebSocketHelper),
		GinWebSocketHandler:  gin.NewHandler(processMessageService, upgrader, muxWebSocketHelper),
	}
}
