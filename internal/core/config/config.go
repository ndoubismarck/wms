package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/stringutil"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

type config struct {
	ctx      types.IContext
	init     atomic.Bool
	data     *types.Config
	watcher  *fsnotify.Watcher
	filePath string
	modeTime time.Time

	callback  atomic.Bool
	callbacks types.ConfigCallbacks
}

func New(ctx types.IContext, callbacks types.ConfigCallbacks) types.IConfig {
	c := &config{
		ctx:       ctx,
		filePath:  "./config/app.json",
		callbacks: callbacks,
	}
	c.init.Store(true)
	return c
}

func (c *config) Get() (*types.Config, error) {
	if err := c.load(); err != nil {
		return nil, err
	}
	if err := c.watch(); err != nil {
		return nil, err
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.data, nil
}

func (c *config) Set(config *types.Config) error {
	c.set(*config)
	return c.write(true)
}

func (c *config) load() error {
	dir := filepath.Dir(c.filePath)
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	write := false
	info, err := os.Stat(c.filePath)
	if err != nil {
		var data types.Config
		data.Server.HTTP.Address = ":8080"
		data.Server.HTTP.DebugEnabled = false
		data.Server.HTTP.MaxMultipartMemory = 20 * 1024 * 1024
		data.Server.HTTP.ForwardedByClientIP = true
		c.set(data)
		write = true
	}
	if info != nil {
		if info.ModTime().UTC().After(c.modeTime) {
			if err = c.read(); err != nil {
				return err
			}
			c.modeTime = info.ModTime().UTC()
		}
	}
	if len(strings.TrimSpace(c.data.Security.Key)) != 32 {
		write = true
		c.data.Security.Key = stringutil.GenerateRandom(32)
	}

	if c.init.Load() {
		c.init.Store(false)
		env, err := c.ctx.DotEnv().Get()
		if err != nil {
			return err
		}
		addr := fmt.Sprintf(":%d", env.HTTP.ServerPort)
		if addr != c.data.Server.HTTP.Address {
			c.data.Server.HTTP.Address = addr
			write = true
		}
		if env.HTTP.ServerDebugEnabled != c.data.Server.HTTP.DebugEnabled {
			c.data.Server.HTTP.DebugEnabled = env.HTTP.ServerDebugEnabled
			write = true
		}
		if env.HTTP.ServerMaxMultipartMemory != c.data.Server.HTTP.MaxMultipartMemory {
			c.data.Server.HTTP.MaxMultipartMemory = env.HTTP.ServerMaxMultipartMemory
			write = true
		}
		if env.HTTP.ServerForwardedByClientIP != c.data.Server.HTTP.ForwardedByClientIP {
			c.data.Server.HTTP.ForwardedByClientIP = env.HTTP.ServerForwardedByClientIP
			write = true
		}
		if env.SQLite != nil {
			if c.data.Database.MariaDB != nil {
				c.data.Database.MariaDB = nil
				write = true
			}
			if c.data.Database.Postgres != nil {
				c.data.Database.Postgres = nil
				write = true
			}
			if c.data.Database.SQLite == nil {
				c.data.Database.SQLite = &types.SqliteDatabaseConfig{}
				write = true
			}
			if env.SQLite.CacheEnabled != c.data.Database.SQLite.CacheEnabled {
				write = true
				c.data.Database.SQLite.CacheEnabled = env.SQLite.CacheEnabled
			}
			if env.SQLite.DebugEnabled != c.data.Database.SQLite.DebugEnabled {
				write = true
				c.data.Database.SQLite.DebugEnabled = env.SQLite.DebugEnabled
			}
			if env.SQLite.EncryptionEnabled != c.data.Database.SQLite.EncryptionEnabled {
				write = true
				c.data.Database.SQLite.EncryptionEnabled = env.SQLite.EncryptionEnabled
			}

		}
		if env.Redis != nil {
			if c.data.Database.Redis == nil {
				c.data.Database.Redis = &types.RedisDatabaseConfig{}
			}
			address := fmt.Sprintf("%s:%d", env.Redis.DatabaseHost, env.Redis.DatabasePort)
			if address != c.data.Database.Redis.Address {
				c.data.Database.Redis.Address = address
				write = true
			}
			username := env.Redis.DatabaseUsername
			if username != c.data.Database.Redis.Username {
				c.data.Database.Redis.Username = username
				write = true
			}
			password := env.Redis.DatabasePassword
			if password != c.data.Database.Redis.Password {
				c.data.Database.Redis.Password = password
				write = true
			}
			databaseName := fmt.Sprintf("%d", env.Redis.DatabaseName)
			if databaseName != c.data.Database.Redis.DatabaseName {
				c.data.Database.Redis.DatabaseName = databaseName
				write = true
			}
			if env.Redis.DebugEnabled != c.data.Database.Redis.DebugEnabled {
				c.data.Database.Redis.DebugEnabled = env.Redis.DebugEnabled
				write = true
			}
		}
		if env.MongoDB != nil {
			if c.data.Database.MongoDB == nil {
				c.data.Database.MongoDB = &types.MongoDBDatabaseConfig{}
			}
			address := fmt.Sprintf("%s:%d", env.MongoDB.DatabaseHost, env.MongoDB.DatabasePort)
			if address != c.data.Database.MongoDB.Address {
				c.data.Database.MongoDB.Address = address
				write = true
			}
			username := env.MongoDB.DatabaseUsername
			if username != c.data.Database.MongoDB.Username {
				c.data.Database.MongoDB.Username = username
				write = true
			}
			password := env.MongoDB.DatabasePassword
			if password != c.data.Database.MongoDB.Password {
				c.data.Database.MongoDB.Password = password
				write = true
			}
			databaseName := env.MongoDB.DatabaseName
			if databaseName != c.data.Database.MongoDB.DatabaseName {
				c.data.Database.MongoDB.DatabaseName = databaseName
				write = true
			}
			if env.MongoDB.DebugEnabled != c.data.Database.MongoDB.DebugEnabled {
				c.data.Database.MongoDB.DebugEnabled = env.MongoDB.DebugEnabled
				write = true
			}
		}
		if env.MariaDB != nil {
			if c.data.Database.SQLite != nil {
				c.data.Database.SQLite = nil
				write = true
			}
			if c.data.Database.MariaDB == nil {
				c.data.Database.MariaDB = &types.MariaDBDatabaseConfig{}
				write = true
			}
			if c.data.Database.Postgres != nil {
				c.data.Database.Postgres = nil
				write = true
			}
			address := fmt.Sprintf("%s:%d", env.MariaDB.DatabaseHost, env.MariaDB.DatabasePort)
			if address != c.data.Database.MariaDB.Address {
				c.data.Database.MariaDB.Address = address
				write = true
			}
			username := env.MariaDB.DatabaseUsername
			if username != c.data.Database.MariaDB.Username {
				c.data.Database.MariaDB.Username = username
				write = true
			}
			password := env.MariaDB.DatabasePassword
			if password != c.data.Database.MariaDB.Password {
				c.data.Database.MariaDB.Password = password
				write = true
			}
			databaseName := env.MariaDB.DatabaseName
			if databaseName != c.data.Database.MariaDB.DatabaseName {
				c.data.Database.MariaDB.DatabaseName = databaseName
				write = true
			}
			if env.MariaDB.CacheEnabled != c.data.Database.MariaDB.CacheEnabled {
				c.data.Database.MariaDB.CacheEnabled = env.MariaDB.CacheEnabled
				write = true
			}
			if env.MariaDB.DebugEnabled != c.data.Database.MariaDB.DebugEnabled {
				c.data.Database.MariaDB.DebugEnabled = env.MariaDB.DebugEnabled
				write = true
			}
		}
		if env.Postgres != nil {
			if c.data.Database.SQLite != nil {
				c.data.Database.SQLite = nil
				write = true
			}
			if c.data.Database.MariaDB != nil {
				c.data.Database.MariaDB = nil
				write = true
			}
			if c.data.Database.Postgres == nil {
				c.data.Database.Postgres = &types.PostgresDatabaseConfig{}
				write = true
			}
			address := fmt.Sprintf("%s:%d", env.Postgres.DatabaseHost, env.Postgres.DatabasePort)
			if address != c.data.Database.Postgres.Address {
				c.data.Database.Postgres.Address = address
				write = true
			}
			username := env.Postgres.DatabaseUsername
			if username != c.data.Database.Postgres.Username {
				c.data.Database.Postgres.Username = username
				write = true
			}
			password := env.Postgres.DatabasePassword
			if password != c.data.Database.Postgres.Password {
				c.data.Database.Postgres.Password = password
				write = true
			}
			databaseName := env.Postgres.DatabaseName
			if databaseName != c.data.Database.Postgres.DatabaseName {
				c.data.Database.Postgres.DatabaseName = databaseName
				write = true
			}
			if env.Postgres.CacheEnabled != c.data.Database.Postgres.CacheEnabled {
				c.data.Database.Postgres.CacheEnabled = env.Postgres.CacheEnabled
				write = true
			}
			if env.Postgres.DebugEnabled != c.data.Database.Postgres.DebugEnabled {
				c.data.Database.Postgres.DebugEnabled = env.Postgres.DebugEnabled
				write = true
			}
		}

	}
	if write {
		if err = c.write(false); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) read() error {
	dir := filepath.Dir(c.filePath)
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("could not create config dir: %s", err)
		}
	}
	if _, err := os.Stat(c.filePath); err == nil {
		bytes, err := os.ReadFile(c.filePath)
		if err != nil {
			return fmt.Errorf("could not load config file: %s", err)
		}
		var data types.Config
		if err = json.Unmarshal(bytes, &data); err != nil {
			return fmt.Errorf("could not decode config file: %s", err)
		}
		c.set(data)
	}
	return nil
}

