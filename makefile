.PHONY: all build dockerbuild dockerimages

ROOT_DIR = $(CURDIR)
BIN_DIR = $(ROOT_DIR)/bin
DOCKER_DIR = $(CURDIR)/dockerbuild
BIN_DIR_DOCKER = $(DOCKER_DIR)/bin


all : build

build:
	go fmt ./mdns-server/...
	go fmt ./mdns-client/...
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/mdns-server ./mdns-server/main.go
	go build -o $(BIN_DIR)/mdns-client ./mdns-client/main.go

dockerbuild:
	mkdir -p $(BIN_DIR_DOCKER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(BIN_DIR_DOCKER)/mdns-server ./mdns-server/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(BIN_DIR_DOCKER)/mdns-client ./mdns-client/main.go

dockerimages: dockerbuild
	cd $(DOCKER_DIR) && docker build -f Dockerfile_mdns-server . -t mdns-server:latest
	cd $(DOCKER_DIR) && docker build -f Dockerfile_mdns-client . -t mdns-client:latest