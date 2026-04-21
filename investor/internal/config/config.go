package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/alekparkhomenko/investor/investor/internal/config/env"
)

var appConfig *Config

type Config struct {
	App    AppSettings
	Logger LoggerSettings
	DB     DBSettings
}

// DBSettings holds database configuration.
type DBSettings struct {
	URL         string
	MaxOpenConns int
	MaxIdleConns int
}

func Load() error {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	appCfg, err := env.NewAppConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &Config{
		App:    appCfg,
		Logger: loggerCfg,
		DB: DBSettings{
			URL:         appCfg.DatabaseURL(),
			MaxOpenConns: 25,
			MaxIdleConns: 10,
		},
	}

	return nil
}

func AppConfig() *Config {
	return appConfig
}
