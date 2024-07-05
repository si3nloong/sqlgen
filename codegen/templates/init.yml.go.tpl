# Where are all the model files located? globs are supported eg  src/**/*.go
src:
{{- range .Source }}
  - {{ . }}
{{- end }}

# Optional: possibly values : `mysql`, `postgres` or `sqlite`
driver: {{ .Driver }}

# Optional: possibly values : `snake_case`, `camelCase` or `PascalCase`
naming_convention: {{ .NamingConvention }}

# Optional: the struct tag for "sqlgen" to read from. Default is `sql`
struct_tag: {{ .Tag }}

# Optional: whether to omit the quote on table name and column name
omit_quote_identifier: {{ .OmitQuoteIdentifier }}

# Optional: whether to omit the getters
omit_getters: {{ .OmitGetters }}

# Optional: to add prefix to getter
getter:
  prefix: {{ .Getter.Prefix }}

# Optional: where should the generated model code go?
exec:
  skip_empty: {{ .Exec.SkipEmpty }}
  filename: {{ .Exec.Filename }}

# Optional: where should the generated database code go?
database:
  dir: {{ .Database.Dir }}
  package: {{ .Database.Package }}
  filename: {{ .Database.Filename }}
  operator:
    filename: {{ .Database.Operator.Filename }}

# Optional: 
no_strict: {{ .NoStrict }}

# Optional: turn on to not generate any file header in generated files
skip_header: {{ .SkipHeader }}

# Optional: 
source_map: {{ .SourceMap }}

# Optional: set to skip running `go mod tidy` when generating server code
skip_mod_tidy: {{ .SkipModTidy }}

models:
