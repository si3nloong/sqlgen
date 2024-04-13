# sqlgen

[![Build](https://github.com/si3nloong/sqlgen/workflows/test/badge.svg?branch=main)](https://github.com/si3nloong/sqlgen/actions?query=workflow%3Atest)
[![Go Report](https://goreportcard.com/badge/github.com/si3nloong/sqlgen)](https://goreportcard.com/report/github.com/si3nloong/sqlgen)
[![Go Coverage](https://codecov.io/gh/si3nloong/sqlgen/branch/main/graph/badge.svg)](https://codecov.io/gh/si3nloong/sqlgen)
[![LICENSE](https://img.shields.io/github/license/si3nloong/sqlgen)](https://github.com/si3nloong/sqlgen/blob/main/LICENSE)

> sqlgen is not an ORM, it's a code generator instead. It parse the go struct and generate the necessary struct methods for you.

## What is sqlgen?

- **sqlgen is based on a Code first approach** — You don't require to write SQL first, but Go code instead.
- **sqlgen enables Codegen** — We generate the boring bits, so you can focus on building your app quickly.
- **sqlgen prioritizes Performance** — Most of the things will define in compile time instead of runtime.
- **sqlgen embrace Generics** — We use generics to eliminate runtime reflection costs and reduce memory allocation.
- **sqlgen eliminates Side Effects** - You will get expected results instead of side effects when mutate your models.

## SQL driver support

| Driver     | Support |
| ---------- | :-----: |
| `mysql`    |   ✅    |
| `postgres` |   ✅    |
| `sqlite`   |   ✅    |

## Quick start

1.  Install sqlgen.

    ```console
    go install github.com/si3nloong/sqlgen/cmd/sqlgen@main
    ```

2.  Define your [model](./docs/MODELS.md) struct.

    <h5 a><strong><code>model/user.go</code></strong></h5>

    ```go
    package model

    import "time"

    type LongText string

    type User struct {
        ID      int64 `sql:",auto_increment"`
        Name    string
        Age     uint8
        Address LongText
        Created time.Time
    }
    ```

3.  Generate the output files.

    ```bash
    # sqlgen generate <source_file>
    sqlgen generate model/user.go
    ```

4.  Generated code will be as follow.

    <h5 a><strong><code>model/generated.go</code></strong></h5>

    ```go
    // Code generated by sqlgen, version v1.0.0-alpha.4. DO NOT EDIT.

    package model

    import (
        "database/sql/driver"
        "time"

        "github.com/si3nloong/sqlgen/sequel/types"
    )

    func (User) CreateTableStmt() string {
        return "CREATE TABLE IF NOT EXISTS `user` (`id` BIGINT NOT NULL AUTO_INCREMENT,`name` VARCHAR(255) NOT NULL,`age` TINYINT UNSIGNED NOT NULL,`address` VARCHAR(255) NOT NULL,`created` DATETIME NOT NULL,PRIMARY KEY (`id`));"
    }
    func (User) AlterTableStmt() string {
        return "ALTER TABLE `user` MODIFY `id` BIGINT NOT NULL AUTO_INCREMENT,MODIFY `name` VARCHAR(255) NOT NULL AFTER `id`,MODIFY `age` TINYINT UNSIGNED NOT NULL AFTER `name`,MODIFY `address` VARCHAR(255) NOT NULL AFTER `age`,MODIFY `created` DATETIME NOT NULL AFTER `address`;"
    }
    func (User) TableName() string {
        return "`user`"
    }
    func (User) Columns() []string {
        return []string{"`id`", "`name`", "`age`", "`address`", "`created`"}
    }
    func (v User) IsAutoIncr() {}
    func (v User) PK() (columnName string, pos int, value driver.Value) {
        return "`id`", 0, int64(v.ID)
    }
    func (v User) Values() []any {
        return []any{int64(v.ID), string(v.Name), int64(v.Age), string(v.Address), time.Time(v.Created)}
    }
    func (v *User) Addrs() []any {
        return []any{types.Integer(&v.ID), types.String(&v.Name), types.Integer(&v.Age), types.String(&v.Address), (*time.Time)(&v.Created)}
    }
    ```

More help to get started:

- [Getting started tutorial](/docs/GET_STARTED.md) - a comprehensive guide to help you get started
- [CLI guide](/docs/CLI.md) - the CLI commands.
- [FAQ](/docs/FAQ.md) - frequent ask questions.
- [Configuration file](/docs/CONFIGURATION.md) - configure code generation.

## Benchmark

<img src="./docs/images/orm_benchmark.jpg" />

## Reporting Issues

If you think you've found a bug, or something isn't behaving the way you think it should, please raise an [issue](https://github.com/si3nloong/sqlgen/issues) on GitHub.

## Contributing

We welcome contributions, Read our [Contribution Guidelines](https://github.com/si3nloong/sqlgen/blob/main/CONTRIBUTING.md) to learn more about contributing to **sqlgen**

## Big Thanks To

Thanks to these awesome companies for their support of Open Source developers ❤

[![GitHub](https://jstools.dev/img/badges/github.svg)](https://github.com/open-source)

## Inspired By

- [You don't need orm in Go](https://medium.com/@enverbisevac/you-dont-need-orm-in-go-9216fb74cdfd)
- [gqlgen](https://github.com/99designs/gqlgen)

## License

[MIT](https://github.com/si3nloong/sqlgen/blob/main/LICENSE)

Copyright (c) 2023-present, SianLoong Lee
