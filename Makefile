run:
	@go run main.go bin/server.exe

build:
	@go build -o bin/server.exe main.go
	@bin/server.exe

build-linux:
	@set GOOS=linux&& set GOARCH=amd64&& go build -v -o bin/server main.go