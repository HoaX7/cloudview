package main

import (
	"cloudview/agents/exporter/core/iox"
	"fmt"

	"gopkg.in/yaml.v3"
)

const configFile = "./config.yaml"

type StaticConfig struct {
	Resource  string `yaml:"resource" json:"resource"`
	Service   string `yaml:"service" json:"service"`
	PublicKey string `yaml:"publicKey" json:"publicKey"`
	Reporting bool   `yaml:"reporting" json:"reporting" default:"true"`
	Config    struct {
		InstanceID string `yaml:"instanceId" json:"instanceId"`
		Region     string `yaml:"region" json:"region"`
		Unit       string `yaml:"unit" json:"unit"`
		Period     struct {
			Cpu  int `yaml:"cpu" json:"cpu"`
			Ram  int `yaml:"ram" json:"ram"`
			Disk int `yaml:"disk" json:"disk"`
		} `yaml:"period" json:"period"`
	} `yaml:"config" json:"config"`
}

var staticConfig StaticConfig

func initialize() {
	// load config from 'config.yaml' file
	file, err := iox.ReadFile(configFile)
	if err != nil {
		err := fmt.Errorf("Error reading 'config.yaml' file make sure it is a valid yaml file %w", err)
		panic(err)
	}

	if err := yaml.Unmarshal(file, &staticConfig); err != nil {
		err := fmt.Errorf("Error reading 'config.yaml' file make sure it is a valid yaml file %w", err)
		panic(err)
	}
	if staticConfig.PublicKey == "" {
		errMsg := "Invalid Public Key. Code 403. Exporter will not report metrics."
		panic(errMsg)
	}
	return
}
