package infra

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ktm-m/playground-go-websocket/constant"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type HTTPServer interface {
	Start()
	Shutdown(ctx context.Context)
	Info() string
	AddHandler(handler ...interface{})
	GetInstance() interface{}
}

type HTTPServerFactory interface {
	CreateServer() HTTPServer
}

type EchoServer struct {
	echo          *echo.Echo
	config        *App
	routeHandlers []EchoRouteHandler
	once          sync.Once
}

type GinServer struct {
	gin           *gin.Engine
	config        *App
	server        *http.Server
	routeHandlers []GinRouteHandler
	once          sync.Once
}

type EchoRouteHandler interface {
	RegisterRoutes(e *echo.Echo)
}

type GinRouteHandler interface {
	RegisterRoutes(e *gin.Engine)
}

type EchoServerFactory struct {
	config *App
}

type GinServerFactory struct {
	config *App
}

func (esf *EchoServerFactory) CreateServer() HTTPServer {
	echoInstance := echo.New()
	echoServer := &EchoServer{
		echo:   echoInstance,
		config: esf.config,
	}
	echoServer.setupMiddleware()
	echoServer.setupHealthEndpoint()

	for _, handler := range echoServer.routeHandlers {
		handler.RegisterRoutes(echoInstance)
	}

	return echoServer
}

func (gsf *GinServerFactory) CreateServer() HTTPServer {
	gin.SetMode(gin.ReleaseMode) // Disable debug mode
	ginInstance := gin.Default()

	err := ginInstance.SetTrustedProxies(gsf.config.TrustProxies)
	if err != nil {
		panic("[INFRA] failed to set trusted proxies")
	}

	ginServer := &GinServer{
		gin:    ginInstance,
		config: gsf.config,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", gsf.config.GinPort),
			Handler: ginInstance,
		},
	}
	ginServer.setupMiddleware()
	ginServer.setupHealthEndpoint()

	for _, handler := range ginServer.routeHandlers {
		handler.RegisterRoutes(ginInstance)
	}

	return ginServer
}

func NewFactory(config *App, serverType string) HTTPServerFactory {
	switch serverType {
	case constant.EchoServer:
		return &EchoServerFactory{config: config}
	case constant.GinServer:
		return &GinServerFactory{config: config}
	default:
		panic("[INFRA] invalid server type")
		return nil
	}
}

func (es *EchoServer) Start() {
	es.once.Do(func() {
		addr := fmt.Sprintf(":%s", es.config.EchoPort)
		es.echo.HideBanner = true
		es.echo.HidePort = true

		for _, route := range es.echo.Routes() {
			log.Printf("[ECHO] %s:%s", route.Method, route.Path)
		}

		if err := es.echo.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("[INFRA] failed to start echo server")
		}
	})
}

func (gs *GinServer) Start() {
	gs.once.Do(func() {
		for _, route := range gs.gin.Routes() {
			log.Printf("[GIN] %s:%s", route.Method, route.Path)
		}

		if err := gs.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("[INFRA] failed to start gin server")
		}
	})
}

func (es *EchoServer) Shutdown(ctx context.Context) {
	if err := es.echo.Shutdown(ctx); err != nil {
		panic("[INFRA] failed to shutdown echo server")
	}
}

func (gs *GinServer) Shutdown(ctx context.Context) {
	if err := gs.server.Shutdown(ctx); err != nil {
		panic("[INFRA] failed to shutdown gin server")
	}
}

func (es *EchoServer) Info() string {
	return fmt.Sprintf("[INFRA] echo server is running on port %s", es.config.EchoPort)
}

func (gs *GinServer) Info() string {
	return fmt.Sprintf("[INFRA] gin server is running on port %s", gs.config.GinPort)
}

func (es *EchoServer) AddHandler(handler ...interface{}) {
	for _, h := range handler {
		if routeHandler, ok := h.(EchoRouteHandler); ok {
			es.routeHandlers = append(es.routeHandlers, routeHandler)
		} else {
			panic("[INFRA] invalid handler type for echo server")
		}
	}
}

func (gs *GinServer) AddHandler(handler ...interface{}) {
	for _, h := range handler {
		if routeHandler, ok := h.(GinRouteHandler); ok {
			gs.routeHandlers = append(gs.routeHandlers, routeHandler)
		} else {
			panic("[INFRA] invalid handler type for gin server")
		}
	}
}

func (es *EchoServer) setupMiddleware() {
	es.echo.Use(
		middleware.Recover(),
		middleware.Logger(),
	)
}

func (gs *GinServer) setupMiddleware() {
	gs.gin.Use(
		gin.Recovery(),
		gin.Logger(),
	)
}

func (es *EchoServer) setupHealthEndpoint() {
	es.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"name":    es.config.Name,
			"port":    es.config.EchoPort,
			"version": es.config.Version,
		})
	})
}

func (gs *GinServer) setupHealthEndpoint() {
	gs.gin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    gs.config.Name,
			"port":    gs.config.GinPort,
			"version": gs.config.Version,
		})
	})
}

func (es *EchoServer) GetInstance() interface{} {
	return es.echo
}

func (gs *GinServer) GetInstance() interface{} {
	return gs.gin
}

func ListenForShutdown(servers []HTTPServer) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	shutdownServer(servers, signalChan)
}

func shutdownServer(servers []HTTPServer, signalChan <-chan os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	<-signalChan
	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		go func(s HTTPServer) {
			defer wg.Done()
			s.Shutdown(ctx)
		}(server)
	}
	wg.Wait()
}
