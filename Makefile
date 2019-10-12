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
	dep ensure -update

run_server:
	@go run go/cmd/server/main.go

build_server:
	@go build -i -v -o go/cmd/server/main.go

run_client:
	@go run go/cmd/client/main.go

build_client:
	@go build -i -v -o go/cmd/client/main.go
