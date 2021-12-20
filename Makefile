export tag=v1.0

build:  $(shell find . -name '*.go')
	echo "building httpserver binary"
	mkdir -p bin/
	CGO_ENABLED=0 GOARCH=amd64 go build -o bin/ ./httpserver ./server.go

release:
	echo "building httpserver container"
	docker build -t registry.cn-beijing.aliyuncs.com/doc01/httpserver:${tag} .

push: release
	echo "pushing doc01/httpserver"
	docker push registry.cn-beijing.aliyuncs.com/doc01/httpserver:${tag} .

deploy:
	echo "deploy httpserver to kubernets"
	kubectl apply -f spec/
