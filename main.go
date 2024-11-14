package main

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/ktm-m/playground-go-websocket/constant"
	"github.com/ktm-m/playground-go-websocket/handler"
	"github.com/ktm-m/playground-go-websocket/helper"
	"github.com/ktm-m/playground-go-websocket/infra"
	"github.com/ktm-m/playground-go-websocket/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

func main() {
	appConfig := infra.InitConfig()

	echoFactory := infra.NewFactory(&appConfig.App, constant.EchoServer)
	ginFactory := infra.NewFactory(&appConfig.App, constant.GinServer)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  appConfig.App.Upgrader.ReadBufferSize,
		WriteBufferSize: appConfig.App.Upgrader.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	socketIOServer := socketio.NewServer(nil)

	muxWebSocketHelper := helper.NewMuxWebSocketHelper()

	services := service.NewService()
	handlers := handler.NewHandler(services.ProcessMessageService, &upgrader, socketIOServer, muxWebSocketHelper)

	echoServer := echoFactory.CreateServer()
	registerEchoHandlers(echoServer, handlers)

	ginServer := ginFactory.CreateServer()
	registerGinHandlers(ginServer, handlers)

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

func registerEchoHandlers(server infra.HTTPServer, handlers *handler.Handler) {
	echoWebSocketHandler := handlers.EchoWebSocketHandler
	echoWebSocketHandler.RegisterRoutes(server.GetInstance().(*echo.Echo))
	server.AddHandler(echoWebSocketHandler)
}

func registerGinHandlers(server infra.HTTPServer, handlers *handler.Handler) {
	ginWebSocketHandler := handlers.GinWebSocketHandler
	ginWebSocketHandler.RegisterRoutes(server.GetInstance().(*gin.Engine))
	server.AddHandler(ginWebSocketHandler)
}
