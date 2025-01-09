package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the configuration for the game server.
type NetWork struct {
	ClientService string `yaml:"ClientService"`
}

type Router struct {
	RouterService string `yaml:"RouterService"`
}

type ServerInfo struct {
	ServerID uint32 `yaml:"ServerID"`
}

type Config struct {
	NetWork    NetWork    `yaml:"Network"`
	Router     Router     `yaml:"Router"`
	ServerInfo ServerInfo `yaml:"ServerInfo"`
}

var (
	Conf Config
)

func init() {
	// Load the configuration file.
	fmt.Println("config init")
	LoadConfig()
}

func LoadConfig() {
	// Load the configuration file.
	configFile := "../conf/service.yaml"

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Failed to open YAML file: %v", err)
	}
	defer file.Close()

	// Decode the configuration file.
	var conf Config
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("Failed to decode YAML file: %v", err)
	}

	Conf = conf
	fmt.Println(Conf)
}
