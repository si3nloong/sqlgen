module github.com/si3nloong/sqlgen/cmd/sqlgen

go 1.27

require (
	github.com/AlecAivazis/survey/v2 v2.3.7
	github.com/go-playground/validator/v10 v10.30.5
	github.com/go-sql-driver/mysql v1.10.1
	github.com/goccy/go-yaml v1.19.2
	github.com/jackc/pgx/v5 v5.11.0
	github.com/mattn/go-sqlite3 v1.14.52
	github.com/samber/lo v1.53.0
	github.com/si3nloong/sqlgen v1.0.0-beta.0.0.20260419173412-8b71ef1df632
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.12.1
	golang.org/x/exp v0.0.0-20261014143641-f0251eb4aa3e
	golang.org/x/tools v0.52.0
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/go-playground/locales v0.14.2 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.3 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/term v0.47.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

replace github.com/si3nloong/sqlgen => ../..
