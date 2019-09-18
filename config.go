package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

//Config represents user defined config
type Config struct {
	Addons     []string          `yaml:"addons"`
	CurseForge []string          `yaml:"curseforge"`
	GitHub     map[string]string `yaml:"github"`
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fdata, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(fdata, &cfg)
	if err != nil {
		return nil, err
	}

	if len(cfg.Addons) > 0 {
		log.Fatal(`Config key "addons" was renamed to "curseforge", please update your config accordingly`)
	}

	return cfg, nil
}
