##
## Build
##

FROM golang:1.16-buster AS build
WORKDIR /app
COPY ["config", "pkg", "go.mod", "go.sum", "main.go", "server.go", "./"]

RUN ls; \
    go env -w GOPROXY="https://goproxy.cn,direct"; \
    CGO_ENABLED=0 GOARCH=amd64 go build -o bin/ ./main.go ./server.go


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