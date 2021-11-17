export tag=v1.0

build:  $(shell find . -name '*.go')
	echo "building httpserver binary"
	mkdir -p bin/
	CGO_ENABLED=0 GOARCH=amd64 go build -o bin/ ./httpserver

release:
	echo "building httpserver container"
	docker build -t cncamp/httpserver:${tag} .

push: release
	echo "pushing cncamp/httpserver"
	docker push cncamp/httpserver:v1.0