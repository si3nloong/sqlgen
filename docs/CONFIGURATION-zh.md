# 配置

`sqlgen` 可以通过 `sqlgen.yml`，`.sqlgen.yml`，`.sqlgen.yaml` 或 `sqlgen.yaml` 文件进行配置，默认情况下会从当前目录加载。

示例：

```yml
# 所有模型文件的位置在哪里？支持通配符，例如 src/**/*.go
src:
  - ./**/*

# 可选：
# 可能的值：`mysql`，`postgres` 或 `sqlite`
driver: mysql

# 可选：
# 可能的值：`snake_case`，`camelCase` 或 `PascalCase`
naming_convention: snake_case

# 可选：
struct_tag: sql

# 可选：
skip_escape: false

# 可选：生成的模型代码应放在哪里？
exec:
  skip_empty: false
  filename: generated.go

# 可选：生成的数据库代码应放在哪里？
database:
  dir: db
  package: db
  filename: db.go
  operator:
    dir: db
    filename: operator.go

# 可选：
strict: true

# 可选：设置为 true 以在生成的文件中不生成任何文件头
skip_header: false

# 可选：
source_map: false
# 可选：设置为跳过在生成服务器代码时运行 `go mod tidy`
# skip_mod_tidy: false
```
