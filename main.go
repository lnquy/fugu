package main

import (
	"github.com/lnquy/fugu/config"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/fugu/server"
)

func main() {
	cfg := config.LoadEnvConfig()

	lvl, _:= log.ParseLevel(cfg.Runtime.LogLevel)
	log.SetLevel(lvl)

	server.Serve(cfg)
}

