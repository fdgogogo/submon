package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Watch string
}

const default_config = `# shooter-subtitle-worker 默认配置文件
# 使用 shooter-subtitle-worker echo_example_config > config.yaml 来保存
watch: "path/to/watch"
`

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

func PrintDefaultConfig() {
	logger.Println(default_config)
}
