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
struct_tag: {{ .Tag }}

# Optional: 
skip_escape: {{ .SkipEscape }}

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
strict: {{ .Strict }}

# Optional: turn on to not generate any file header in generated files
skip_header: {{ .SkipHeader }}

# Optional: 
source_map: {{ .SourceMap }}

# Optional: set to skip running `go mod tidy` when generating server code
# skip_mod_tidy: {{ .SkipModTidy }}