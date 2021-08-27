package config

type Config struct {
	PostgresDB PsqlDatabaseConfig
	App        AppConfig
}

type PsqlDatabaseConfig struct {
	Host     string `envconfig:"WALLET_DB_HOST"`
	Name     string `envconfig:"WALLET_DB_NAME"`
	Dialect  string `envconfig:"WALLET_DB_DIALECT"`
	User     string `envconfig:"WALLET_DB_USER"`
	Password string `envconfig:"WALLET_DB_PASSWORD"`
	Port     string `envconfig:"WALLET_DB_PORT"`
	SSLMode  string `envconfig:"WALLET_DB_SSL_MODE"`
}

type AppConfig struct {
	Port string `envconfig:"WALLET_PORT"`
}
