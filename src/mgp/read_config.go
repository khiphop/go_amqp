package mgp

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


type Env struct {
	Env         string `yaml:"env"`
	ServiceRole string `yaml:"service_role"`
}


// Cfg yaml file map，capitalize
type Cfg struct {
	HttpPort int `yaml:"http_port"`
	Amqp     AmqpCfg
}

type AmqpCfg struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	VHost        string `yaml:"vhost"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	QueuePrefix  string `yaml:"queue_prefix"`
	QueueCount   int    `yaml:"queue_count"`
	TransferUrl  string `yaml:"transfer_url"`
	QueueStartNo int    `yaml:"queue_start_no"`
}

const (
	configFilePath = "./config/config.yaml"
	envFilePath    = "./config/env.yaml"
)

func ReadConfig() Cfg {
	var st Cfg

	config, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Print(err)
	}

	err1 := yaml.Unmarshal(config, &st)
	if err1 != nil {
		fmt.Println("error")
	}

	return st
}

func ReadEnv() Env {
	var st Env

	config, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		fmt.Print(err)
	}

	err1 := yaml.Unmarshal(config, &st)
	if err1 != nil {
		fmt.Println("error")
	}

	return st
}


