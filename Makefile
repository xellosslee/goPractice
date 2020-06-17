run:
	@go run main.go bin/server.exe

build:
	@go build -o bin/server.exe main.go
	@bin/server.exe