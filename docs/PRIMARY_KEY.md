# Primary Key

## Auto Increment Key

```go
type Location struct {
	ID   int64 `sql:",pk,auto_increment"`
	Name string
}
```

It will generate the following codes :

```go
...
func (Location) IsAutoIncr() {}
func (v *Location) ScanAutoIncr(val int64) error {
	v.ID = int64(val)
	return nil
}
func (v Location) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (Location) Columns() []string {
	return []string{"id", "name"} // 2
}
func (v Location) Values() []any {
	return []any{
		v.Name, // 1 - name
	}
}
func (v *Location) Addrs() []any {
	return []any{
		&v.ID,   // 0 - id
		&v.Name, // 1 - name
	}
}
func (Location) InsertColumns() []string {
	return []string{"name"} // 1
}
func (Location) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v Location) InsertOneStmt() (string, []any) {
	return "INSERT INTO `location` (`name`) VALUES (?);", []any{v.Name}
}
func (v Location) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`name` FROM `location` WHERE `id` = ? LIMIT 1;", []any{v.ID}
}
func (v Location) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `location` SET `name` = ? WHERE `id` = ?;", []any{v.Name, v.ID}
}
...
```

## Composite Key

```go
type Location struct {
    CountryCode string `sql:",pk"`
    PartyID     string `sql:",pk"`
}
```

It will generate the following codes :

```go
...
func (v Location) CompositeKey() ([]string, []int, []any) {
	return []string{"country_code", "party_id"}, []int{0, 1}, []any{v.CountryCode, v.PartyID}
}
func (v Location) InsertOneStmt() (string, []any) {
	return "INSERT INTO `location` (`country_code`,`party_id`,`name`) VALUES (?,?,?);", v.Values()
}
func (v Location) FindOneByPKStmt() (string, []any) {
	return "SELECT `country_code`,`party_id`,`name` FROM `location` WHERE (`country_code`,`party_id`) = (?,?) LIMIT 1;", []any{v.CountryCode, v.PartyID}
}
func (v Location) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `location` SET `name` = ? WHERE (`country_code`,`party_id`) = (?,?);", []any{v.Name, v.CountryCode, v.PartyID}
}
...
```

## Primary Key

```go
type Location struct {
    Location string `sql:",pk"`
}
```

It will generate the following codes :

```go
...
func (v Location) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (Location) Columns() []string {
	return []string{"id", "name"} // 2
}
func (v Location) InsertOneStmt() (string, []any) {
	return "INSERT INTO `location` (`id`,`name`) VALUES (?,?);", v.Values()
}
func (v Location) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`name` FROM `location` WHERE `id` = ? LIMIT 1;", []any{v.ID}
}
func (v Location) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `location` SET `name` = ? WHERE `id` = ?;", []any{v.Name, v.ID}
}
...
```
