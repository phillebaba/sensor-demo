.PHONY: all proto build_server build_client

IMAGE ?= phillebaba/sensor-demo
VERSION ?= $(shell git describe --tags --always --dirty)
#TAG ?= $(VERSION)
TAG ?= "latest"
PLATFORMS ?= "linux/amd64,linux/arm64,linux/arm"

all: build_server build_client

dep:
	@cd go; dep ensure

npm:
	@cd web; npm install

proto: npm dep
	@mkdir -p web/proto
	@protoc -I proto --plugin=protoc-gen-ts=./web/node_modules/.bin/protoc-gen-ts --js_out=import_style=commonjs,binary:web/proto --ts_out=service=grpc-web:web/proto --go_out=plugins=grpc:go/pkg/api proto/temperature.proto

run_web: proto
	@npm run --prefix ./web start

docker_push_web: proto
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	@docker buildx rm cross
	@docker buildx create --use --name cross --platform $(PLATFORMS)
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE)-web:$(TAG) --push -f web/docker/Dockerfile web

run_server: dep
	@go run go/cmd/server/main.go

run_client: dep
	@go run go/cmd/client/main.go

docker_push_go:
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	@docker buildx rm cross
	@docker buildx create --use --name cross --platform $(PLATFORMS)
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE):$(TAG) --push go
