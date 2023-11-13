{{- reserveImport "strings" }}
{{- reserveImport "database/sql/driver" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{ range .Models }}
{{- $hasCustomTabler := .HasTableName -}}
{{- $structName := .GoName -}}
func (v {{ $structName }}) CreateTableStmt() string {
	return {{ createTable "v" . }}
}
func ({{ $structName }}) AlterTableStmt() string {
	return {{ quote (alterTable .) }}
}
{{ if eq $hasCustomTabler false -}}
func ({{ $structName }}) TableName() string {
	return {{ quote (wrap .TableName) }}
}
func (v {{ $structName }}) InsertOneStmt() string {
	return {{ insertOneStmt . }}
}
{{ end -}}
func ({{ $structName }}) InsertVarQuery() string {
	return {{ quote (varStmt .) }}
}
{{ if eq .HasColumn false -}}
func ({{ $structName }}) Columns() []string {
	return {{ "[]string{" }}{{- range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ quote (wrap $f.ColumnName) }}{{ end }}{{- "}" }}
}
{{ end -}}
{{ if ne .PK nil -}}
{{ if .PK.IsAutoIncr -}}
func ({{ $structName }}) IsAutoIncr() {}
{{ end -}}
func (v {{ $structName }}) PK() (columnName string, pos int, value driver.Value) {
	return {{ quote (wrap .PK.Field.ColumnName) }}, {{ .PK.Field.Index }}, {{ castAs .PK.Field }}
}
{{ if eq $hasCustomTabler false -}}
func (v {{ $structName }}) FindByPKStmt() string {
	return {{ findByPKStmt . }}
}
{{ end -}}
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
	return sequel.Column[{{ $return.Type }}]({{ quote (wrap $f.ColumnName) }}, v.{{ .GoPath }}, func(vi {{ $return.Type }}) driver.Value { return {{ castAs $f "vi" }} })
}
{{ end -}}
{{ "" }}
{{ end }}