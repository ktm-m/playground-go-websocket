package outbound

import "github.com/labstack/echo/v4"

type EchoWebSocketHandlerPort interface {
	RegisterRoutes(e *echo.Echo)
	ServeHTML(c echo.Context) error
	EchoGorillaMuxWebSocket(c echo.Context) error
	EchoSocketIOWebSocket(c echo.Context) error
}
