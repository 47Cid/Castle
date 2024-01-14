package config

import (
	"encoding/json"
	"os"
)

var config *Config

type Config struct {
	ListenPort         string            `json:"listenPort"`
	RemoteServerAddr   string            `json:"remoteServerAddr"`
	ProxyLogFileName   string            `json:"proxyLogFileName"`
	WAFLogFileName     string            `json:"wafLogFileName"`
	DestinationToLabel map[string]string `json:"destinationToLabel"`
}

func GetListenPort() string {
	return config.ListenPort
}

func GetRemoteServerAddr() string {
	return config.RemoteServerAddr
}

func GetProxyLogFile() string {
	return config.ProxyLogFileName
}

func GetWAFLogFile() string {
	return config.WAFLogFileName
}

func GetLabel(destination string) string {
	return config.DestinationToLabel[destination]
}

func LoadConfig(filename string) error {
	// Open the config file
	configFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// Decode the config file
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}
