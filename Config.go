package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ConfigFile string `yaml:"-"`
	Domain     string `yaml:"domain,omitempty"`
	Port       int    `yaml:"port,omitempty"`
	Https      bool   `yaml:"https,omitempty"`
	Log        struct {
		Datetime bool `yaml:"datetime,omitempty"`
		SrcFile  bool `yaml:"srcfile,omitempty"`
		Info     bool `yaml:"info,omitempty"`
		Warning  bool `yaml:"warning,omitempty"`
		Error    bool `yaml:"error,omitempty"`
		Trace    bool `yaml:"trace,omitempty"`
		Debug    bool `yaml:"debug,omitempty"`
	} `yaml:"log,omitempty"`
}

func initFlag() {
	flag.StringVar(&conf.ConfigFile, "config", "config.yaml", "config file")
	flag.IntVar(&conf.Port, "port", 8090, "port")
	flag.BoolVar(&conf.Https, "https", true, "https mode")
	flag.BoolVar(&conf.Log.Datetime, "log-datetime", false, "log datetime enable")
	flag.BoolVar(&conf.Log.SrcFile, "log-srcfile", true, "log source file enable")
	flag.BoolVar(&conf.Log.Info, "log-info", true, "log info enable")
	flag.BoolVar(&conf.Log.Warning, "log-warning", true, "log warning enable")
	flag.BoolVar(&conf.Log.Error, "log-error", true, "log error enable")
	flag.BoolVar(&conf.Log.Trace, "log-trace", true, "log trace enable")
	flag.BoolVar(&conf.Log.Debug, "log-debug", false, "log debug enable")
	flag.Usage = usage
}

func (c *Config) Read(fn string) error {
	buf, err := os.ReadFile(fn)
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
		fmt.Println(err.Error())
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
		fmt.Printf("cannot read config file: %v", err)
	}

	err = conf.checkRequired()
	if err != nil {
		fmt.Printf("config check failed: %v", err)
	}
}
