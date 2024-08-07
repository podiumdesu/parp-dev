package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	PrivateKeyFilePath string `json:"privateKeyFilePath"`
	BcWsEndpoint       string `json:"bcWsEndpoint"`
	BcRpcEndpoint      string `json:"bcRpcEndpoint"`
	ServerEndpoint     string `json:"serverEndpoint"`
	ContractAddress    string `json:"contractAddress"`
}

func LoadConfig(filename string) (*Config, error) {
	// Read configuration file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse configuration from JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
