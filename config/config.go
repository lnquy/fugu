package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/lnquy/fugu/modules/global"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	Config struct {
		*Runtime
		*Server
	}

	Runtime struct {
		IsDebugging bool   `envconfig:"DEBUGGING" default:"true"`
		LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
	}

	Server struct {
		IP       string        `envconfig:"SERVER_IP" default:"127.0.0.1"`
		Port     string        `envconfig:"SERVER_PORT" default:"3333"`
		TLSKey   string        `envconfig:"SERVER_TLS_KEY"`
		TLSCert  string        `envconfig:"SERVER_TLS_CERT"`
		RTimeout time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"15s"`
		WTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"15s"`
	}
)

func LoadEnvConfig() *Config {
	var cfg Config
	if err := envconfig.Process(global.EnvironmentPrefix, &cfg); err != nil {
		log.Fatalf("config: unable to load config for %T: %s", cfg, err)
	}
	return &cfg
}

func (s *Server) GetFullAddr() string {
	if s.Port == "" {
		return s.IP
	}
	return fmt.Sprintf("%s:%s", s.IP, s.Port)
}
