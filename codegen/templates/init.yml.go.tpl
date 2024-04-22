# Where are all the model files located? globs are supported eg  src/**/*.go
src:
{{- range .Source }}
  - {{ . }}
{{- end }}

# Optional: 
# Possibly values : `mysql`, `postgres` or `sqlite`
driver: {{ .Driver }}

# Optional: 
# Possibly values : `snake_case`, `camelCase` or `PascalCase`
naming_convention: {{ .NamingConvention }}

# Optional:
# The struct tag for "sqlgen" to read from. Default is `sql`
struct_tag: {{ .Tag }}

# Optional:
# Whether to omit the quote on table name and column name
omit_quote_identifier: {{ .QuoteIdentifier }}

# Optional:
# Whether to omit the getters
omit_getter: {{ .OmitGetter }}

# Optional: to add prefix to getter
getter:
  prefix: "Prefix"

# Optional: Where should the generated model code go?
exec:
  skip_empty: {{ .Exec.SkipEmpty }}
  filename: {{ .Exec.Filename }}

# Optional: Where should the generated database code go?
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