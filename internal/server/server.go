package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/server/handlers"
	"server/internal/server/helpers"
	"server/internal/server/middleware"
	"server/internal/services"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	conf "github.com/tavsec/gin-healthcheck/config"
)

type server struct {
	ctx        types.IContext
	engine     *gin.Engine
	helper     *helpers.Helper
	server     *http.Server
	address    string
	services   *services.Services
	connection net.PacketConn
	middleware *middleware.Middleware
	eventsChan chan *types.ServerEvent
}

func New(ctx types.IContext, services *services.Services) types.IServer {
	s := &server{
		ctx:        ctx,
		services:   services,
		eventsChan: make(chan *types.ServerEvent),
	}
	ctx.Hooks().OnStop(s.stop)
	ctx.Hooks().OnStart(s.start)
	return s
}

func (s *server) start() error {
	s.ctx.Logger().Debug("starting http server...")
	cnf, err := s.ctx.Config().Get()
	if err != nil {
		return err
	}
	if cnf.Server.HTTP.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine = gin.New()
	s.engine.MaxMultipartMemory = int64(cnf.Server.HTTP.MaxMultipartMemory)
	s.engine.ForwardedByClientIP = cnf.Server.HTTP.ForwardedByClientIP
	s.helper = helpers.New(s.ctx, s.services)
	s.address = cnf.Server.HTTP.Address

	s.middleware = middleware.New(s.ctx, s.services, s.helper)
	s.engine.Use(
		s.middleware.Cors(), s.middleware.Cache(),
		s.middleware.Swagger(), s.middleware.Logger(), s.middleware.Recovery(),
	)

	if err = handlers.Setup(s.ctx, s.engine, s.services, s.helper, s.middleware, s.eventsChan); err != nil {
		return err
	}

	if err = healthcheck.New(s.engine, conf.DefaultConfig(), []checks.Check{}); err != nil {
		s.ctx.Logger().Error(err)
	}

	s.server = &http.Server{
		Addr:    s.address,
		Handler: s.engine,
	}
	go func(s *server) {
		time.Sleep(time.Second * 1)
		if err = s.server.ListenAndServe(); err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "closed") {
				s.ctx.Logger().Fatal(err)
			}
		}
	}(s)
	return s.waitStart()
}

func (s *server) stop() error {
	if s.server != nil {
		s.ctx.Logger().Debug("stopping http server...")
		if err := s.server.Shutdown(context.Background()); err != nil {
			s.ctx.Logger().Warn(err)
			return nil
		}
		s.ctx.Logger().Debug("http server stopped")
	}
	return nil
}

func (s *server) waitStart() error {
	ctx, ctxCancelFn := context.WithDeadline(
		context.Background(),
		time.Now().Add(15*time.Second),
	)
	defer ctxCancelFn()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout while starting http server on %s", s.address)
		default:
			conn, err := net.DialTimeout("tcp", s.address, time.Second*2)
			if err == nil && conn != nil {
				if err := conn.Close(); err != nil {
					s.ctx.Logger().Warn(err)
				}
				s.ctx.Logger().Debugf("http server started on %s", s.address)
				return nil
			}
			time.Sleep(time.Second)
		}
	}
}

func (s *server) SendEvent(event *types.ServerEvent) {
	s.eventsChan <- event
}
