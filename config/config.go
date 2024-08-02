package config

import (
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	Env     string `env:"ENV" envDefault:"dev"`
	OpenAi  OpenAi
	Grpc    Grpc
	MsGraph MsGraph
}

type OpenAi struct {
	Model  string `env:"model,notEmpty" envDefault:"whisper-1"`
	Secret string `env:"secret,notEmpty"`
}

type Grpc struct {
	Port      int           `env:"port" envDefault:"44045"`
	Timeout   time.Duration `env:"timeout" envDefault:"5s"`
	InputFile string        `env:"input_file,notEmpty" envDefault:"audio/received.wav"`
}

type MsGraph struct {
	Token string `env:"token,notEmpty"`
}

func Load() *Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic("Could not read config. " + err.Error())
	}

	return &cfg
}
