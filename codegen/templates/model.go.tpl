{{- reserveImport "database/sql/driver" }}

{{ range .Models -}}
const CreateTable{{ .GoName }} = `
-- Migration script for table {{ quote .Name }}
CREATE TABLE `+ {{ quote (wrap .Name) }} +`IF NOT EXISTS (
	{{ range $i, $f := .Fields }}
	{{- $f.Name }} VARCHAR(20){{- if $i }}{{ ", " }}{{- end }}
	{{- end }}
)
`

func ({{ .GoName }}) Table() string {
	return {{ quote .Name }}
}

func ({{ .GoName }}) Columns() []string {
	return {{ "[]string{" }}{{- range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ quote $f.Name }}{{ end }}{{- "}" }}
}

{{ if ne .PK nil -}}
func ({{ .GoName }}) IsAutoIncr() bool {
	{{- if hasTag .PK "auto" "auto_increment" }}
	return true
	{{- else }}
	return false
	{{- end }}
}

func (v {{ .GoName }}) PK() (int, any) {
	{{- if isValuer .PK }}
    return {{ .PK.Index }}, ((driver.Valuer)(v.{{ .PK.GoName }}))
	{{- else }}
	return {{ .PK.Index }}, {{ cast "v" .PK }}
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