package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/constant"
	"github.com/ktm-m/playground-go-websocket/helper"
	"github.com/ktm-m/playground-go-websocket/internal/port/inbound"
	"github.com/ktm-m/playground-go-websocket/internal/port/outbound"
	"log"
	"sync"
)

var (
	mu      sync.Mutex
	clients = make(map[*websocket.Conn]struct{})
)

type handler struct {
	processMessageService inbound.ProcessMessagePort
	upgrader              *websocket.Upgrader
	muxWebSocketHelper    helper.MuxWebSocketHelper
}

func (h *handler) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/gin")
	group.GET("/gorilla-mux", h.GorillaMuxWebSocket)
}

func NewHandler(
	processMessageService inbound.ProcessMessagePort,
	upgrader *websocket.Upgrader,
	muxWebSocketHelper helper.MuxWebSocketHelper,
) outbound.GinWebSocketHandlerPort {
	return &handler{
		processMessageService: processMessageService,
		upgrader:              upgrader,
		muxWebSocketHelper:    muxWebSocketHelper,
	}
}

func (h *handler) GorillaMuxWebSocket(c *gin.Context) {
	conn, err := h.muxWebSocketHelper.UpgradeConnection(h.upgrader, c)
	if err != nil {
		log.Println("[HANDLER] failed to upgrade connection:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			log.Println("[HANDLER] failed to close connection:", err)
			return
		}
	}(conn)

	h.muxWebSocketHelper.AddClient(&mu, conn, clients)
	defer h.muxWebSocketHelper.RemoveClient(&mu, conn, clients)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[HANDLER] failed to read message:", err)
			break
		}

		resp, err := h.processMessageService.ProcessMessage(string(msg), constant.GinServer)
		if err != nil {
			log.Println("[HANDLER] failed to process message:", err)
			break
		}

		h.muxWebSocketHelper.BroadCastMessage(&mu, []byte(resp), clients)
	}
}
