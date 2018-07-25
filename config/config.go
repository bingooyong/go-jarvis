package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"github.com/sirupsen/logrus"
)

var config *Config

func init() {
	path, _ := os.Getwd()

	config = &Config{
		Server: Server{
			Port: 8080,
		},
		Deploy: Deploy{
			Path: path,
		},
		LogPath: "",
	}
}

type Server struct {
	Port int
}

type Deploy struct {
	Path string
}

type Config struct {
	Server  Server
	Deploy  Deploy
	LogPath string `yaml:"log-path"`
}

func Instance() *Config {
	return config
}

func (config *Config) Load(path string) {
	cBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal([]byte(cBytes), config)
	if err != nil {
		panic(err.Error())
	}

	config.ConfigLog()
}

func (config *Config) ConfigLog() {
	logFile := config.LogPath
	if logFile == "" {
		return
	}
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0777)
	if err == nil {
		logrus.StandardLogger().Out = file
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}
