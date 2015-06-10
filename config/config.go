// Package config reads configuration information from a file.
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
type AWSConfig struct {
	S3BucketName string
	S3KeyPrefix  string
	AccessKey    string
	SecretKey    string
}
type Config struct {
	ZDConf      ZendeskConfig
	AWSConf     AWSConfig
	MaxFileSize uint64
}

var cfg *Config = nil

// readConfig reads the config from `cfgFileName` and Unmarshals the data
// into a Config struct
func readConfig() *Config {
	cfg = &Config{}
	bytes, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		log.Error("Error while reading file %v: %+v", cfgFileName, err)
		return &Config{}
	}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		log.Error("Error unmarshaling config file json: %+v", err)
	}
	return cfg
}

// GetConfig determines whether the config singleton is already initialized.
// If it is initialized, it is simply returned. Else, a call to readConfig is made.
func GetConfig() Config {
	if cfg == nil {
		cfg = readConfig()
	}
	return *cfg
}
