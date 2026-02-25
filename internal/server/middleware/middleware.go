package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/stringutil"
	"server/internal/server/helpers"
	"server/internal/services"
	"server/internal/services/core/setup"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	ctx        types.IContext
	helpers    *helpers.Helper
	services   *services.Services
	loggerChan chan *gin.Context
}

func New(ctx types.IContext, services *services.Services, helpers *helpers.Helper) *Middleware {
	mw := &Middleware{
		ctx:        ctx,
		helpers:    helpers,
		services:   services,
		loggerChan: make(chan *gin.Context, 1000000),
	}
	go func(ctx types.IContext, middleware *Middleware) {
		for {
			ginCtx := <-middleware.loggerChan
			log := []string{ginCtx.Request.Method}
			pth := ginCtx.Request.URL.Path
			sze := ginCtx.Writer.Size()
			sts := ginCtx.Writer.Status()
			if sze < 0 {
				sze = 0
			}
			msg := strings.TrimSpace(ginCtx.Errors.ByType(gin.ErrorTypePrivate).String())
			if len(ginCtx.Request.URL.RawQuery) > 0 {
				pth = fmt.Sprintf("%s?%s", pth, ginCtx.Request.URL.RawQuery)
			}
			log = append(log, pth)
			log = append(log, fmt.Sprintf("%d", sts))
			log = append(log, fmt.Sprintf("%dB", sze))
			log = append(log, ginCtx.ClientIP())
			if len(msg) > 0 {
				log = append(log, msg)
			}
			ctx.Logger().HTTP(strings.Join(log, " | "))
		}
	}(ctx, mw)
	return mw
}

func (m *Middleware) Cors() gin.HandlerFunc {
	crs := make(map[string]string)
	crs["Access-Control-Allow-Origin"] = "*"
	crs["Access-Control-Allow-Headers"] = "*"
	crs["Access-Control-Expose-Headers"] = "Authorization"
	crs["Access-Control-Allow-Credentials"] = "true"
	crs["Access-Control-Allow-Methods"] = "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH"
	return func(context *gin.Context) {
		for key, val := range crs {
			context.Writer.Header().Set(key, val)
		}
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(204)
			return
		}
		context.Next()
	}
}

func (m *Middleware) Cache() gin.HandlerFunc {
	return func(context *gin.Context) {
		cacheEnabled := false
		if strings.HasPrefix(context.Request.URL.Path, "/files/") {
			cacheEnabled = true
		}
		if strings.HasPrefix(context.Request.URL.Path, "/assets/") {
			cacheEnabled = true
		}
		if cacheEnabled {
			context.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		context.Next()
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := types.ApiResponseBody{}
		success, err := m.helpers.Auth().Login(context, helpers.AuthTypePrivateApi)
		if err != nil {
			response.Code = types.ServiceResultCodeInternal
			context.IndentedJSON(http.StatusInternalServerError, response)
			context.Abort()
			m.ctx.Logger().Error(err)
			return
		}
		if !success {
			response.Code = types.ServiceResultCodeUnauthorized
			context.IndentedJSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}
		context.Next()
	}
}

func (m *Middleware) Setup() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := types.ApiResponseBody{}
		checkSetup, err := m.services.Setup().Get(setup.GetData{})
		if err != nil {
			response.Code = types.ServiceResultCodeInternal
			context.IndentedJSON(http.StatusInternalServerError, response)
			context.Abort()
			m.ctx.Logger().Error(err)
			return
		}
		if checkSetup.Code != types.ServiceResultCodeSuccess {
			response.Code = checkSetup.Code
			response.Data = checkSetup.Payload
			context.IndentedJSON(http.StatusForbidden, response)
			context.Abort()
			return
		}
		if !checkSetup.Payload.Setup.Database {
			response.Code = types.ServiceResultCodeInvalidSetup
			response.Data = checkSetup.Payload
			context.IndentedJSON(http.StatusForbidden, response)
			context.Abort()
			return
		}
		if !checkSetup.Payload.Setup.Admin {
			response.Code = types.ServiceResultCodeInvalidSetup
			response.Data = checkSetup.Payload
			context.IndentedJSON(http.StatusForbidden, response)
			context.Abort()
			return
		}
		context.Next()
	}
}

func (m *Middleware) OnlyOwner() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := types.ApiResponseBody{}
		account, authenticated := m.helpers.Auth().GetUser(context)
		if !authenticated || len(account.ID) == 0 {
			response.Code = types.ServiceResultCodeUnauthorized
			context.IndentedJSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}
		path := context.Request.URL.Path
		if regexp.MustCompile(`/coaches/(\d+)`).MatchString(path) {

		}

		if regexp.MustCompile(`/athletes/(\d+)`).MatchString(path) {

		}

		if regexp.MustCompile(`/packages/(\d+)`).MatchString(path) {

		}

		if regexp.MustCompile(`/locations/(\d+)`).MatchString(path) {

		}
		context.Next()
	}
}

func (m *Middleware) AllowedRoles(roles []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		response := types.ApiResponseBody{}
		account, authenticated := m.helpers.Auth().GetUser(context)
		if !authenticated || len(account.ID) == 0 {
			response.Code = types.ServiceResultCodeUnauthorized
			context.IndentedJSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}
		if !stringutil.InSlice(string(account.Role), roles) {
			response.Code = types.ServiceResultCodeForbidden
			context.IndentedJSON(http.StatusForbidden, response)
			context.Abort()
			return
		}
		context.Next()
	}
}

func (m *Middleware) Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		m.loggerChan <- context
	}
}

func (m *Middleware) Swagger() gin.HandlerFunc {
	return func(context *gin.Context) {
		//host := context.ts.Request.Host
		//swagger.AppInfo.Host = host
	}
}

func (m *Middleware) Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
