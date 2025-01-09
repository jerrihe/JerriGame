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
	RouterService map[string]string `yaml:"RouterService"`
}

type RedisInfo struct {
	Addr        string `yaml:"Addr"`
	Port        int    `yaml:"Port"`
	Password    string `yaml:"Password"`
	DB          int    `yaml:"DB"`
	MaxIdle     int    `yaml:"MaxIdle"`
	MaxActive   int    `yaml:"MaxActive"`
	IdleTimeout int    `yaml:"IdleTimeout"`
	ConnTimeout int    `yaml:"ConnTimeout"`
	ReadTimeout int    `yaml:"ReadTimeout"`
}

type ServerInfo struct {
	ServerID int32 `yaml:"ServerID"`
}

type Config struct {
	NetWork         NetWork             `yaml:"Network"`
	RelationService map[string][]string `yaml:"RelationService"`
	ServerInfo      ServerInfo          `yaml:"ServerInfo"`
	Redis           RedisInfo           `yaml:"Redis"`
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

func GetRelationService() map[string][]string {
	return Conf.RelationService
}
