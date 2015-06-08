package config

import (
	"encoding/json"
	"io/ioutil"
)

const cfgFileName = "./zendesk-exporter.json"

type ZendeskConfig struct {
	BaseUrl  string
	Username string
	Password string
}

var cfg *ZendeskConfig = nil

func readConfig() *ZendeskConfig {
	cfg = ZendeskConfig{}
	bytes, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return &ZendeskConfig{}
	}
	json.Unmarshal(bytes, &cfg)
	return &cfg
}

func GetZDConfig() ZendeskConfig {
	if cfg == nil {
		cfg = readConfig()
	}
	return *cfg
}
