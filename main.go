package main

import (
	"os"

	"github.com/hudl/zendesk-exporter/config"
	"github.com/hudl/zendesk-exporter/logging"
	"github.com/hudl/zendesk-exporter/ticketwriter"
	"github.com/hudl/zendesk-exporter/zendesk"

	"github.com/hudl/ZeGo/zego"
	golog "github.com/op/go-logging"
)

var log = golog.MustGetLogger("main")

func main() {
	logging.Configure()
	log.Info("Starting...")
	if len(os.Args) != 2 {
		log.Error("Usage: %s <starttime>", os.Args[0])
		return
	}
	startTime := os.Args[1]
	if startTime == "" {
		log.Error("No start time specified. Exiting.")
		return
	}
	log.Info("Using startTime: %s", startTime)

	cfg := config.GetConfig()
	aconf := cfg.AWSConf
	tickWrt := ticketwriter.New(cfg.MaxFileSize, "ticks_", aconf.S3BucketName,
		aconf.S3KeyPrefix, aconf.AccessKey, aconf.SecretKey, aconf.KinesisStream,
		cfg.ICConf.App, cfg.ICConf.APIKey)
	auth := zego.Auth{cfg.ZDConf.Username, cfg.ZDConf.Password, cfg.ZDConf.BaseUrl}
	poller := zendesk.Poller{
		auth,
		startTime,
		tickWrt,
		0,
	}
	poller.Poll()
}
