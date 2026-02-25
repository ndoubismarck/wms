package types

type Config struct {
	Server struct {
		HTTP struct {
			Address             string `json:"address"`
			DebugEnabled        bool   `json:"debug_enabled"`
			MaxMultipartMemory  uint64 `json:"max_multipart_memory"`
			ForwardedByClientIP bool   `json:"forwarded_by_client_ip"`
		} `json:"http"`
	} `json:"server"`
	Security struct {
		Key string `json:"key"`
	} `json:"security"`
	Database struct {
		Redis    *RedisDatabaseConfig    `json:"redis,omitempty"`
		SQLite   *SqliteDatabaseConfig   `json:"sqlite,omitempty"`
		MariaDB  *MariaDBDatabaseConfig  `json:"mariadb,omitempty"`
		MongoDB  *MongoDBDatabaseConfig  `json:"mongodb,omitempty"`
		Postgres *PostgresDatabaseConfig `json:"postgres,omitempty"`
	} `json:"database"`
	Services struct {
		Ebay *struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"ebay,omitempty"`
	} `json:"services"`
}

type ConfigCallbacks struct {
	OnChangeFn func(config Config)
}

type IConfig interface {
	Get() (*Config, error)
	Set(config *Config) error
}

type RedisDatabaseConfig struct {
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	DebugEnabled bool   `json:"debug_enabled"`
}

type SqliteDatabaseConfig struct {
	CacheEnabled      bool `json:"cache_enabled"`
	DebugEnabled      bool `json:"debug_enabled"`
	EncryptionEnabled bool `json:"encryption_enabled"`
}

type MongoDBDatabaseConfig struct {
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	CacheEnabled bool   `json:"cache_enabled"`
	DebugEnabled bool   `json:"debug_enabled"`
}
type MariaDBDatabaseConfig struct {
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	CacheEnabled bool   `json:"cache_enabled"`
	DebugEnabled bool   `json:"debug_enabled"`
}

type PostgresDatabaseConfig struct {
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	CacheEnabled bool   `json:"cache_enabled"`
	DebugEnabled bool   `json:"debug_enabled"`
}
