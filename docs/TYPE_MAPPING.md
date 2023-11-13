# Data Types

A scalar message field can have one of the following types – the table shows the type specified in the .proto file, and the corresponding type in the automatically generated class:

|                                  | mysql              | postgres (not stable) | sqlite (not stable) |
| -------------------------------- | ------------------ | --------------------- | ------------------- |
| `~string`                        | VARCHAR(255)       | VARCHAR(255)          | VARCHAR(255)        |
| `~[]byte`                        | BLOB               | BLOB                  | BLOB                |
| `~[16]byte`                      | BINARY(16)         | BINARY(16)            | BINARY(16)          |
| `sql.RawBytes`                   | VARCHAR(255)       | VARCHAR(255)          | VARCHAR(255)        |
| `~int`                           | INTEGER            | INTEGER               | INTEGER             |
| `~int8`                          | TINYINT            | TINYINT               | TINYINT             |
| `~int16`                         | SMALLINT           | SMALLINT              | SMALLINT            |
| `~int32`                         | MEDIUMINT          | MEDIUMINT             | MEDIUMINT           |
| `~int64`                         | BIGINT             | BIGINT                | BIGINT              |
| `~uint`                          | INTEGER UNSIGNED   | INTEGER UNSIGNED      | INTEGER UNSIGNED    |
| `~uint8`                         | TINYINT UNSIGNED   | TINYINT UNSIGNED      | TINYINT UNSIGNED    |
| `~uint16`                        | SMALLINT UNSIGNED  | SMALLINT UNSIGNED     | SMALLINT UNSIGNED   |
| `~uint32`                        | MEDIUMINT UNSIGNED | MEDIUMINT UNSIGNED    | MEDIUMINT UNSIGNED  |
| `~uint64`                        | BIGINT UNSIGNED    | BIGINT UNSIGNED       | BIGINT UNSIGNED     |
| `~float32`                       | FLOAT              | FLOAT                 | FLOAT               |
| `~float64`                       | FLOAT              | FLOAT                 | FLOAT               |
| `cloud.google.com/go/civil.Date` | DATE               | DATE                  | DATE                |
| `time.Time`                      | DATETIME           | DATETIME              | DATETIME            |
| `struct`                         | JSON               | JSON                  | JSON                |
| `array`                          | JSON               | JSON                  | JSON                |
| `slice`                          | JSON               | JSON                  | JSON                |
