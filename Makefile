.PHONY: all build-frontend build-go build clean fmt test

all: build

fmt:
	go fmt ./...

test: fmt
	go test ./...

build-frontend:
	cd frontend && npm install && npm run build
	mkdir -p pkg/api/ui
	cp -r frontend/dist/* pkg/api/ui/

build-go: fmt test
	go build -o bin/sm3-migration-tool main.go

build: build-frontend build-go

container-build:
	docker build -t quay.io/acaglio/sm3-migration-tool:latest -f Containerfile .

container-push:
	docker push quay.io/acaglio/sm3-migration-tool:latest

push: build container-build container-push

clean:
	rm -rf frontend/dist
	rm -rf pkg/api/ui
	rm -f sm3-migration-tool
