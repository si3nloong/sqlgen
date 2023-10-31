# Where are all the model files located? globs are supported eg  src/**/*.go
src:
  - ./examples/**/*

# Optional: 
driver: mysql

# Optional: 
namingConvention: snake_case

# Optional: 
tag: sql

# Optional: Where should any generated code go?
exec:
  filename: generated.go

# Optional: 
database:
  dir: db
  package: db
  filename: "{name}.resolvers.go"

# Optional: 
strict: true

# Optional: 
includeHeader: true

# Optional: 
sourceMap: false
