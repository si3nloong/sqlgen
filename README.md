# sqlgen

[![Build](https://github.com/si3nloong/sqlgen/workflows/test/badge.svg?branch=main)](https://github.com/si3nloong/sqlgen/actions?query=workflow%3Atest)
[![Go Report](https://goreportcard.com/badge/github.com/si3nloong/sqlgen)](https://goreportcard.com/report/github.com/si3nloong/sqlgen)
[![Go Coverage](https://codecov.io/gh/si3nloong/sqlgen/branch/main/graph/badge.svg)](https://codecov.io/gh/si3nloong/sqlgen)
[![LICENSE](https://img.shields.io/github/license/si3nloong/sqlgen)](https://github.com/si3nloong/sqlgen/blob/main/LICENSE)

> sqlgen is not an ORM, it's a compiler. It make mapping to go struct without any extra costs (reflection) incur.

## What is sqlgen?

- **sqlgen is based on a Schema first approach** — You get to Define your Model using the go struct.
- **sqlgen enables Codegen** — We generate the boring bits, so you can focus on building your app quickly.
- **sqlgen prioritizes performance** — We use generics to eliminate runtime reflection costs.
- **sqlgen prioritizes Type safety**
- **sqlgen eliminates Side Effects** - You will get expected results instead of side effects when mutate your models.

## Quick start

1.  Install sqlgen.

    ```console
    go install github.com/si3nloong/sqlgen
    ```

2.  Define your struct.

    ```go
    import "time"

    type User struct {
        ID      int64
        Name    string
        Age     uint8
        Created time.Time
    }
    ```

3.  Generate the output files.

    ```console
    sqlgen generate <source_file>
    ```

More help to get started:

- [Getting started tutorial](/docs/README.md) - a comprehensive guide to help you get started
- [Reference docs](/docs/API.md) for the APIs

## Reporting Issues

If you think you've found a bug, or something isn't behaving the way you think it should, please raise an [issue](https://github.com/si3nloong/sqlgen/issues) on GitHub.

## Contributing

We welcome contributions, Read our [Contribution Guidelines](https://github.com/si3nloong/sqlgen/blob/main/CONTRIBUTING.md) to learn more about contributing to **sqlgen**

## Big Thanks To

Thanks to these awesome companies for their support of Open Source developers ❤

[![GitHub](https://jstools.dev/img/badges/github.svg)](https://github.com/open-source)

## License

[MIT](https://github.com/si3nloong/sqlgen/blob/main/LICENSE)

Copyright (c) 2023-present, SianLoong Lee
