##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app


COPY ./ .

RUN go env -w GO111MODULE="on"; \
    go env -w GOPROXY="https://goproxy.cn,direct"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/httpserver -ldflags="-w -s" ./main.go ./server.go

##
## Deploy
##
FROM scratch

WORKDIR /app

COPY config config/

COPY --from=build /app/bin/httpserver /app/httpserver

EXPOSE 8000

ENV VERSION=1.0

LABEL lan="golang" app="httpserver"

ENTRYPOINT ["/app/httpserver"]