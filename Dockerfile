##
## Build
##

FROM golang:1.16-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY pkg/httpserver ./httpserver

RUN ls; \
    go env -w GOPROXY="https://goproxy.cn,direct"; \
    CGO_ENABLED=0 GOARCH=amd64 go build -o bin/ ./httpserver


##
## Deploy
##
FROM scratch
WORKDIR /

COPY --from=build /pkg/httpserver /httpserver

EXPOSE 8000

ENV VERSION=1.0

LABEL lan="golang" app="httpserver"

ENTRYPOINT ["/httpserver"]