package config

import (
	"os"
	"path/filepath"

	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/internal/strfmt"
	"github.com/si3nloong/sqlgen/sequel"
	"gopkg.in/yaml.v3"

	_ "github.com/si3nloong/sqlgen/sequel/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/postgres"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/sqlite"
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

const (
	DefaultConfigFile    = "sqlgen.yml"
	DefaultGeneratedFile = "generated.go"
	DefaultStructTag     = "sql"
)

var cfgFilenames = []string{DefaultConfigFile, ".sqlgen.yml", ".sqlgen.yaml", "sqlgen.yaml"}

type Config struct {
	Source           []string       `yaml:"src"`
	Driver           sqlDriver      `yaml:"driver"`
	NamingConvention naming         `yaml:"naming_convention,omitempty"`
	Tag              string         `yaml:"struct_tag,omitempty"`
	Strict           bool           `yaml:"strict"`
	SkipEscape       bool           `yaml:"skip_escape"`
	Exec             ExecConfig     `yaml:"exec"`
	Database         DatabaseConfig `yaml:"database"`
	SkipHeader       bool           `yaml:"skip_header"`
	SourceMap        bool           `yaml:"source_map"`
	SkipModTidy      bool           `yaml:"skip_mod_tidy"`
}

type ExecConfig struct {
	Filename string `yaml:"filename"`
}

type DatabaseConfig struct {
	Package  string `yaml:"package"`
	Dir      string `yaml:"dir"`
	Filename string `yaml:"filename"`
}

func (c *Config) init() {
	c.Source = []string{"./**/*"}
	c.NamingConvention = SnakeCase
	c.Tag = DefaultStructTag
	c.Driver = MySQL
	c.Strict = true
	c.Exec.Filename = DefaultGeneratedFile
	c.Database.Package = "db"
	c.Database.Dir = "db"
	c.Database.Filename = "db.go"
}

func (c Config) Clone() *Config {
	newConfig := &Config{}
	newConfig.init()
	if len(c.Source) > 0 {
		newConfig.Source = make([]string, len(c.Source))
		copy(newConfig.Source, c.Source)
	}
	if c.Driver != "" {
		newConfig.Driver = c.Driver
	}
	if c.NamingConvention != "" {
		newConfig.NamingConvention = c.NamingConvention
	}
	if c.Tag != "" {
		newConfig.Tag = c.Tag
	}
	if c.Exec.Filename != "" {
		newConfig.Exec.Filename = c.Exec.Filename
	}
	if c.Database.Dir != "" {
		newConfig.Database.Dir = c.Database.Dir
	}
	if c.Database.Package != "" {
		newConfig.Database.Package = c.Database.Package
	}
	if c.Database.Filename != "" {
		newConfig.Database.Filename = c.Database.Filename
	}
	if c.Strict != newConfig.Strict {
		newConfig.Strict = c.Strict
	}
	if c.SkipHeader != newConfig.SkipHeader {
		newConfig.SkipHeader = c.SkipHeader
	}
	if c.SourceMap != newConfig.SourceMap {
		newConfig.SourceMap = c.SourceMap
	}
	if c.SkipModTidy != newConfig.SkipModTidy {
		newConfig.SkipModTidy = c.SkipModTidy
	}
	return newConfig
}

func (c Config) Dialect() sequel.Dialect {
	d, ok := sequel.GetDialect(string(c.Driver))
	if !ok {
		panic("sqlgen: missing dialect, please register your dialect first")
	}
	return d
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
