copy:
	rm -rf ./cmd/sqlgen/codegen/sequel.go.tpl
	cp -rf ./sequel/sequel.go ./cmd/sqlgen/codegen/sequel.go.tpl
build:
	go build -o sqlgen -ldflags="-s -w" ./main.go