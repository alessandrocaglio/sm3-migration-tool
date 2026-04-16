.PHONY: all build-frontend build-go build clean

all: build

build-frontend:
	cd frontend && npm install && npm run build
	mkdir -p pkg/api/ui
	cp -r frontend/dist/* pkg/api/ui/

build-go:
	go build -o sm3-migration-tool main.go

build: build-frontend build-go

container-build:
	docker build -t quay.io/acaglio/sm3-migration-tool:latest -f Containerfile .

clean:
	rm -rf frontend/dist
	rm -rf pkg/api/ui
	rm -f sm3-migration-tool
