package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HTTPConfig struct {
	Address string        `yaml:"address" env:"API_ADDRESS" env-default:":8081"`
	Timeout time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
}

type DBConfig struct {
	Address         string        `yaml:"address" env:"DB_ADDRESS"`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"20"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"10"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"1h"`
}

type ITMOConfig struct {
	AdapterURL string        `yaml:"adapter_url" env:"ITMO_ADAPTER_URL" env-default:"http://127.0.0.1:35601"`
	APIToken   string        `yaml:"api_token" env:"ITMO_ADAPTER_API_TOKEN"`
	Timeout    time.Duration `yaml:"timeout" env:"ITMO_ADAPTER_TIMEOUT" env-default:"10s"`
}

type NATSConfig struct {
	URL      string `yaml:"url" env:"NATS_URL" env-default:"nats://localhost:4222"`
	Stream   string `yaml:"stream" env:"NATS_STREAM" env-default:"SCHEDULER"`
	Consumer string `yaml:"consumer" env:"NATS_CONSUMER" env-default:"SCHEDULER_SYNC"`
}

type Config struct {
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`

	HTTP HTTPConfig `yaml:"api_server"`
	DB   DBConfig   `yaml:"database"`
	ITMO ITMOConfig `yaml:"itmo"`
	NATS NATSConfig `yaml:"nats"`
}

func MustLoad(configPath string) Config {
	var cfg Config

	_ = godotenv.Load()

	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			var pe *os.PathError
			if !errors.As(err, &pe) {
				log.Fatalf("cannot read config %q: %s", configPath, err)
			}
		}
	}

	// Environment variables override values from config.yaml
	if err := cleanenv.UpdateEnv(&cfg); err != nil {
		log.Fatalf("cannot read env: %s", err)
	}

	if cfg.DB.Address == "" {
		log.Fatal("DB_ADDRESS is required")
	}

	if cfg.ITMO.APIToken == "" {
		log.Fatal("ITMO_ADAPTER_API_TOKEN is required")
	}

	return cfg
}
