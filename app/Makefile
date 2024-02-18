build:
	@go build -o bin/app

build-linux:
	@GOARCH=amd64 GOOS=linux go build -o bin/app

run: build
	@./bin/app

test:
	@go test -v ./..