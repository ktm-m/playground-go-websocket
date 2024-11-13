package gin

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/internal/port/inbound"
	"github.com/ktm-m/playground-go-websocket/internal/port/outbound"
)

type ginWebSocketHandler struct {
	processMessageService inbound.ProcessMessagePort
	upgrader              *websocket.Upgrader
	socketIO              *socketio.Server
}

func (h *ginWebSocketHandler) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/gin")
	group.GET("/gorilla-mux", h.GinGorillaMuxWebSocket)
	group.GET("/socket-io", h.GinSocketIOWebSocket)
}

func NewGinWebSocketHandler(processMessageService inbound.ProcessMessagePort, upgrader *websocket.Upgrader, socketIO *socketio.Server) outbound.GinWebSocketHandlerPort {
	return &ginWebSocketHandler{
		processMessageService: processMessageService,
		upgrader:              upgrader,
		socketIO:              socketIO,
	}
}

func (h *ginWebSocketHandler) GinGorillaMuxWebSocket(c *gin.Context) {

}

func (h *ginWebSocketHandler) GinSocketIOWebSocket(c *gin.Context) {}
