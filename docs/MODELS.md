# Declaring Models

## Fields Tags

Format

```
`sql:"[column name],[tag options...]"`
```

Tags are case insensitive, however `snake_case` is preferred. Tags are optional to use when declaring models, `sqlgen` supports the following tags:

| Tag Name         | Alias | Description                                                                                                      | Example               |
| ---------------- | ----- | ---------------------------------------------------------------------------------------------------------------- | --------------------- |
| `primary_key`    | `pk`  | Specifies column as primary key                                                                                  | sql:",pk"             |
| `auto_increment` | -     | Specifies column auto incrementable                                                                              | sql:",auto_increment" |
| `binary`         | -     | Specifies column value using serializer `encoding.BinaryMarshaler` and deserializer `encoding.BinaryUnmarshaler` | sql:",binary"         |
| `size`           | -     | Specifies column data size/length                                                                                | sql:",size:10"        |

<!-- | `unique`         | -     | Specifies column as unique                                                                                       | sql:",unique"                |
| `unsigned`       | -     | Specifies column as unsigned, only applicable for `constraints.Float`                                            | sql:",unsigned"              |
| `convert`        | -     | Specifies serializer for how to serialize and deserialize data into db, e.g:                                     | sql:",convert:"              |
| `type`           | -     | Column data type                                                                                                 | sql:",type:CHAR(3) NOT NULL" | -->
