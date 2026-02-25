package postgresql

import (
	"server/internal/core/shared/types"

	"gorm.io/gorm"
)

type database struct {
	ctx  types.IContext
	conn *gorm.DB
}

func New(ctx types.IContext) types.IPostgreSQLDatabase {
	d := &database{
		ctx: ctx,
	}
	ctx.Hooks().OnStop(d.close)
	ctx.Hooks().OnStart(d.open)
	return d
}

func (d *database) open() error {
	return nil
}

func (d *database) close() error {
	return nil
}

func (d *database) Conn() (*gorm.DB, bool, error) {
	if d.conn == nil {
		return nil, false, nil
	}
	return d.conn, true, nil
}
