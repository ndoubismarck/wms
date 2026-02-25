package helpers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"server/internal/core/shared/types"
	"server/internal/services"
)

type SwaggerHelper struct {
	ctx      types.IContext
	services *services.Services
}

func (h *SwaggerHelper) Handler(name string) gin.HandlerFunc {
	//ToDo: Disable in production
	return ginSwagger.WrapHandler(swaggerFiles.NewHandler(),
		ginSwagger.InstanceName(name),
	)
}

func (h *SwaggerHelper) abortHandler(status int, err error) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.AbortWithStatus(status)
		if err != nil {
			h.ctx.Logger().Error(err)
		}
	}
}
