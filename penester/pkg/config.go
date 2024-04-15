package pkg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	ConfPath  string
	Website   string
	Pipeline  PipelineStep
	OutputDir string
}

type PipelineStep struct {
	PortScan *PortScan `yaml:"portScan"`
	DirScan  *DirScan  `yaml:"dirScan"`
}

type PortScan struct {
	Start int    `yaml:"start"`
	End   int    `yaml:"end"`
	Type  string `yaml:"type"`
}

type DirScan struct {
	DirTxt string `yaml:"dirTxt"`
}

func NewConfig(confPath string) *Config {
	return &Config{
		ConfPath: confPath,
	}
}

func (c *Config) LoadConfig() {
	yamlFile, err := ioutil.ReadFile(c.ConfPath)
	if err != nil {
		log.Printf("Failed to read YAML file: %v", err)
		return
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("Failed to unmarshal YAML: %v", err)
		return
	}

	log.Printf("Config Loaded: %+v", c.Pipeline.PortScan)
}
