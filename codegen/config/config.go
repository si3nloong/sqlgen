package config

import (
	"os"
	"path/filepath"

	"github.com/si3nloong/sqlgen/internal/fileutil"
	"gopkg.in/yaml.v3"
)

type sqlDriver string

const (
	MySQL    sqlDriver = "mysql"
	Postgres sqlDriver = "postgres"
	Sqlite   sqlDriver = "sqlite"
)

type naming string

const (
	SnakeCase  naming = "snake_case"
	CamelCase  naming = "camelCase"
	PascalCase naming = "PascalCase"
)

var cfgFilenames = []string{".sqlgen.yml", ".sqlgen.yaml", "sqlgen.yml", "sqlgen.yaml"}

type Config struct {
	Source           []string  `yaml:"src"`
	Driver           sqlDriver `yaml:"driver"`
	NamingConvention naming    `yaml:"namingConvention,omitempty"`
	Tag              string    `yaml:"tag,omitempty"`
	Strict           bool      `yaml:"strict"`
	IncludeHeader    bool      `yaml:"includeHeader"`
	SourceMap        bool      `yaml:"sourceMap"`
}

func (c *Config) init() {
	c.Source = []string{"./**/*"}
	c.NamingConvention = SnakeCase
	c.Tag = "sql"
	c.Driver = MySQL
	c.Strict = true
	c.IncludeHeader = true
}

func LoadConfigFrom(src string) (*Config, error) {
	cfg := new(Config)
	cfg.init()
	if err := decodeToFile(src, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func DefaultConfig() *Config {
	cfg := new(Config)
	cfg.init()
	file, exists := findCfgInDir(fileutil.Getpwd())
	if !exists {
		return cfg
	}
	if err := decodeToFile(file, cfg); err != nil {
		panic(err)
	}
	return cfg
}

func decodeToFile(src string, cfg *Config) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(cfg); err != nil {
		return err
	}
	return nil
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
