package handlers

import (
	"server/internal/core/shared/types"
	"server/internal/server/handlers/app"
	"server/internal/server/handlers/web"
	"server/internal/server/helpers"
	"server/internal/server/middleware"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

func Setup(ctx types.IContext, engine *gin.Engine, services *services.Services, helpers *helpers.Helper, middleware *middleware.Middleware, eventsChan chan *types.ServerEvent) error {
	if err := app.Setup(ctx, engine, services, helpers, middleware, eventsChan); err != nil {
		return err
	}
	if err := web.Setup(ctx, engine, services, helpers, middleware); err != nil {
		return err
	}
	return nil
}
