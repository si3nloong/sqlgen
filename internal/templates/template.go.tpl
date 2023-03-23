package {{ .GoPkg }}

func ({{ .GoName }}) Table() string {
	return {{ quote .Name }}
}

func ({{ .GoName }}) Columns() []string {
    return {{ "[]string{" }}{{- range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ quote $f.Name }}{{ end }}{{- "}" }}
}

func (v {{ .GoName }}) Values() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ cast (print `v.` $f.Name) $f }}{{ end }}{{- `}` }}
}

func (v *{{ .GoName }}) Addrs() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ `, ` }}{{ end }}{{ addr (print `&v.` $f.Name) $f }}{{ end }}{{- `}` }}
}