func (c *config) set(data types.Config) {
	c.data = &data
}

func (c *config) write(callback bool) error {
	c.ctx.Logger().Debug("updating config file...")
	c.callback.Store(callback)
	bytes, err := json.MarshalIndent(c.data, "", "    ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(c.filePath, bytes, os.ModePerm); err != nil {
		return err
	}
	info, err := os.Stat(c.filePath)
	if err != nil {
		return fmt.Errorf("could not write config file: %s", err)
	}
	c.modeTime = info.ModTime().UTC()
	c.ctx.Logger().Debug("config file successfully updated")
	return nil
}

func (c *config) watch() error {
	if c.watcher == nil {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		c.watcher = watcher
		go func() {
			for {
				select {
				case evt, ok := <-c.watcher.Events:
					if !ok {
						return
					}
					if evt.Has(fsnotify.Write) || evt.Has(fsnotify.Create) {
						if c.data != nil && c.callback.Load() && c.callbacks.OnChangeFn != nil {
							c.callbacks.OnChangeFn(*c.data)
						}
					}
				case err, ok := <-c.watcher.Errors:
					if !ok {
						return
					}
					if err != nil {
						c.ctx.Logger().Error(err)
					}
				}
			}
		}()
		if err = c.watcher.Add(c.filePath); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) validate() error {
	if len(c.data.Server.HTTP.Address) == 0 {
		return errors.New("invalid http server address")
	}
	if len(c.data.Security.Key) != 32 {
		return errors.New("invalid security key")
	}
	return nil
}
