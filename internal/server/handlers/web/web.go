package web

import (
	"io"
	"io/fs"
	"net/http"
	"server/frontend"
	"server/internal/core/shared/types"
	"server/internal/server/helpers"
	"server/internal/server/middleware"
	"server/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type handler struct {
	ctx      types.IContext
	helpers  *helpers.Helper
	services *services.Services
}

func Setup(ctx types.IContext, engine *gin.Engine, services *services.Services, helpers *helpers.Helper, _ *middleware.Middleware) error {
	h := &handler{
		ctx:      ctx,
		helpers:  helpers,
		services: services,
	}
	spa, err := fs.Sub(frontend.Web, "web/dist")
	if err != nil {
		return err
	}
	spaStaticFs := http.FS(spa)
	spaFileServer := http.FileServer(spaStaticFs)

	engine.GET("/files/images/*any", h.serveStaticImage)
	engine.HEAD("/files/images/*any", h.serveStaticImage)
	engine.NoRoute(func(context *gin.Context) {
		url := context.Request.URL
		path := url.Path
		skip := []string{
			"/api/v1", "/app/v1", "/.well-known",
		}
		for _, val := range skip {
			if strings.HasPrefix(path, val) {
				context.Next()
				return
			}
		}
		var asset bool
		if strings.HasPrefix(path, "/assets") {
			asset = true
		}
		if strings.HasPrefix(path, "/favicon") {
			asset = true
		}
		if asset {
			spaFileServer.ServeHTTP(context.Writer, context.Request)
		} else {
			f, err := spa.Open("index.html")
			if err != nil {
				ctx.Logger().Error(err)
				return
			}
			defer func(f fs.File) {
				_ = f.Close()
			}(f)
			content, err := io.ReadAll(f)
			if err != nil {
				ctx.Logger().Error(err)
				return
			}
			context.Data(http.StatusOK, "text/html; charset=utf-8", content)
		}
	})

	return nil
}
