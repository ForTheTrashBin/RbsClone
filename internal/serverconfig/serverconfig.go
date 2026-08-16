package serverconfig

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

//-----------------------------------------------------------------------------
// Configuration for the postgres database
//-----------------------------------------------------------------------------

type DBConfig struct {
	User     string `env:"USER,notEmpty"`          // The User's name
	Password string `env:"PASSWORD,notEmpty"`      // The Password of the user
	Host     string `env:"HOST,notEmpty"`          // The database host
	Port     uint16 `env:"PORT" envDefault:"5432"` // The port the database is running on (defaults to 5432)
	Database string `env:"DATABASE,notEmpty"`      // The name of the database to connect to
}

//-----------------------------------------------------------------------------
// Configuration for tzhe routers and http(s) servers
//-----------------------------------------------------------------------------

type RTConfig struct {
	HTTPPort  uint16 `env:"HTTP_PORT" envDefault:"8080"`
	HTTPSPort uint16 `env:"HTTPS_PORT" envDefault:"8443"`
}

//-----------------------------------------------------------------------------

type Config struct {
	DB       DBConfig   `envPrefix:"DB_"`
	RT       RTConfig   `envPrefix:"RT_"`
	LogLevel slog.Level `env:"LOGLEVEL" envDefault:"info"`
}

//-----------------------------------------------------------------------------
// Administrativ functions
//-----------------------------------------------------------------------------

func (cfg *Config) InitGoDotEnv() error {

	return godotenv.Load()
}

func (cfg *Config) InitCaarlos0() error {

	return env.Parse(cfg)
}

//-----------------------------------------------------------------------------
// Access functions
//-----------------------------------------------------------------------------

func (cfg Config) GetHTTPPort() uint16 {

	return cfg.RT.HTTPPort
}

func (cfg Config) GetHTTPSPort() uint16 {

	return cfg.RT.HTTPSPort
}
