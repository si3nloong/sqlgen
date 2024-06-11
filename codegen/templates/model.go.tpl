{{- reserveImport "strings" }}
{{- reserveImport "database/sql/driver" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- $omitGetters := .OmitGetters -}}
{{ range .Models }}
{{- $hasCustomTabler := .HasTableName -}}
{{- $hasCustomDatabaser := .HasDatabaseName -}}
{{- $hasNotOnlyPK := .HasNotOnlyPK -}}
{{- $structName := .GoName -}}
func (v {{ $structName }}) CreateTableStmt() string {
	return {{ createTable "v" . }}
}
{{ if and (not $hasCustomDatabaser) (ne .DatabaseName nil) -}}
func ({{ $structName }}) DatabaseName() string {
	return {{ quote (quoteIdentifier .DatabaseName) }}
}
{{ end -}}
{{ if not $hasCustomTabler -}}
func ({{ $structName }}) TableName() string {
	return {{ quote (quoteIdentifier .TableName) }}
}
func ({{ $structName }}) InsertOneStmt() string {
	return {{ insertOneStmt . }}
}
{{ end -}}
{{- /* postgres will not generate this because the argument number alway different */ -}}
{{ if isStaticVar -}}
func ({{ $structName }}) InsertVarQuery() string {
	return {{ quote (varStmt .) }}
}
{{ end -}}
{{ if not .HasColumn -}}
func ({{ $structName }}) Columns() []string {
	return {{ "[]string{" }}{{- range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ quote (quoteIdentifier $f.ColumnName) }}{{ end }}{{- "}" }}
}
{{ end -}}
{{ if ne .Keys nil -}}
func ({{ $structName }}) HasPK() {}
{{ end -}}
{{ if ne .Keys nil -}}
{{- /* If it has auto increment key */ -}}
{{ if .IsAutoIncr -}}
func ({{ $structName }}) IsAutoIncr() {}
{{ end -}}
func (v {{ $structName }}) PK() ([]string, []int, []any) {
	return {{ `[]string{` }}{{ range $i, $f := .Keys }}{{- if $i }}{{ ", " }}{{ end }}{{ quote (quoteIdentifier $f.ColumnName) }}{{ end }}{{- `},` }}{{ `[]int{` }}{{ range $i, $f := .Keys }}{{- if $i }}{{ ", " }}{{ end }}{{ $f.Index }}{{ end }}{{- `},` }}{{ `[]any{` }}{{ range $i, $f := .Keys }}{{- if $i }}{{ ", " }}{{ end }}{{ castAs $f }}{{ end }}{{- `}` }}
}
{{ if (and (not $hasCustomTabler) ($hasNotOnlyPK)) -}}
{{- /* If it has static table and columns other than key */ -}}
func ({{ $structName }}) FindByPKStmt() string {
	return {{ findByPKStmt . }}
}
{{ end -}}
{{ if (and (not $hasCustomTabler) ($hasNotOnlyPK)) -}}
{{- /* If it has static table and columns other than key */ -}}
func ({{ $structName }}) UpdateByPKStmt() string {
	return {{ updateByPKStmt . }}
}
{{ end -}}
{{ end -}}
func (v {{ $structName }}) Values() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ ", " }}{{ end }}{{ castAs $f }}{{ end }}{{- `}` }}
}
func (v *{{ $structName }}) Addrs() []any {
	return {{ `[]any{` }}{{ range $i, $f := .Fields }}{{- if $i }}{{ `, ` }}{{ end }}{{ addrOf "v" $f }}{{ end }}{{- `}` }}
}
{{ if not $omitGetters -}}
{{ range $f := .Fields -}}
{{- $return := getFieldTypeValue $f -}}
func (v {{ $structName }}) {{ $return.FuncName }}() (sequel.ColumnValuer[{{ $return.Type }}]) {
	return sequel.Column{{ typeConstraint $return }}({{ quote (quoteIdentifier $f.ColumnName) }}, v.{{ .GoPath }}, func(val {{ $return.Type }}) driver.Value { return {{ castAs $f "val" }} })
}
{{ end -}}
{{ end -}}
{{ "" }}
{{ end }}