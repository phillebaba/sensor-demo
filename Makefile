.PHONY: all

VERSION := $(shell git describe --tags --always --dirty)

IMAGE_REGISTRY := "phillebaba"
IMAGE_NAME := "sensor-demo"
IMAGE_TAG_NAME := $(VERSION)
IMAGE_REPO := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG := $(IMAGE_REPO):$(IMAGE_TAG_NAME)

PLATFORMS := "linux/amd64,linux/arm64,linux/arm"

all: image

dep:
	@cd go; dep ensure

npm:
	@cd web; npm install

proto: npm dep
	@mkdir -p web/proto
	@protoc -I proto --plugin=protoc-gen-ts=./web/node_modules/.bin/protoc-gen-ts --js_out=import_style=commonjs,binary:web/proto --ts_out=service=grpc-web:web/proto --go_out=plugins=grpc:go/pkg/api proto/temperature.proto

run_web: proto
	@npm run --prefix ./web start

run_server: dep
	@go run go/cmd/server/main.go

run_client: dep
	@go run go/cmd/client/main.go

image:
	@docker build -t $(IMAGE) go
	@docker build -t $(IMAGE) -f web/docker/Dockerfile web

push: proto
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	@docker buildx rm cross
	@docker buildx create --use --name cross --platform $(PLATFORMS)
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE_TAG) --push go
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE)-web:$(TAG) --push -f web/docker/Dockerfile web
