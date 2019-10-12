.PHONY: all proto build_server build_client

all: build_server build_client

proto:
	@mkdir -p web/proto
	@protoc -I proto --plugin=protoc-gen-ts=./web/node_modules/.bin/protoc-gen-ts --js_out=import_style=commonjs,binary:web/proto --ts_out=service=grpc-web:web/proto --go_out=plugins=grpc:go/pkg/api proto/temperature.proto

run_web: proto
	@npm run --prefix ./web start:dev

build_web: proto
	@npm run --prefix ./web build

dep:
	@cd go; dep ensure

run_server: dep
	@go run go/cmd/server/main.go

build_server: dep
	@go build -i -v -o go/cmd/server/main.go

run_client: dep
	@go run go/cmd/client/main.go

build_client: dep
	@go build -i -v -o go/cmd/client/main.go
