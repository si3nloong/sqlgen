copy:
	rm -rf ./codegen/sequel.go.tpl
	cp -rf ./sequel/model.go ./codegen/sequel.go.tpl
build:
	go build -o sqlgen -ldflags="-s -w" ./main.go