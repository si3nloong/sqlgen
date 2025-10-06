# Data Types

A scalar message field can have one of the following types – the table shows the type specified in the .proto file, and the corresponding type in the automatically generated class:

|                            | **mysql**          | **postgres**     | **sqlite** (not stable) |
| -------------------------- | ------------------ | ---------------- | ----------------------- |
| `~string`                  | VARCHAR(255)       | varchar(255)     | TEXT                    |
| `~bool`                    | BOOLEAN            | bool             | BOOLEAN                 |
| `~byte`                    | CHAR(1)            | char(1)          | TEXT                    |
| `~rune`                    | CHAR(1)            | char(1)          | TEXT                    |
| `~[]byte`                  | BLOB               | blob             | BLOB                    |
| `~[16]byte`                | BINARY(16)         | BINARY(16)       | BINARY(16)              |
| `~int`                     | INTEGER            | int4             | INTEGER                 |
| `~int8`                    | TINYINT            | int2             | INTEGER                 |
| `~int16`                   | SMALLINT           | int2             | INTEGER                 |
| `~int32`                   | MEDIUMINT          | int4             | INTEGER                 |
| `~int64`                   | BIGINT             | int8             | INTEGER                 |
| `~uint`                    | INTEGER UNSIGNED   | int4             | INTEGER                 |
| `~uint8`                   | TINYINT UNSIGNED   | int2             | INTEGER                 |
| `~uint16`                  | SMALLINT UNSIGNED  | int2             | INTEGER                 |
| `~uint32`                  | MEDIUMINT UNSIGNED | int4             | INTEGER                 |
| `~uint64`                  | BIGINT UNSIGNED    | int8             | INTEGER                 |
| `~float32`                 | FLOAT              | real             | FLOAT                   |
| `~float64`                 | FLOAT              | double precision | FLOAT                   |
| `~[...]rune`               | CHAR(**:size**)    | char(**:size**)  | TEXT                    |
| `~[...]byte`               | CHAR(**:size**)    | char(**:size**)  | TEXT                    |
| `~[]rune`                  | VARCHAR(255)       | text             | TEXT                    |
| `~[]string`                | JSON               | text[]           | TEXT                    |
| `~[]bool`                  | JSON               | bool[]           | TEXT                    |
| `~[]int`                   | JSON               | int4[]           | TEXT                    |
| `~[]int8`                  | JSON               | int2[]           | TEXT                    |
| `~[]int16`                 | JSON               | int2[]           | TEXT                    |
| `~[]int32`                 | JSON               | int4[]           | TEXT                    |
| `~[]int64`                 | JSON               | int8[]           | TEXT                    |
| `~[]uint`                  | JSON               | int4[]           | TEXT                    |
| `~[]uint8`                 | JSON               | int2[]           | TEXT                    |
| `~[]uint16`                | JSON               | int2[]           | TEXT                    |
| `~[]uint32`                | JSON               | int4[]           | TEXT                    |
| `~[]uint64`                | JSON               | int8[]           | TEXT                    |
| `~[]float32`               | JSON               | double[]         | TEXT                    |
| `~[]float64`               | JSON               | double[]         | TEXT                    |
| `struct`                   | JSON               | json             | TEXT                    |
| `array`                    | JSON               | json             | TEXT                    |
| `slice`                    | JSON               | json             | TEXT                    |
| `time.Time`                | DATETIME           | timestamptz(6)   | TEXT                    |
| `database/sql.RawBytes`    | VARCHAR(255)       | VARCHAR(255)     | TEXT                    |
| `encoding/json.RawMessage` | JSON               | json             | TEXT                    |
| `any`                      | JSON               | json             | TEXT                    |
