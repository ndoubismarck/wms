package database

import (
	"server/internal/core/database/mysql"
	"server/internal/core/database/postgresql"
	"server/internal/core/database/sqlite"
	"server/internal/core/shared/types"
)

type database struct {
	mysql      types.IMySQLDatabase
	sqlite     types.ISQLiteDatabase
	postgresql types.IPostgreSQLDatabase
}

func New(ctx types.IContext) types.IDatabase {
	return &database{
		mysql:      mysql.New(ctx),
		sqlite:     sqlite.New(ctx),
		postgresql: postgresql.New(ctx),
	}
}

func (d *database) MySQL() types.IMySQLDatabase {
	return d.mysql
}

func (d *database) SQLite() types.ISQLiteDatabase {
	return d.sqlite
}

func (d *database) PostgreSQL() types.IPostgreSQLDatabase {
	return d.postgresql
}
