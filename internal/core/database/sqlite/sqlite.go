package sqlite

import (
	"log"
	"os"
	"path/filepath"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/fileutil"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
)

type database struct {
	ctx  types.IContext
	conn *gorm.DB
}

func New(ctx types.IContext) types.ISQLiteDatabase {
	d := &database{
		ctx: ctx,
	}
	ctx.Hooks().OnStop(d.close)
	ctx.Hooks().OnStart(d.open)
	return d
}

func (d *database) open() error {
	conf, err := d.ctx.Config().Get()
	if err != nil {
		return err
	}
	if conf.Database.SQLite == nil {
		return nil
	}
	path := "data/.database/app.db"
	pathDir := filepath.Dir(path)
	d.ctx.Logger().Debugf("connecting to sqlite database %s...", path)
	if !fileutil.Exists(pathDir) {
		if err = os.MkdirAll(pathDir, os.ModePerm); err != nil {
			return err
		}
	}
	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: dblogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			dblogger.Config{
				Colorful:                  true,
				LogLevel:                  dblogger.Error,
				SlowThreshold:             time.Second * 10,
				ParameterizedQueries:      true,
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		return err
	}
	d.conn = conn
	d.ctx.Logger().Debugf("successfully connected to sqlite database %s", path)
	return nil
}

func (d *database) close() error {
	return nil
}

func (d *database) Conn() (*gorm.DB, bool, error) {
	if d.conn == nil {
		conf, err := d.ctx.Config().Get()
		if err != nil {
			return nil, false, err
		}
		if conf.Database.SQLite == nil {
			return nil, false, nil
		}
		if err = d.open(); err != nil {
			return nil, false, err
		}
	}
	return d.conn, true, nil
}
