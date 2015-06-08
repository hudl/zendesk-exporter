package zendesk

import (
	"github.com/adamar/ZeGo/zego"
	"github.com/hudl/zendesk-exporter/config"
)

type ZDPoller struct {
	Config    config.ZendeskConfig
	Auth      zego.Auth
	StartTime string
}

func (p *ZDPoller) Poll() {

}
