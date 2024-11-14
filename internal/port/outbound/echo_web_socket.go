package outbound

import "github.com/labstack/echo/v4"

type EchoWebSocketHandlerPort interface {
	RegisterRoutes(e *echo.Echo)
	ServeHTML(c echo.Context) error
	GorillaMuxWebSocket(c echo.Context) error
	SocketIOWebSocket(c echo.Context) error
}
