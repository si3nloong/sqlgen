# Where are all the model files located? globs are supported eg  src/**/*.go
src:
{{- range .Source }}
  - {{ . }}
{{- end }}

# Optional: 
driver: {{ .Driver }}

# Optional: 
naming_convention: {{ .NamingConvention }}

# Optional: 
struct_tag: {{ .Tag }}

# Optional: Where should any generated code go?
exec:
  filename: {{ .Exec.Filename }}

# Optional: 
database:
  dir: {{ .Database.Dir }}
  package: {{ .Database.Package }}
  filename: {{ .Database.Filename }}

# Optional: 
strict: {{ .Strict }}

# Optional: 
skip_header: {{ .SkipHeader }}

# Optional: 
source_map: {{ .SourceMap }}

# Optional: set to skip running `go mod tidy` when generating server code
# skip_mod_tidy: {{ .SkipModTidy }}