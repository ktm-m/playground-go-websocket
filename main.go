package main

import (
	"github.com/ktm-m/playground-go-websocket/constant"
	"github.com/ktm-m/playground-go-websocket/infra"
	"sync"
)

func main() {
	appConfig := infra.InitConfig()

	echoFactory := infra.NewFactory(&appConfig.App, constant.EchoServer)
	ginFactory := infra.NewFactory(&appConfig.App, constant.GinServer)

	servers := []infra.HTTPServer{
		echoFactory.CreateServer(),
		ginFactory.CreateServer(),
	}

	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		go func(s infra.HTTPServer) {
			defer wg.Done()
			s.Start()
		}(server)
	}

	go infra.ListenForShutdown(servers)
	wg.Wait()
}

//type ChatServer struct {
//	mu      sync.Mutex
//	clients map[*websocket.Conn]struct{}
//}
//
//func (cs *ChatServer) AddClient(conn *websocket.Conn) {
//	cs.mu.Lock() // Lock the mutex to prevent concurrent access to the map
//	defer cs.mu.Unlock()
//
//	cs.clients[conn] = struct{}{}
//}
//
//func (cs *ChatServer) RemoveClient(conn *websocket.Conn) {
//	cs.mu.Lock()
//	defer cs.mu.Unlock()
//
//	delete(cs.clients, conn)
//}
//
//func (cs *ChatServer) Broadcast(msg []byte) {
//	cs.mu.Lock()
//	defer cs.mu.Unlock()
//
//	for conn := range cs.clients {
//		err := conn.WriteMessage(websocket.TextMessage, msg)
//		if err != nil {
//			log.Println("[MAIN] failed to write message:", err)
//			continue
//		}
//	}
//}

//upgrader := websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//}
//
//cs := &ChatServer{
//	clients: make(map[*websocket.Conn]struct{}),
//}
//
//http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "chat.html")
//})
//
//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//	// Upgrade the HTTP connection to a WebSocket connection
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println("[MAIN] failed to upgrade connection:", err)
//		return
//	}
//	defer func(conn *websocket.Conn) {
//		err = conn.Close()
//		if err != nil {
//			log.Println("[MAIN] failed to close connection:", err)
//		}
//	}(conn)
//
//	cs.AddClient(conn)
//	defer cs.RemoveClient(conn)
//
//	for {
//		_, msg, err := conn.ReadMessage()
//		if err != nil {
//			log.Println("[MAIN] failed to read message:", err)
//			return
//		}
//
//		cs.Broadcast(msg)
//	}
//})
//err := http.ListenAndServe(":8080", nil)
//if err != nil {
//	log.Fatalf("[MAIN] failed to start server: %s", err.Error())
//}
