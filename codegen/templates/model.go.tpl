{{- reserveImport "database/sql/driver" }}

{{ range .Models -}}
// Implements `sql.Valuer` interface.
func ({{ .GoName }}) Table() string {
	return {{ quote .Name }}
}

func ({{ .GoName }}) Columns() []string {
	return {{ "[]string{" }}{{- range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ quote $f.Name }}{{ end }}{{- "}" }}
}

{{ if ne .PK nil -}}
func (v {{ .GoName }}) Key() (driver.Value, error) {
	{{ if isValuer .PK -}}
    return ((driver.Valuer)(v.{{ .PK.GoName }})).Value()
	{{ else -}}
	return v.{{ .PK.GoName }}, nil
	{{- end }}
}
{{ end }}
func (v {{ .GoName }}) Values() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ cast "v" $f }}{{ end }}{{- `}` }}
}

func (v *{{ .GoName }}) Addrs() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ `, ` }}{{ end }}{{ addr "v" $f }}{{ end }}{{- `}` }}
}

{{ end }}