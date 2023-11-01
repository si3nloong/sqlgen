package config

import (
	"os"
	"path/filepath"

	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/internal/strfmt"
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

const ConfigFile = "sqlgen.yml"

var cfgFilenames = []string{ConfigFile, ".sqlgen.yml", ".sqlgen.yaml", "sqlgen.yaml"}

type Config struct {
	Source           []string  `yaml:"src"`
	Driver           sqlDriver `yaml:"driver"`
	NamingConvention naming    `yaml:"naming_convention,omitempty"`
	Tag              string    `yaml:"tag,omitempty"`
	Strict           bool      `yaml:"strict"`
	Exec             struct {
		Filename string `yaml:"filename"`
	} `yaml:"exec"`
	Database struct {
		Package  string `yaml:"package"`
		Dir      string `yaml:"dir"`
		Filename string `yaml:"filename"`
	} `yaml:"database"`
	SkipHeader  bool `yaml:"skip_header"`
	SourceMap   bool `yaml:"source_map"`
	SkipModTidy bool `yaml:"skip_mod_tidy"`
}

func (c *Config) init() {
	c.Source = []string{"./**/*"}
	c.NamingConvention = SnakeCase
	c.Tag = "sql"
	c.Driver = MySQL
	c.Strict = true
	c.Exec.Filename = "generated.go"
	c.Database.Package = "db"
	c.Database.Dir = "db"
	c.Database.Filename = "db.go"
}

func (c Config) RenameFunc() func(string) string {
	switch c.NamingConvention {
	case SnakeCase:
		return strfmt.ToSnakeCase
	case CamelCase:
		return strfmt.ToCamelCase
	case PascalCase:
		return strfmt.ToPascalCase
	default:
		return func(s string) string { return s }
	}
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
