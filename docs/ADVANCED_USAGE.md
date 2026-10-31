# Advanced Usages

### Multiple Insert

```go
package main

import (
    "context"
    "database/sql"
    "time"

    "cloud.google.com/go/civil"
    _ "github.com/go-sql-driver/mysql"
    "db"
    "model"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    dbConn, err := sql.Open("mysql", "root:abcd1234@/sqlbench?parseTime=true")
    if err != nil {
        panic(err)
    }
    defer dbConn.Close()

    birthDate, _ := civil.ParseDate("1995-01-28")

    if _, err := db.InsertInto(ctx, dbConn, []model.User{
        {Name: "John Doe", Gender: model.Male, BirthDate: birthDate, Created: time.Now()},
        {Name: "YY", Gender: model.Female, BirthDate: birthDate, Created: time.Now()},
        {Name: "Yoman", Gender: model.Male, BirthDate: birthDate, Created: time.Now()},
    }); err != nil {
        panic(err)
    }
}
```

### Multiple Upsert

```go
package main

import (
    "context"
    "database/sql"
    "time"

    "cloud.google.com/go/civil"
    _ "github.com/go-sql-driver/mysql"
    "db"
    "model"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    dbConn, err := sql.Open("mysql", "root:abcd1234@/sqlbench?parseTime=true")
    if err != nil {
        panic(err)
    }
    defer dbConn.Close()

    birthDate, _ := civil.ParseDate("1995-01-28")

    if _, err := db.Upsert(ctx, dbConn, []model.User{
        {Name: "John Doe", Gender: model.Male, BirthDate: birthDate, Created: time.Now()},
        {Name: "YY", Gender: model.Female, BirthDate: birthDate, Created: time.Now()},
        {Name: "Yoman", Gender: model.Male, BirthDate: birthDate, Created: time.Now()},
    }); err != nil {
        panic(err)
    }
}
```

### QueryOneStmt

```go
    /*
    SELECT `id`, `name`, `birth_date`, `gender`, `address`, `created`
    FROM `user` WHERE `gender` = 0 AND `birth_date` >= "1995-01-28"
    ORDER BY `created` DESC LIMIT 50;
    */
    users, err := db.QueryOneStmt[model.User](context.Background(), dbConn, func(v User) {
        return db.SelectOneStmt{
            Select:    v.Columns(),
            FromTable: v.TableName(),
            Where: db.And(
                db.Equal(v.GetGender(), model.Female),
                db.GreaterThanOrEqual(v.GetBirthDate(), birthDate),
            ),
            OrderBy: []sequel.OrderByClause{
                db.Desc(v.GetCreated()),
            },
        }
    })
    if err != nil {
        panic(err)
    }
```

### QueryStmt

```go
    /*
    SELECT `id`, `name`, `birth_date`, `gender`, `address`, `created`
    FROM `user` WHERE `gender` = 0 AND `birth_date` >= "1995-01-28"
    ORDER BY `created` DESC LIMIT 50;
    */
    users, err := db.QueryStmt[model.User](context.Background(), dbConn, func(v User) {
        return db.SelectStmt{
            Select:    v.Columns(),
            FromTable: v.TableName(),
            Where: db.And(
                db.Equal(v.GetGender(), model.Female),
                db.GreaterThanOrEqual(v.GetBirthDate(), birthDate),
            ),
            OrderBy: []sequel.OrderByClause{
                db.Desc(v.GetCreated()),
            },
            Limit: 50,
        }
    })
    if err != nil {
        panic(err)
    }
```

### ExecStmt with UPDATE

```go
    /*
    UPDATE `user` SET `name` = "John" WHERE `id` = 1;
    */
    if _, err := db.ExecStmt(context.Background(), dbConn, func(v User) db.UpdateStmt {
		return db.UpdateStmt{
			Where: db.Equal(v.ColumnID(), 1),
    		Set: []sequel.SetClause{
				db.Set(v.ColumnName(), "John"),
			},
		}
	}); err != nil {
        panic(err)
    }
```

### ExecStmt with DELETE

