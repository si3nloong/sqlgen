package config

import (
	"maps"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/si3nloong/sqlgen/internal/strfmt"

	_ "github.com/si3nloong/sqlgen/sequel/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/postgres"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/sqlite"
)

type SqlDriver string

const (
	MySQL    SqlDriver = "mysql"
	Postgres SqlDriver = "postgres"
	Sqlite   SqlDriver = "sqlite"
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
	Source           []string  `yaml:"src"`
	Driver           SqlDriver `yaml:"driver"`
	NamingConvention naming    `yaml:"naming_convention,omitempty"`
	Tag              string    `yaml:"struct_tag,omitempty"`
	// Whether to quote the table name and column name
	OmitQuoteIdentifier bool              `yaml:"omit_quote_identifier,omitempty"`
	OmitGetters         bool              `yaml:"omit_getters,omitempty"`
	Strict              *bool             `yaml:"strict,omitempty"`
	Exec                ExecConfig        `yaml:"exec"`
	Getter              GetterConfig      `yaml:"getter"`
	Migration           *MigrationConfig  `yaml:"migration"`
	Database            *DatabaseConfig   `yaml:"database"`
	SourceMap           bool              `yaml:"source_map"`
	SkipHeader          bool              `yaml:"skip_header"`
	SkipModTidy         bool              `yaml:"skip_mod_tidy"`
	Models              map[string]*Model `yaml:"models"`
}

type Model struct {
	DataType   string `yaml:"data_type"`
	Scanner    string `yaml:"scan" validate:"required_with=Valuer"`
	SQLScanner string `yaml:"sql_scan" validate:"required_with=SQLValuer"`
	Valuer     string `yaml:"value" validate:"required_with=Scanner"`
	SQLValuer  string `yaml:"sql_value" validate:"required_with=SQLScanner"`
}

type ExecConfig struct {
	// Skip `generated.go` file being generated if the input has no matching struct
	SkipEmpty bool `yaml:"skip_empty"`
	// Declare the name of generated go file. Default is `generated.go`
	Filename string `yaml:"filename"`
}

type GetterConfig struct {
	Prefix string `yaml:"prefix"`
}

type MigrationConfig struct {
}

type DatabaseConfig struct {
	Package  string                  `yaml:"package"`
	Dir      string                  `yaml:"dir"`
	Filename string                  `yaml:"filename"`
	Operator *DatabaseOperatorConfig `yaml:"operator"`
}

type DatabaseOperatorConfig struct {
	Package  string `yaml:"package"`
	Dir      string `yaml:"dir"`
	Filename string `yaml:"filename"`
}

func (c *Config) init() {
	strict := true
	c.Source = []string{"./**/*"}
	c.NamingConvention = SnakeCase
	c.Tag = DefaultStructTag
	c.Driver = MySQL
	c.Strict = &strict
	c.Exec.Filename = DefaultGeneratedFile
	c.Getter.Prefix = "Get"
	c.Database = new(DatabaseConfig)
	c.Database.Package = "db"
	c.Database.Dir = "db"
	c.Database.Filename = "db.go"
	c.Database.Operator = new(DatabaseOperatorConfig)
	c.Database.Operator.Package = c.Database.Package
	c.Database.Operator.Dir = c.Database.Dir
	c.Database.Operator.Filename = "operator.go"
	c.Models = make(map[string]*Model)
}

func (c *Config) initIfEmpty() {
	if c.Source == nil {
		c.Source = []string{"./**/*"}
	}
	if c.Driver == "" {
		c.Driver = MySQL
	}
	if c.NamingConvention == "" {
		c.NamingConvention = SnakeCase
	}
	if c.Tag == "" {
		c.Tag = DefaultStructTag
	}
	if c.Getter.Prefix == "" {
		c.Getter.Prefix = "Get"
	}
	if c.Exec.Filename == "" {
		c.Exec.Filename = DefaultGeneratedFile
	}

	if c.Database == nil {
		c.Database = new(DatabaseConfig)
	}
	if c.Database.Operator == nil {
		c.Database.Operator = new(DatabaseOperatorConfig)
	}
	if c.Database.Package == "" {
		c.Database.Package = "db"
	}
	if c.Database.Dir == "" {
		c.Database.Dir = "db"
	}
	if c.Database.Filename == "" {
		c.Database.Filename = "db.go"
	}
	if c.Database.Operator.Package == "" {
		c.Database.Operator.Package = c.Database.Package
	}
	if c.Database.Operator.Dir == "" {
		c.Database.Operator.Dir = c.Database.Dir
	}
	if c.Database.Operator.Filename == "" {
		c.Database.Operator.Filename = "operator.go"
	}
	if c.Models == nil {
		c.Models = make(map[string]*Model)
	}
}

func (c *Config) Merge(mapCfg *Config) *Config {
	if mapCfg.Source != nil {
		c.Source = append([]string{}, mapCfg.Source...)
	}
	if mapCfg.Driver != "" {
		c.Driver = mapCfg.Driver
	}
	if mapCfg.NamingConvention != "" {
		c.NamingConvention = mapCfg.NamingConvention
	}
	if mapCfg.Strict != nil {
		strict := *mapCfg.Strict
		c.Strict = &strict
	}
	if mapCfg.Tag != "" {
		c.Tag = mapCfg.Tag
	}
	if mapCfg.Getter.Prefix != "" {
		c.Getter.Prefix = mapCfg.Getter.Prefix
	}
	if mapCfg.Exec.SkipEmpty {
		c.Exec.SkipEmpty = true
	}
	if mapCfg.Exec.Filename != "" {
		c.Exec.Filename = mapCfg.Exec.Filename
	}
	if mapCfg.Database != nil {
		c.Database = &DatabaseConfig{
			Dir:      mapCfg.Database.Dir,
			Package:  mapCfg.Database.Package,
			Filename: mapCfg.Database.Filename,
		}

		if mapCfg.Database.Operator != nil {
			c.Database.Operator = &DatabaseOperatorConfig{
				Dir:      mapCfg.Database.Operator.Dir,
				Package:  mapCfg.Database.Operator.Package,
				Filename: mapCfg.Database.Operator.Filename,
			}
		}
	}
	if mapCfg.SkipHeader {
		c.SkipHeader = true
	}
	if mapCfg.SkipModTidy {
		c.SkipModTidy = true
	}
	if mapCfg.Models != nil {
		maps.Copy(c.Models, mapCfg.Models)
	}
	c.initIfEmpty()
	return c
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
