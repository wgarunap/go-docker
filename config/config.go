package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Conf struct {
	DebugPort   int `yaml:"debug_port"`
	MetricsPort int `yaml:"metrics_port"`
}

var Config Conf

func init() {
	b, err := ioutil.ReadFile("config/config.yaml") // just pass the file name
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.UnmarshalStrict(b, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
