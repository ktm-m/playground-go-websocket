package gin

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/constant"
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
	socketIO              *socketio.Server
}

func (h *handler) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/gin")
	group.GET("/gorilla-mux", h.GorillaMuxWebSocket)
	group.GET("/socket-io", h.SocketIOWebSocket)
}

func NewHandler(processMessageService inbound.ProcessMessagePort, upgrader *websocket.Upgrader, socketIO *socketio.Server) outbound.GinWebSocketHandlerPort {
	return &handler{
		processMessageService: processMessageService,
		upgrader:              upgrader,
		socketIO:              socketIO,
	}
}

func (h *handler) GorillaMuxWebSocket(c *gin.Context) {
	conn, err := h.upgradeConnection(h.upgrader, c)
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

	h.addClient(conn)
	defer h.removeClient(conn)

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

		h.broadcastMessage([]byte(resp), clients)
	}
}

func (h *handler) SocketIOWebSocket(c *gin.Context) {
	h.socketIO.OnEvent("/", "message", func(conn socketio.Conn, msg string) {
		resp, err := h.processMessageService.ProcessMessage(msg, constant.GinServer)
		if err != nil {
			log.Println("[HANDLER] failed to process message:", err)
			conn.Emit("error", err)
			return
		}

		conn.Emit("message", resp)
	})

	h.socketIO.ServeHTTP(c.Writer, c.Request)
}

func (h *handler) upgradeConnection(upgrader *websocket.Upgrader, c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[HANDLER] failed to upgrade connection:", err)
		return nil, err
	}

	return conn, nil
}

func (h *handler) addClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	clients[conn] = struct{}{}
}

func (h *handler) removeClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, conn)
}

func (h *handler) broadcastMessage(msg []byte, clients map[*websocket.Conn]struct{}) {
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("[HANDLER] failed to write message:", err)
			continue
		}
	}
}
