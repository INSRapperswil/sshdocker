.PHONY: all dev-start dev-up dev-down build clean

all: build

dev-start: dev-down dev-up
		go run ./cmd/sshdocker -c sshdocker-dev -u admin -p password -s /bin/bash

dev-up:
		docker run -d -it --name sshdocker-dev ubuntu:focal bash

dev-down:
		docker rm -f sshdocker-dev

build:
		mkdir -p bin
		CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/sshdocker ./cmd/sshdocker

clean:
		rm -rf bin