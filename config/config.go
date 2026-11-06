package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type RentalConfig struct {
	ServiceName     string
	ServiceBootTime time.Time
	Profile         string `envconfig:"PROFILE" default:"local"`
	ContextTimeout  int    `envconfig:"CONTEXT_TIMEOUT" default:"5" desc:"context_timeout"`
	GinMode         string `envconfig:"GIN_MODE" default:"debug"`
	RootPath        string `envconfig:"LUXE_ROOT_PATH" default:"." yaml:"rootPath"` // for the console
	// RootPath string `envconfig:"LUXE_ROOT_PATH" default:"C:/Program Files/LuxeQuoteSystem" yaml:"rootPath"` // for the service

	// Logging
	LogOutput     string `envconfig:"LOG_OUTPUT" default:"debug.log" yaml:"logOutput"`
	LogLevel      string `envconfig:"LOG_LEVEL" default:"debug" yaml:"logLevel"`
	LogStructured bool   `envconfig:"LOG_STRUCTURED" default:"true" yaml:"logStructured"`

	// DB Setup
	DBSetup bool `envconfig:"DB_SETUP" default:"true" yaml:"dbSetup"`

	// Web Servers
	ServerAddress         string `envconfig:"SERVER_ADDRESS" default:":9999" desc:"localhost:9999"`
	InternalServerAddress string `envconfig:"INTERNAL_SERVER_ADDRESS" default:":9998" desc:"localhost:9998"`

	// Postgresql
	PostgresHost     string `envconfig:"POSTGRES_HOST" yaml:"postgresHost" default:"aws-0-eu-central-1.pooler.supabase.com" required:"true" desc:"127.0.0.1"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" yaml:"postgresPort" default:"5432" required:"true" desc:"5432"`
	PostgresUser     string `envconfig:"POSTGRES_USER" yaml:"postgresUser" default:"postgres.piemowbzscfgslggpjgp" required:"true" desc:"user"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" yaml:"postgresPassword" default:"kS9ensPK8d4Tomf6" required:"true" desc:"password"`
	PostgresDatabase string `envconfig:"POSTGRES_DB" yaml:"postgresDb" default:"postgres" required:"false" desc:"rental_db"`
	PostgresSsl      string `envconfig:"POSTGRES_SSL" yaml:"postgresSsl" required:"false" default:"prefer" desc:"ssl"`
}

const serviceName = "luxe-backend"

var cfg *RentalConfig

func GetInstance() *RentalConfig {
	return cfg
}

func InitConfigFromEnvYaml() (*RentalConfig, error) {
	if cfg == nil {
		_, err := InitConfigFromEnv()
		if err != nil {
			return cfg, err
		}

		cfgFileName := fmt.Sprintf("%s/config.%s.yml", cfg.RootPath, cfg.Profile)
		file, err := os.Open(cfgFileName)
		if err != nil {
			return cfg, fmt.Errorf("failed to read profile config from %s : %s", cfgFileName, err.Error())
		}

		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(cfg)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse the config from %s : %s", cfgFileName, err.Error())
		}
	}

	return cfg, nil
}

func InitConfigFromEnv() (*RentalConfig, error) {
	if cfg == nil {
		cfg = new(RentalConfig)
		cfg.ServiceBootTime = time.Now()
		cfg.ServiceName = serviceName

		err := envconfig.Process(cfg.ServiceName, cfg)
		if err != nil {
			return cfg, fmt.Errorf("failed to read environment variables configuration : %s", err.Error())
		}
	}

	return cfg, nil
}
