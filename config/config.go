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
	cfg = &ZendeskConfig{}
	bytes, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		log.Error("Error while reading file %v: %+v", cfgFileName, err)
		return &ZendeskConfig{}
	}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		log.Error("Error unmarshaling config file json: %+v", err)
	}
	return cfg
}

func GetZDConfig() ZendeskConfig {
	if cfg == nil {
		cfg = readConfig()
	}
	return *cfg
}
