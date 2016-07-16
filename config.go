package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"fmt"
	"github.com/mitchellh/go-homedir"
)

const comments = `# submon 默认配置文件
# 使用 submon echo_example_config > config.yaml 来保存`

type WatchConfig struct {
	Path       string
	MaxRetry   int
	NoFullScan bool
	Recursive  bool
}

type Config struct {
	Watch     []WatchConfig
	Lang      string
	LogFormat string
	LogLevel  string
	Debug     bool
}

func NewConfig() (config Config) {
	config.Lang = "chn"
	config.LogFormat = "color"
	config.LogLevel = "INFO"
	config.Debug = false
	return
}

func ExampleConfig() (y string) {
	c := NewConfig()
	c.Watch = []WatchConfig{
		WatchConfig{Path:"path/to/watch", MaxRetry: 3, NoFullScan:false, Recursive:false},
		WatchConfig{Path:"another/path/to/watch", MaxRetry:3, NoFullScan:true, Recursive:true}}
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
