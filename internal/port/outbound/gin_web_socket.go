package outbound

import (
	"github.com/gin-gonic/gin"
)

type GinWebSocketHandlerPort interface {
	RegisterRoutes(e *gin.Engine)
	GinGorillaMuxWebSocket(c *gin.Context)
	GinSocketIOWebSocket(c *gin.Context)
}
