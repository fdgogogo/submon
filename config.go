package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const comments = `# submon 默认配置文件
# 使用 submon echo_example_config > config.yaml 来保存`

type WatchDir struct {
	Path       string
	NoFullScan bool
	Recursive  bool
}

type WatchConfig struct {
	Dirs     []WatchDir
	MaxRetry int
}

type Config struct {
	Watch     WatchConfig
	Lang      string
	Workers   int
	LogFormat string
	LogLevel  string
	Debug     bool
}

func NewConfig() (config Config) {
	config.Lang = "chn"
	config.Workers = 4
	config.LogFormat = "color"
	config.LogLevel = "INFO"
	config.Debug = false
	config.Watch.MaxRetry = 3
	return
}

func ExampleConfig() (y string) {
	c := NewConfig()
	c.Watch.Dirs = []WatchDir{
		WatchDir{Path: "path/to/watch", NoFullScan: false, Recursive: false},
		WatchDir{Path: "another/path/to/watch", NoFullScan: true, Recursive: true}}
	d, _ := yaml.Marshal(&c)
	y = comments + "\n" + string(d)
	return
}

func ReadConfigFile(configFilePath string) (config Config) {
	logger.Info("Using config file: " + configFilePath)

	var configData []byte
	path, err := homedir.Expand(configFilePath)
	if err != nil {
		logger.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) || os.IsPermission(err) {
		// 文件是否存在
		logger.Info("Config file does not exists or not readable, using default config.")
		config = NewConfig()
		return
	} else {
		// 读取文件失败
		configData, err = ioutil.ReadFile(path)
		if err != nil {
			logger.Fatal(err)
		}
	}

	yaml_err := yaml.Unmarshal(configData, &config)
	if yaml_err != nil {
		logger.Fatalf("Malformed config file %s", configFilePath)
	}
	return
}

func PrintDefaultConfig() {
	fmt.Println(ExampleConfig())
}
