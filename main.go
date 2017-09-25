package main

import (
	"github.com/lnquy/fugu/config"
	"github.com/lnquy/fugu/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadEnvConfig()
	lvl, _ := log.ParseLevel(cfg.Runtime.LogLevel) // Default debug
	log.SetLevel(lvl)

	server.Serve(cfg)
}
