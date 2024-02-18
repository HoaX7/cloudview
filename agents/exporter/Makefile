build:
	@go build -o bin/agents/exporter

build-linux:
# -ldflags "-s -w" you can use this to make the binary a little smaller but it is not recommended
	@GOARCH=amd64 GOOS=linux go build -o exporter

build-windows:
	@GOARCH=amd64 GOOS=windows go build -o bin/linux/exporter.exe

run: build
	@./bin/agents/exporter

test:
	@go test -v ./..