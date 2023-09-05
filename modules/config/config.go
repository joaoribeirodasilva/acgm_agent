package config

import (
	"fmt"
	"os"
	"path/filepath"

	"biqx.com.br/acgm_agent/modules/cmd"
	"gopkg.in/yaml.v1"
)

const DEFAULT_DIR = "/etc/acgm"
const DEV_DIR = "./etc/acgm"
const DEFAULT_FILE = "agent.yaml"

type Config struct {
	Version  float64        `json:"version" yaml:"version"`
	Database ConfigDatabase `json:"database" yaml:"database"`
	Logs     ConfigLogs     `json:"logs" yaml:"logs"`
	Metrics  ConfigMetrics  `json:"metrics" yaml:"metrics"`
	Settings ConfigSettings `json:"settings" yaml:"settings"`
	Nginx    Nginx          `json:"nginx" yaml:"nginx"`
	path     string         `json:"-" yaml:"-"`
	//dev      bool           `json:"-" yaml:"-"`
}

func (c *Config) Read(options *cmd.CmdOptions) error {

	var err error
	path := ""

	if options.ConfigFile != nil && *options.ConfigFile != "" {
		path, err = filepath.Abs(*options.ConfigFile)
		if err != nil {
			panic(fmt.Sprintf("ERROR: configuration file path error! - %s", err.Error()))
		}
		if !exists(path) {
			panic(fmt.Sprintf("ERROR: configuration file %s does not exist!", path))
		}
		c.path = path
	} else {
		var pathEtc string
		pathEtc, err = filepath.Abs(fmt.Sprintf("%s/%s", DEFAULT_DIR, DEFAULT_FILE))
		if err != nil {
			panic(fmt.Sprintf("ERROR: configuration file path error! - %s", err.Error()))
		}
		if !exists(pathEtc) {
			var pathLocal string
			pathLocal, err = filepath.Abs(fmt.Sprintf("%s/%s", DEV_DIR, DEFAULT_FILE))
			if err != nil {
				panic(fmt.Sprintf("ERROR: configuration file path error! - %s", err.Error()))
			}

			if !exists(pathLocal) {
				panic(fmt.Sprintf("ERROR: no configuration file found at %s and %s!", pathEtc, pathLocal))
			}
			c.path = pathLocal
		} else {
			c.path = pathEtc
		}
	}

	var fileBytes []byte

	fileBytes, err = os.ReadFile(c.path)
	if err != nil {
		panic(fmt.Sprintf("ERROR: reading configuration file %s! - %s", c.path, err.Error()))
	}

	err = yaml.Unmarshal(fileBytes, c)
	if err != nil {
		panic(fmt.Sprintf("ERROR: parsing configuration file %s! - %s", c.path, err.Error()))
	}

	return nil
}

func (c *Config) Write() error {

	fileBytes, err := yaml.Marshal(c)
	if err != nil {
		panic(fmt.Sprintf("ERROR: building configuration yaml! - %s", err.Error()))
	}

	err = os.WriteFile(c.path, fileBytes, 0622)
	if err != nil {
		panic(fmt.Sprintf("ERROR: saving configuration file %s! - %s", c.path, err.Error()))
	}

	return nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
