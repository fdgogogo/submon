package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Watch string
	Lang  string
}

const defaultConfig = `# shooter-subtitle-worker 默认配置文件
# 使用 shooter-subtitle-worker echo_example_config > config.yaml 来保存
watch: path/to/watch
lang: chn
`

func ReadConfig(configFilePath string) (config Config) {
	logger.Println("Using config file: " + configFilePath)

	var configData []byte
	var err error

	if configFilePath == "" {
		logger.Println("No config file specified, using default config.")
		configData = []byte(defaultConfig)
	} else {
		configData, err = ioutil.ReadFile(configFilePath)
		if err != nil {
			if err == os.ErrNotExist {
				logger.Println("Config file does not exists, using default config.")
				configData = []byte(defaultConfig)
			} else {
				panic(err)
			}
		}
	}
	yaml_err := yaml.Unmarshal(configData, &config)
	if yaml_err != nil {
		panic(yaml_err)
	}
	return
}

func PrintDefaultConfig() {
	logger.Println(defaultConfig)
}