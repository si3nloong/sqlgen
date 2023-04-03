sqlgen not only generate your model functions, but it also provide bunch of helper functions to ease your work on manipulate SQL query.

- CRUD
  - [INSERT single record.](#insert-single-or-multiple-records)
  - [INSERT multiple records.](#insert-single-or-multiple-records)
  - [UPDATE single record using primery key.](#update-single-record)
  - UPDATE multiple records.
  - [DELETE single record using primary key.](#delete-single-record)
  - [TRUNCATE table.](#delete-single-record)

```go
import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ...

db, err := sql.Open("mysql", "user:password@/dbname")
if err != nil {
	panic(err)
}
```

### SELECT single record

```go
import (
    "context"
    "github.com/si3nloong/sqlgen/sql"
)

// ...

result, err := sql.SelectOneFrom[Model](context.TODO(), db)
if err != nil {
    panic(err)
}
```

### SELECT multiple records

```go
import (
    "context"
    "github.com/si3nloong/sqlgen/sql"
)

// ...

models, err := sql.SelectFrom[Model](context.TODO(), db)
if err != nil {
    panic(err)
}
```

### INSERT single or multiple records

```go
import (
    "context"
    "github.com/si3nloong/sqlgen/sql"
)

// ...

models := []*Model{&{}, &{}}
result, err := sql.InsertInto(context.TODO(), db, models)
if err != nil {
    panic(err)
}
```

### UPDATE single record

`UpdateOne` using interface `KeyValuer` as the primary identify on which record to perform update.

```go
import (
    "context"
    "github.com/si3nloong/sqlgen/sql"
)

// ...

model := Model{}
model.Name = "new name"

result, err := sql.UpdateOne(context.TODO(), db, model)
if err != nil {
    panic(err)
}
```

### DELETE single record

`DeleteOne` using interface `KeyValuer` as the primary identify on which record to perform update.

```go
import (
    "context"
    "github.com/si3nloong/sqlgen/sql"
)

// ...

model := Model{}

result, err := sql.DeleteOne(context.TODO(), db, model)
if err != nil {
    panic(err)
}
```
