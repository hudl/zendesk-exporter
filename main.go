package main

import (
	"os"

	"github.com/hudl/zendesk-exporter/config"
	"github.com/hudl/zendesk-exporter/logging"
	"github.com/hudl/zendesk-exporter/ticketwriter"
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

	cfg := config.GetConfig()
	aconf := cfg.AWSConf
	tickWrt := ticketwriter.New(cfg.MaxFileSize, "ticks_", aconf.S3BucketName, aconf.S3KeyPrefix, aconf.AccessKey, aconf.SecretKey)
	auth := zego.Auth{cfg.ZDConf.Username, cfg.ZDConf.Password, cfg.ZDConf.BaseUrl}
	poller := zendesk.Poller{
		auth,
		startTime,
		tickWrt,
		0,
	}
	poller.Poll()
}
