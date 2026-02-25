package types

type DotEnv struct {
	Log struct {
		Level string `env:"LEVEL" envDefault:"info"`
	} `envPrefix:"LOG_"`
	HTTP struct {
		ServerPort                uint32 `env:"SERVER_PORT" envDefault:"8080"`
		ServerDebugEnabled        bool   `env:"SERVER_DEBUG_ENABLED" envDefault:"false"`
		ServerMaxMultipartMemory  uint64 `env:"SERVER_MAX_MULTIPART_MEMORY" envDefault:"20971520"`
		ServerForwardedByClientIP bool   `env:"SERVER_FORWARDED_BY_CLIENT_IP" envDefault:"true"`
	} `envPrefix:"HTTP_"`
	SQLite *struct {
		CacheEnabled      bool `env:"CACHE_ENABLED" envDefault:"false"`
		DebugEnabled      bool `env:"DEBUG_ENABLED" envDefault:"false"`
		EncryptionEnabled bool `env:"ENCRYPTION_ENABLED" envDefault:"false"`
	} `envPrefix:"SQLITE_"`
	Redis *struct {
		Enabled          bool   `env:"ENABLED" envDefault:"false"`
		DebugEnabled     bool   `env:"DEBUG_ENABLED" envDefault:"false"`
		DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
		DatabasePort     uint32 `env:"DATABASE_PORT" envDefault:"6379"`
		DatabaseName     int    `env:"DATABASE_NAME" envDefault:"0"`
		DatabasePrefix   string `env:"DATABASE_PREFIX" envDefault:"app:"`
		DatabaseUsername string `env:"DATABASE_USERNAME" envDefault:""`
		DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:""`
	} `envPrefix:"REDIS_"`
	MariaDB *struct {
		Enabled          bool   `env:"ENABLED" envDefault:"false"`
		CacheEnabled     bool   `env:"CACHE_ENABLED" envDefault:"false"`
		DebugEnabled     bool   `env:"DEBUG_ENABLED" envDefault:"false"`
		DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
		DatabasePort     uint32 `env:"DATABASE_PORT" envDefault:"3306"`
		DatabaseName     string `env:"DATABASE_NAME" envDefault:"app"`
		DatabaseUsername string `env:"DATABASE_USERNAME" envDefault:"root"`
		DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:""`
	} `envPrefix:"MARIADB_"`
	MongoDB *struct {
		Enabled          bool   `env:"ENABLED" envDefault:"false"`
		DebugEnabled     bool   `env:"DEBUG_ENABLED" envDefault:"false"`
		DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
		DatabasePort     uint32 `env:"DATABASE_PORT" envDefault:"27017"`
		DatabaseName     string `env:"DATABASE_NAME" envDefault:"app"`
		DatabaseUsername string `env:"DATABASE_USERNAME" envDefault:"admin"`
		DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:""`
	} `envPrefix:"MONGODB_"`
	Postgres *struct {
		SSLMode          bool   `env:"SSL_MODE" envDefault:"false"`
		TimeZone         string `env:"TIMEZONE" envDefault:"UTC"`
		CacheEnabled     bool   `env:"CACHE_ENABLED" envDefault:"false"`
		DebugEnabled     bool   `env:"DEBUG_ENABLED" envDefault:"false"`
		DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
		DatabasePort     uint32 `env:"DATABASE_PORT" envDefault:"5432"`
		DatabaseName     string `env:"DATABASE_NAME" envDefault:"app"`
		DatabaseUsername string `env:"DATABASE_USERNAME" envDefault:"postgres"`
		DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:""`
	} `envPrefix:"POSTGRES_"`
}

type IDotEnv interface {
	Get() (*DotEnv, error)
}
