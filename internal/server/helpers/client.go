package helpers

import (
	"server/internal/core/shared/types"
	"server/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClientHelper struct {
	ctx      types.IContext
	services *services.Services
}

func (h *ClientHelper) GetIPAddress(context *gin.Context) (string, error) {
	clientIP := strings.TrimSpace(context.GetHeader("X-Forwarded-For"))
	if len(clientIP) == 0 {
		clientIP = context.Request.RemoteAddr
	}
	clientIPs := strings.Split(clientIP, ",")
	if len(clientIPs) > 0 {
		clientIP = clientIPs[0]
	}
	if strings.Contains(clientIP, ":") {
		split := strings.Split(clientIP, ":")
		if len(split) > 0 {
			value := split[0]
			value = strings.ReplaceAll(value, "[", "")
			value = strings.ReplaceAll(value, "]", "")
			clientIP = value
		}
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) == 0 {
		clientIP = context.ClientIP()
	}
	return clientIP, nil
}
