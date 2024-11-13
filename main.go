package main

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/constant"
	"github.com/ktm-m/playground-go-websocket/handler"
	"github.com/ktm-m/playground-go-websocket/infra"
	"github.com/ktm-m/playground-go-websocket/internal/service"
	"sync"
)

func main() {
	appConfig := infra.InitConfig()

	echoFactory := infra.NewFactory(&appConfig.App, constant.EchoServer)
	ginFactory := infra.NewFactory(&appConfig.App, constant.GinServer)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  appConfig.App.Upgrader.ReadBufferSize,
		WriteBufferSize: appConfig.App.Upgrader.WriteBufferSize,
	}
	socketIOServer := socketio.NewServer(nil)

	services := service.NewService()
	handlers := handler.NewHandler(services.ProcessMessageService, &upgrader, socketIOServer)

	echoServer := echoFactory.CreateServer()
	echoServer.AddHandler(handlers.EchoWebSocketHandler)

	ginServer := ginFactory.CreateServer()
	ginServer.AddHandler(handlers.GinWebSocketHandler)

	servers := []infra.HTTPServer{
		echoServer,
		ginServer,
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

//http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "chat.html")
//})
//
