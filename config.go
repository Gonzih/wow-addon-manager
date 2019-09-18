package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//Config represents user defined config
type Config struct {
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

	return cfg, nil
}
