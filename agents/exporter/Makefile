build:
	@go build -o bin/agents/exporter

build-linux:
	@GOARCH=amd64 GOOS=linux go build -o exporter

build-windows:
	@GOARCH=amd64 GOOS=windows go build -o bin/linux/exporter.exe

run: build
	@./bin/agents/exporter

test:
	@go test -v ./..