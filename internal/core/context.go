package core

import (
	"os"
	"server/internal/core/cache"
	"server/internal/core/config"
	"server/internal/core/database"
	"server/internal/core/dotenv"
	"server/internal/core/events"
	"server/internal/core/hooks"
	"server/internal/core/logger"
	"server/internal/core/shared/types"
	"server/internal/server"
	"server/internal/services"
	"server/internal/services/providers"
	"syscall"
)

type context struct {
	app      types.App
	cache    types.ICache
	hooks    types.IHooks
	server   types.IServer
	events   types.IEvents
	config   types.IConfig
	logger   types.ILogger
	dotenv   types.IDotEnv
	database types.IDatabase
}

func NewContext(app types.App, hooksChan chan types.Hook) types.IContext {
	c := &context{
		app:    app,
		hooks:  hooks.New(hooksChan),
		events: events.New(),
	}
	c.server = server.New(c, services.New(c, providers.New(c)))
	return c
}

func (c *context) App() types.App {
	return c.app
}

func (c *context) Cache() types.ICache {
	if c.cache == nil {
		c.cache = cache.New(c)
	}
	return c.cache
}

func (c *context) Hooks() types.IHooks {
	return c.hooks
}

func (c *context) Events() types.IEvents {
	return c.events
}

func (c *context) Config() types.IConfig {
	if c.config == nil {
		c.config = config.New(c, types.ConfigCallbacks{
			OnChangeFn: func(config types.Config) {
				process, err := os.FindProcess(syscall.Getpid())
				if err != nil {
					c.Logger().Fatal(err)
				}
				if err = process.Signal(syscall.SIGHUP); err != nil {
					c.Logger().Fatal(err)
				}
			},
		})
	}
	return c.config
}

func (c *context) Server() types.IServer {
	return c.server
}

func (c *context) Logger() types.ILogger {
	if c.logger == nil {
		c.logger = logger.New(c)
	}
	return c.logger
}

func (c *context) DotEnv() types.IDotEnv {
	if c.dotenv == nil {
		c.dotenv = dotenv.New(c)
	}
	return c.dotenv
}

func (c *context) Database() types.IDatabase {
	if c.database == nil {
		c.database = database.New(c)
	}
	return c.database
}
