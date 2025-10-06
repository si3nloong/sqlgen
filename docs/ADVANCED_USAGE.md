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

    if _, err := db.UpsertInto(ctx, dbConn, []model.User{
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
    users, err := db.QueryOneStmt[model.User](ctx, dbConn, func(v User) {
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
    users, err := db.QueryStmt[model.User](ctx, dbConn, func(v User) {
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

### ExecStmt

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
