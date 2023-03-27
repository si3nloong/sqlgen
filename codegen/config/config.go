package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var cfgFilenames = []string{".sqlgen.yml", "sqlgen.yml", "sqlgen.yaml"}

type Config struct {
	SrcDir           string `yaml:"srcDir"`
	NamingConvention string `yaml:"namingConvention,omitempty"`
	Tag              string `yaml:"tag,omitempty"`
}

func pwd() string {
	path, _ := os.Getwd()
	return path
}

func (c *Config) init() {
	c.SrcDir = pwd()
	c.NamingConvention = "snakecase"
	c.Tag = "sql"
}

func DefaultConfig() *Config {
	cfg := new(Config)
	cfg.init()
	filesrc, exists := findCfgInDir(pwd())
	if !exists {
		return cfg
	}
	f, _ := os.Open(filesrc)
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(cfg); err != nil {
		panic(err)
	}
	return cfg
}

func findCfgInDir(dir string) (string, bool) {
	for _, cfgName := range cfgFilenames {
		path := filepath.Join(dir, cfgName)
		if isFileExists(path) {
			return path, true
		}
	}
	return "", false
}

func isFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
