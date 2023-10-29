{{- reserveImport "database/sql/driver" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{ range .Models }}
func ({{ .GoName }}) CreateTableStmt() string {
	return {{ quote (createTable .) }}
}
func ({{ .GoName }}) AlterTableStmt() string {
	return {{ quote (alterTable .) }}
}
{{ if .HasTableName -}}
func ({{ .GoName }}) TableName() string {
	return {{ quote .TableName }}
}
{{- end }}
{{ if .HasColumn -}}
func ({{ .GoName }}) Columns() []string {
	return {{ "[]string{" }}{{- range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ quote $f.ColumnName }}{{ end }}{{- "}" }}
}
{{- end }}
{{ if ne .PK nil -}}
func (v {{ .GoName }}) IsAutoIncr() bool {
	return {{ .PK.IsAutoIncr }}
}
func (v {{ .GoName }}) PK() (columnName string, pos int, value driver.Value) {
	return {{ quote .PK.Field.ColumnName }}, {{ .PK.Field.Index }}, {{ castAs "v" .PK.Field }}
}
{{- end }}
func (v {{ .GoName }}) Values() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ castAs "v" $f }}{{ end }}{{- `}` }}
}
func (v *{{ .GoName }}) Addrs() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ `, ` }}{{ end }}{{ addrOf "v" $f }}{{ end }}{{- `}` }}
}
{{ end }}