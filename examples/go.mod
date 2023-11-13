module github.com/si3nloong/sqlgen/examples

go 1.21

toolchain go1.21.3

require (
	cloud.google.com/go v0.110.10
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gofrs/uuid/v5 v5.0.0
	github.com/google/uuid v1.4.0
	github.com/jaswdr/faker v1.16.0
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/si3nloong/sqlgen v0.0.0-20230404062952-cbb69b02fc6a
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.38.1 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/tools v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/si3nloong/sqlgen => ../
