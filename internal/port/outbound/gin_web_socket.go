package outbound

import (
	"github.com/gin-gonic/gin"
)

type GinWebSocketHandlerPort interface {
	RegisterRoutes(e *gin.Engine)
	GorillaMuxWebSocket(c *gin.Context)
	SocketIOWebSocket(c *gin.Context)
}
