# Getting Started

> A comprehensive guide to use sqlgen.

1.  Install sqlgen.

    ```console
    go install github.com/si3nloong/sqlgen
    ```

2.  Define your struct.

    <h5 a><strong><code>model/user.go</code></strong></h5>

    ```go
    package model

    import "time"

    type LongText string

    type User struct {
        ID int64
        Name LongText
        Age uint8
        Created time.Time
    }
    ```

3.  Generate the output files.

    ```console
    sqlgen generate ./model/user.go
    ```

4.  This will generate the output file.

    <h5 a><strong><code>model/generated.go</code></strong></h5>

    ```go
    // Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

    package model

    import (
        "github.com/si3nloong/sqlgen/sql/types"
    )

    // Implements `sql.Valuer` interface.
    func (User) Table() string {
        return "user"
    }

    func (User) Columns() []string {
        return []string{"id", "name", "age", "created"}
    }

    func (v User) Values() []any {
        return []any{v.ID, string(v.Name), int64(v.Age), v.Created}
    }

    func (v *User) Addrs() []any {
        return []any{&v.ID, types.String(&v.Name), types.Integer(&v.Age), &v.Created}
    }

    ```

    The code generator will help you generate the necessary boilerplate codes for you to utilise with `github.com/si3nloong/sqlgen/sql` package later on.