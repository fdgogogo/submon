package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Watch string
}

func ReadConfig(configFilePath string) (config Config) {
	logger.Println("Using config file: " + configFilePath)
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	yaml_err := yaml.Unmarshal(configData, &config)
	if yaml_err != nil {
		panic(yaml_err)
	}
	return
}
