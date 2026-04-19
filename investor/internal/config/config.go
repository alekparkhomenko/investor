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
	}

	return nil
}

func AppConfig() *Config {
	return appConfig
}
