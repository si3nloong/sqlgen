package config

import (
	"os"
	"path/filepath"

	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/internal/gosyntax"
	"gopkg.in/yaml.v3"
)

var cfgFilenames = []string{".sqlgen.yml", ".sqlgen.yaml", "sqlgen.yml", "sqlgen.yaml"}

type Config struct {
	SrcDir           string `yaml:"srcDir"`
	Driver           string `yaml:"driver" survey:"driver"`
	NamingConvention string `yaml:"namingConvention,omitempty" survey:"namingConvention"`
	Tag              string `yaml:"tag,omitempty" survey:"tag"`
	IncludeHeader    *bool  `yaml:"includeHeader"`
}

func (c *Config) init() {
	c.SrcDir = "."
	c.NamingConvention = "snake_case"
	c.Tag = "sql"
	c.Driver = "mysql"
	c.IncludeHeader = gosyntax.PtrOf(true)
}

func DefaultConfig() *Config {
	cfg := new(Config)
	cfg.init()
	file, exists := findCfgInDir(fileutil.Getpwd())
	if !exists {
		return cfg
	}
	f, _ := os.Open(file)
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(cfg); err != nil {
		panic(err)
	}
	return cfg
}

func findCfgInDir(dir string) (string, bool) {
	for _, cfgName := range cfgFilenames {
		path := filepath.Join(dir, cfgName)
		if fileutil.IsFileExists(path) {
			return path, true
		}
	}
	return "", false
}
