{{- reserveImport "strings" }}
{{- reserveImport "database/sql/driver" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{ range .Models }}
{{- $structName := .GoName -}}
func (v {{ $structName }}) CreateTableStmt() string {
	return {{ createTable "v" . }}
}
func ({{ $structName }}) AlterTableStmt() string {
	return {{ quote (alterTable .) }}
}
{{ if eq .HasTableName false -}}
func ({{ $structName }}) TableName() string {
	return {{ quote (wrap .TableName) }}
}
{{ end -}}
func (v {{ $structName }}) InsertOneStmt() string {
	return {{ insertOneStmt . }}
}
func ({{ $structName }}) InsertVarQuery() string {
	return {{ quote (varStmt .) }}
}
{{ if eq .HasColumn false -}}
func ({{ $structName }}) Columns() []string {
	return {{ "[]string{" }}{{- range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ quote (wrap $f.ColumnName) }}{{ end }}{{- "}" }}
}
{{ end -}}
{{ if ne .PK nil -}}
func (v {{ $structName }}) IsAutoIncr() bool {
	return {{ .PK.IsAutoIncr }}
}
func (v {{ $structName }}) PK() (columnName string, pos int, value driver.Value) {
	return {{ quote (wrap .PK.Field.ColumnName) }}, {{ .PK.Field.Index }}, {{ castAs .PK.Field }}
}
func (v {{ $structName }}) FindByPKStmt() string {
	return {{ findByPKStmt . }}
}
{{ end -}}
func (v {{ $structName }}) Values() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ castAs $f }}{{ end }}{{- `}` }}
}
func (v *{{ $structName }}) Addrs() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ `, ` }}{{ end }}{{ addrOf "v" $f }}{{ end }}{{- `}` }}
}
{{ range $f := .Fields -}}
{{- $return := getFieldTypeValue $f -}}
func (v {{ $structName }}) {{ $return.FuncName }}() (sequel.ColumnValuer[{{ $return.Type }}]) {
	return sequel.Column[{{ $return.Type }}]({{ quote (wrap $f.ColumnName) }}, v.{{ .GoName }}, func(vi {{ $return.Type }}) driver.Value { return {{ castAs $f "vi" }} })
}
{{ end -}}
{{ end }}