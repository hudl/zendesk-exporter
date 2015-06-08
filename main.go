package main

import (
	"os"

	"github.com/adamar/ZeGo/zego"
	"github.com/hudl/zendesk-exporter/config"
	"github.com/hudl/zendesk-exporter/zendesk"
)

func main() {
	startTime := os.Args[1]
	if startTime == "" {
		return
	}

	cfg := config.GetZDConfig()
	auth := zego.Auth{cfg.Username, cfg.Password, cfg.BaseUrl}
	poller := zendesk.ZDPoller{
		cfg,
		auth,
		startTime,
	}
	poller.Poll()
}
