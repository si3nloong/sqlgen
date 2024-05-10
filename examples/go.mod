module github.com/si3nloong/sqlgen/examples

go 1.21

toolchain go1.21.3

require (
	cloud.google.com/go v0.113.0
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gofrs/uuid/v5 v5.0.0
	github.com/google/uuid v1.6.0
	github.com/jaswdr/faker v1.16.0
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/si3nloong/sqlgen v1.0.0-alpha.3.0.20231118095154-390f9683bb93
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.39.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/si3nloong/sqlgen => ../
