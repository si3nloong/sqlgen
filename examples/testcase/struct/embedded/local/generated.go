package embedded

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "name", "z", "created", "ok"} // 5
}
func (v B) Values() []any {
	return []any{
		v.a.ID,       // 0 - id
		v.a.Name,     // 1 - name
		v.a.Z,        // 2 - z
		v.ts.Created, // 3 - created
		v.ts.OK,      // 4 - ok
	}
}
func (v *B) Addrs() []any {
	return []any{
		&v.a.ID,       // 0 - id
		&v.a.Name,     // 1 - name
		&v.a.Z,        // 2 - z
		&v.ts.Created, // 3 - created
		&v.ts.OK,      // 4 - ok
	}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?)" // 5
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO `b` (`id`,`name`,`z`,`created`,`ok`) VALUES (?,?,?,?,?);", v.Values()
}
func (v B) IDValue() any {
	return v.a.ID
}
func (v B) NameValue() any {
	return v.a.Name
}
func (v B) ZValue() any {
	return v.a.Z
}
func (v B) CreatedValue() any {
	return v.ts.Created
}
func (v B) OKValue() any {
	return v.ts.OK
}
func (v B) ColumnID() sequel.ColumnClause {
	return sequel.BasicColumn("id", v.a.ID)
}
func (v B) ColumnName() sequel.ColumnClause {
	return sequel.BasicColumn("name", v.a.Name)
}
func (v B) ColumnZ() sequel.ColumnClause {
	return sequel.BasicColumn("z", v.a.Z)
}
func (v B) ColumnCreated() sequel.ColumnClause {
	return sequel.BasicColumn("created", v.ts.Created)
}
func (v B) ColumnOK() sequel.ColumnClause {
	return sequel.BasicColumn("ok", v.ts.OK)
}
