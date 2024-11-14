package echo

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/constant"
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

type echoWebSocketHandler struct {
	processMessageService inbound.ProcessMessagePort
	upgrader              *websocket.Upgrader
	socketIO              *socketio.Server
}

func (h *echoWebSocketHandler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/echo")
	group.GET("/html", h.ServeHTML)
	group.GET("/gorilla-mux", h.EchoGorillaMuxWebSocket)
	group.GET("/socket-io", h.EchoSocketIOWebSocket)
}

func NewEchoWebSocketHandler(processMessageService inbound.ProcessMessagePort, upgrader *websocket.Upgrader, socketIO *socketio.Server) outbound.EchoWebSocketHandlerPort {
	return &echoWebSocketHandler{
		processMessageService: processMessageService,
		upgrader:              upgrader,
		socketIO:              socketIO,
	}
}

func (h *echoWebSocketHandler) ServeHTML(c echo.Context) error {
	return c.File("./html/chat.html")
}

func (h *echoWebSocketHandler) EchoGorillaMuxWebSocket(c echo.Context) error {
	conn, err := h.upgradeConnection(h.upgrader, c)
	if err != nil {
		log.Println("[HANDLER] failed to upgrade connection:", err)
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			log.Println("[HANDLER] failed to close connection:", err)
		}
	}(conn)

	h.addClient(conn)
	defer h.removeClient(conn)

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

		h.broadcastMessage([]byte(resp), clients)
	}
}

func (h *echoWebSocketHandler) EchoSocketIOWebSocket(c echo.Context) error {
	h.socketIO.OnEvent("/", "message", func(conn socketio.Conn, msg string) {
		resp, err := h.processMessageService.ProcessMessage(msg, constant.EchoServer)
		if err != nil {
			log.Println("[HANDLER] failed to process message:", err)
			conn.Emit("error", err)
			return
		}

		conn.Emit("response", resp)
	})

	h.socketIO.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *echoWebSocketHandler) upgradeConnection(upgrader *websocket.Upgrader, c echo.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("[HANDLER] failed to upgrade connection:", err)
		return nil, err
	}

	return conn, nil
}

func (h *echoWebSocketHandler) addClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	clients[conn] = struct{}{}
}

func (h *echoWebSocketHandler) removeClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, conn)
}

func (h *echoWebSocketHandler) broadcastMessage(msg []byte, clients map[*websocket.Conn]struct{}) {
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
