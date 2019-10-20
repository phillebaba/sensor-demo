.PHONY: all proto

VERSION := $(shell git describe --tags --always --dirty)

IMAGE_REGISTRY := "phillebaba"
IMAGE_NAME := "sensor-demo"
IMAGE_TAG_NAME := $(VERSION)
IMAGE_REPO := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG_GO := $(IMAGE_REPO):$(IMAGE_TAG_NAME)
IMAGE_TAG_WEB := $(IMAGE_REPO)-web:$(IMAGE_TAG_NAME)

PLATFORMS := "linux/amd64,linux/arm64,linux/arm"

all: image

proto:
	@mkdir -p web/proto
	@docker run --rm -v ${PWD}:${PWD} -w ${PWD} phillebaba/protoc-docker:6a297cf -I./proto --js_out=import_style=commonjs,binary:web/proto --ts_out=service=grpc-web:web/proto --go_out=plugins=grpc:go/pkg/api proto/temperature.proto

run_web: proto
	@cd web; npm install
	@npm run --prefix ./web start

dep:
	@go get github.com/golang/dep/cmd/dep
	@cd go; dep ensure -vendor-only

run_server: dep proto
	@go run go/cmd/server/main.go

run_client: dep proto
	@go run go/cmd/client/main.go

image:
	@docker build -t $(IMAGE_TAG_GO) go
	@docker build -t $(IMAGE_TAG_WEB) -f web/docker/Dockerfile web

push: proto
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	@docker buildx create --use --name cross --platform $(PLATFORMS) --node cross0
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE_TAG_GO) --push go
	@docker buildx build --platform $(PLATFORMS) -t $(IMAGE_TAG_WEB) --push -f web/docker/Dockerfile web
