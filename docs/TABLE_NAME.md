# Custom Table Name

```go
package model

import (
    "github.com/si3nloong/sqlgen/sequel"
)

type User struct {
	sequel.Table `sql:"CustomUserTableName"`
	Name string
}
```

To rename your table name, you may embed the struct `sequel.Table` and provide the value in struct tag.
