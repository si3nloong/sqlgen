# Declaring Models

## Fields Tags

Format

```
`sql:"[column name],[tag options...]"`
```

Example:

```go
type User struct{
    ID   int64  `sql:",pk,auto_increment"`
    Name string `sql:",size:10"`
}
```

Tags are case insensitive, however `snake_case` is preferred. Tags are optional to use when declaring models, `sqlgen` supports the following tags:

| Tag Name         | Alias | Description                                                                                                      | Example                                                                |
| ---------------- | ----- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `primary_key`    | `pk`  | Specifies column as primary key                                                                                  | `sql:",pk"`                                                            |
| `auto_increment` | -     | Specifies column auto incrementable                                                                              | `sql:",auto_increment"`                                                |
| `binary`         | -     | Specifies column value using serializer `encoding.BinaryMarshaler` and deserializer `encoding.BinaryUnmarshaler` | `sql:",binary"`                                                        |
| `size`           | -     | Specifies column data size/length                                                                                | `sql:",size:10"`                                                       |
| `encode`         | -     | Specifies custom `sql.Valuer`                                                                                    | `sql:",encode:github.com/si3nloong/sqlgen/encoding.MarshalStringList"` |
| `decode`         | -     | Specifies custom `sql.Scanner`                                                                                   | `sql:",decode:github.com/si3nloong/sqlgen/types.Bool"`                 |

<!-- | `unique`         | -     | Specifies column as unique                                                                                       | sql:",unique"                |
| `unsigned`       | -     | Specifies column as unsigned, only applicable for `constraints.Float`                                            | sql:",unsigned"              |
| `type`           | -     | Column data type                                                                                                 | sql:",type:CHAR(3) NOT NULL" | -->
