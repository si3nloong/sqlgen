# Configuration

`sqlgen` can be configured using a `sqlgen.yml`, `sqlgen.yml`, `.sqlgen.yaml` or `sqlgen.yaml` file, by default it will be loaded from the current directory.

Example :

```yml
# Where are all the model files located? globs are supported eg  src/**/*.go
src:
  - ./**/*

# Optional:
# Possibly values : `mysql`, `postgres` or `sqlite`
driver: mysql

# Optional:
# Possibly values : `snake_case`, `camelCase` or `PascalCase`
naming_convention: snake_case

# Optional:
struct_tag: sql

# Optional:
skip_escape: false

# Optional: Where should the generated model code go?
exec:
  skip_empty: false
  filename: generated.go

# Optional: Where should the generated database code go?
database:
  dir: db
  package: db
  filename: db.go
  operator:
    dir: db
    filename: operator.go

# Optional:
strict: true

# Optional: turn on to not generate any file header in generated files
skip_header: false

# Optional:
source_map: false
# Optional: set to skip running `go mod tidy` when generating server code
# skip_mod_tidy: false
```
