copy:
	cp -rf ./sequel/sequel.go ./codegen/sequel.go.tpl
build:
	go build -o sqlgen -ldflags="-s -w" ./main.go