# sqlgen

[![Build](https://github.com/si3nloong/sqlgen/workflows/test/badge.svg?branch=main)](https://github.com/si3nloong/sqlgen/actions?query=workflow%3Atest)
[![Go Report](https://goreportcard.com/badge/github.com/si3nloong/sqlgen)](https://goreportcard.com/report/github.com/si3nloong/sqlgen)
[![Go Coverage](https://codecov.io/gh/si3nloong/sqlgen/branch/main/graph/badge.svg)](https://codecov.io/gh/si3nloong/sqlgen)
[![LICENSE](https://img.shields.io/github/license/si3nloong/sqlgen)](https://github.com/si3nloong/sqlgen/blob/main/LICENSE)

> sqlgen is not an ORM, it's a code generator instead. It parse the go struct and generate the necessary methods on struct for you.

## What is sqlgen?

- **sqlgen is based on a Code first approach** — You don't require to write SQL first, but Go code instead.
- **sqlgen enables Codegen** — We generate the boring bits, so you can focus on building your app quickly.
- **sqlgen prioritizes performance** — Most of the things will define in compile time instead of runtime.
- **sqlgen embrace generics** — We use generics to eliminate runtime reflection costs and reduce memory allocation.
- **sqlgen eliminates Side Effects** - You will get expected results instead of side effects when mutate your models.

## Quick start

1.  Install sqlgen.

    ```console
    go install github.com/si3nloong/sqlgen
    ```

2.  Define your struct.

    `models/user.go`

    ```go
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
    sqlgen generate models/user.go
    ```

More help to get started:

- [Getting started tutorial](/docs/GET_STARTED.md) - a comprehensive guide to help you get started
- [CLI guide](/docs/CLI.md) for the CLI command.

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
