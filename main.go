package main

import (
	"os"

	"github.com/hudl/zendesk-exporter/config"
	"github.com/hudl/zendesk-exporter/logging"
	"github.com/hudl/zendesk-exporter/zendesk"

	"github.com/adamar/ZeGo/zego"
	golog "github.com/op/go-logging"
)

var log = golog.MustGetLogger("main")

func main() {
	logging.Configure()
	log.Info("Starting...")

	startTime := os.Args[1]
	if startTime == "" {
		log.Error("No start time specified. Exiting.")
		return
	}
	log.Info("Using startTime: %s", startTime)

	cfg := config.GetZDConfig()
	log.Info("Config: %+v", cfg)
	auth := zego.Auth{cfg.Username, cfg.Password, cfg.BaseUrl}
	poller := zendesk.Poller{
		auth,
		startTime,
	}
	poller.Poll()
}