```go
    /*
    DELETE FROM `user` WHERE `id` = 1;
    */
	if _, err := db.ExecStmt(context.Background(), dbConn, func(v User) db.DeleteStmt {
		return db.DeleteStmt{
			Where: db.Equal(v.ColumnID(), 1),
		}
	}); err != nil {
        panic(err)
    }
```

### Pagination

```go
    /*
    SELECT `id`, `name`, `birth_date`, `gender`, `address`, `created`
    FROM `user` WHERE `gender` = 0 AND `birth_date` >= "1995-01-28"
    ORDER BY `created` DESC LIMIT 50;
    */
    p := db.Paginate[model.User](func(v User) {
        return db.SelectStmt{
            Select:    v.Columns(),
            FromTable: v.TableName(),
            Where: db.And(
                db.Equal(v.GetGender(), model.Female),
                db.GreaterThanOrEqual(v.GetBirthDate(), birthDate),
            ),
            OrderBy: []sequel.OrderByClause{
                db.Desc(v.GetCreated()),
            },
            Limit: 50,
        }
    })
    // Loop through user table by using cursor-based pagination
    for result, err := range p.Next(context.Background(), dbConn) {
        if err != nil {
            panic(err)
        }
        for _, v := range result {
            log.Println(v)
        }
    }
```

### JOIN queries

Both generated `SelectStmt` and `SelectOneStmt` support ordered joins. Use
`sequel.JoinClause` callbacks, convenience helpers, or chainable statement methods.
Each callback writes SQL using the same statement writer as `WhereClause`. Regenerate your database package to expose the
new statement fields and methods.

```go
on := db.OnColumn(
    db.ColumnRef("users.id"),
    db.ColumnRef("o.user_id"),
)

// Composable callbacks returned by convenience helpers.
stmt := db.SelectStmt{
    Select:    []string{"users.id", "users.name"},
    FromTable: "users",
    Joins: []sequel.JoinClause{
        db.LeftJoin("orders AS o", on),
    },
}

// Equivalent chainable form; methods return a copy.
stmt = (db.SelectStmt{
    Select:    []string{"users.id", "users.name"},
    FromTable: "users",
}).LeftJoin("orders AS o", on)
```

Available helpers and statement methods are `InnerJoin`, `LeftJoin`, `RightJoin`,
`FullJoin`, and `CrossJoin`. `CrossJoin` takes only a table.
The statement's `Join(joins...)` method appends callbacks in order.
`sequel.JoinClause` is shared; JOIN helpers belong to each generated driver package.
MySQL does not generate the `FullJoin` helper or statement
methods. PostgreSQL and SQLite generate them; joins still require support from
your database version.

The `on` argument accepts existing `WhereClause` callbacks and generated `And`/`Or` helpers.
Use `OnColumn` to compare two columns, and existing value predicates (or
`w.Var(value)` in a callback) to bind values. JOIN and WHERE values share the same
placeholder sequence. Keep optional-side filters in `On` when unmatched rows
should remain in an outer join.

Table names, aliases, and `ColumnRef` are trusted SQL text, just like `FromTable`
and `Select`; supply any dialect-specific identifier quoting yourself. Never use
user input as identifier text. Generated model columns also work with
`OnColumn`, but qualify columns when names occur in multiple tables.

Pass the resulting statement through `QueryStmt` or `QueryOneStmt`. Explicitly
select columns in the order expected by the destination's `Addrs()` method.
To return fields from multiple tables, use a result type implementing that
scanner contract. Outer-joined fields must accept SQL NULL (for example through
`sql.NullString`). JOINs do not automatically populate nested models or collapse
multiple result rows into one parent.

Custom callbacks can be mixed with the built-in helpers. Include a leading space
and use `w.Var` for bound values:

```go
custom := sequel.JoinClause(func(w sequel.StmtWriter) {
    w.WriteString(" LEFT JOIN orders AS o ON o.user_id = users.id AND o.status = ")
    w.WriteString(w.Var("paid"))
})
stmt = stmt.Join(custom)
```

Callbacks render directly; the database validates SQL during execution.
