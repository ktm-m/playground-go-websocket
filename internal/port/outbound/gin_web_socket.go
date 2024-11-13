package outbound

import "github.com/gin-gonic/gin"

type GinWebSocketHandlerPort interface {
	GinGorillaMuxWebSocket(c *gin.Context)
	GinSocketIOWebSocket(c *gin.Context)
}
