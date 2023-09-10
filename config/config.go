package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
	"time"
)

type (
	Config struct {
		HTTPServer `yaml:"http_server"`
		ZapLogger  `yaml:"zap_logger"`
		PostgresDB `yaml:"postgres"`
	}

	HTTPServer struct {
		Port            string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	}

	ZapLogger struct {
		Level           string   `yaml:"level"`
		Encoding        string   `yaml:"encoding"`
		OutputPath      []string `yaml:"output_path"`
		ErrorOutputPath []string `yaml:"error_output_path"`
	}

	PostgresDB struct {
		MaxPoolSize  int           `yaml:"max_pool_size"`
		ConnTimeout  time.Duration `yaml:"conn_timeout"`
		ConnAttempts int           `yaml:"conn_attempts"`
		URL          string        `env-required:"true" env:"PG_URL"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
