package echo

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/constant"
	"github.com/ktm-m/playground-go-websocket/helper"
	"github.com/ktm-m/playground-go-websocket/internal/port/inbound"
	"github.com/ktm-m/playground-go-websocket/internal/port/outbound"
	"github.com/labstack/echo/v4"
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
	socketIO              *socketio.Server
	muxWebSocketHelper    helper.MuxWebSocketHelper
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/echo")
	group.GET("/html", h.ServeHTML)
	group.GET("/gorilla-mux", h.GorillaMuxWebSocket)
	group.GET("/socket-io", h.SocketIOWebSocket)
}

func NewHandler(
	processMessageService inbound.ProcessMessagePort,
	upgrader *websocket.Upgrader,
	socketIO *socketio.Server,
	muxWebSocketHelper helper.MuxWebSocketHelper,
) outbound.EchoWebSocketHandlerPort {
	return &handler{
		processMessageService: processMessageService,
		upgrader:              upgrader,
		socketIO:              socketIO,
		muxWebSocketHelper:    muxWebSocketHelper,
	}
}

func (h *handler) ServeHTML(c echo.Context) error {
	return c.File("./html/chat.html")
}

func (h *handler) GorillaMuxWebSocket(c echo.Context) error {
	conn, err := h.muxWebSocketHelper.UpgradeConnection(h.upgrader, c)
	if err != nil {
		log.Println("[HANDLER] failed to upgrade connection:", err)
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			log.Println("[HANDLER] failed to close connection:", err)
		}
	}(conn)

	h.muxWebSocketHelper.AddClient(&mu, conn, clients)
	defer h.muxWebSocketHelper.RemoveClient(&mu, conn, clients)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[HANDLER] failed to read message:", err)
			return nil
		}

		resp, err := h.processMessageService.ProcessMessage(string(msg), constant.EchoServer)
		if err != nil {
			log.Println("[HANDLER] failed to process message:", err)
			return nil
		}

		h.muxWebSocketHelper.BroadCastMessage(&mu, []byte(resp), clients)
	}
}

func (h *handler) SocketIOWebSocket(c echo.Context) error {
	h.socketIO.OnEvent("/", "message", func(conn socketio.Conn, msg string) {
		resp, err := h.processMessageService.ProcessMessage(msg, constant.EchoServer)
		if err != nil {
			log.Println("[HANDLER] failed to process message:", err)
			conn.Emit("error", err)
			return
		}

		conn.Emit("message", resp)
	})

	h.socketIO.ServeHTTP(c.Response(), c.Request())
	return nil
}
