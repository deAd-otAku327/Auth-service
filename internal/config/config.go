package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server `yaml:"server"`
	DBConn `yaml:"db-conn"`
}

type Server struct {
	Host                 string        `yaml:"host" env:"HOST"`
	Port                 string        `yaml:"port" env:"PORT"`
	AccessTokenSecretKey string        `env:"ACCESS_TOKEN_SECRET" env-required:"true"`
	AccessTokenExpire    time.Duration `yaml:"access_token_expire" env-default:"12h"`
	RefreshTokenExpire   time.Duration `yaml:"refresh_token_expire" env-default:"48h"`
	LogLevel             string        `yaml:"log_level" env-default:"info"`
	AsyncHashingLimit    int           `yaml:"async_hashing_limit" env-default:"10"`
}

type DBConn struct {
	URL          string `env:"DB_URL" env-required:"true"`
	MaxOpenConns int    `yaml:"max_open_conns" env-default:"15"`
}

func New(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
