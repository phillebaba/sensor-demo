FROM golang:1.13-alpine AS build

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/phillebaba/sensor-demo/go/
WORKDIR /go/src/github.com/phillebaba/sensor-demo/go/
RUN dep ensure -vendor-only

COPY . /go/src/github.com/phillebaba/sensor-demo/go/
RUN CGO_ENABLED=0 go build -o /bin/sensor-demo-client cmd/client/main.go
RUN CGO_ENABLED=0 go build -o /bin/sensor-demo-server cmd/server/main.go

FROM scratch
COPY --from=build /bin/sensor-demo-client /bin/sensor-demo-client
COPY --from=build /bin/sensor-demo-server /bin/sensor-demo-server
