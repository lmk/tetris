package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ConfigFile string `yaml:"-"`
	Domain     string `yaml:"domain,omitempty"`
	Port       int    `yaml:"port,omitempty"`
	IsDebug    bool   `yaml:"debug,omitempty"`
}

func (c *Config) Read(fn string) error {
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("cannot read config %s: %v", fn, err)
	}

	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return fmt.Errorf("cannot unmarshal config %s: %v", fn, err)
	}

	return nil
}

func (c *Config) makePretty() string {

	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		Info.Printf(err.Error())
	}

	return string(buf)
}

func (c *Config) checkRequired() error {

	if c.Domain == "" {
		c.Domain = "localhost"
	}

	if c.Port == 0 {
		c.Port = 8090
	}

	return nil
}

func initConf() {
	err := conf.Read(conf.ConfigFile)
	if err != nil {
		Error.Printf("cannot read config file: %v", err)
	}

	err = conf.checkRequired()
	if err != nil {
		Error.Printf("config check failed: %v", err)
	}

	Info.Printf("config: %s", conf.makePretty())
}
