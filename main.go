package main

import (
	"github.com/lnquy/fugu/config"
	"github.com/lnquy/fugu/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadEnvConfig()

	lvl, _ := log.ParseLevel(cfg.Runtime.LogLevel)
	log.SetLevel(lvl)

	server.Serve(cfg)
}
