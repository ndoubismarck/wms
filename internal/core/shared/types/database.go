package types

import (
	"gorm.io/gorm"
)

type Dialect string

const (
	MySQLDialect      Dialect = "mysql"
	SQLiteDialect             = "sqlite"
	PostgreSQLDialect         = "postgresql"
)

// todo add redis and mongodb

type IDatabase interface {
	MySQL() IMySQLDatabase
	SQLite() ISQLiteDatabase
	PostgreSQL() IPostgreSQLDatabase
}

type IMySQLDatabase interface {
	Conn() (*gorm.DB, bool, error)
}

type ISQLiteDatabase interface {
	Conn() (*gorm.DB, bool, error)
}

type IPostgreSQLDatabase interface {
	Conn() (*gorm.DB, bool, error)
}
