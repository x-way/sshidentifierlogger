.PHONY: all build format lint clean

all: format lint build

build:
	go build .
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o sshidentifierlogger.aarch64 .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o sshidentifierlogger.amd64 .

format:
	go install mvdan.cc/gofumpt@latest
	gofumpt -w .

lint:
	go vet -v ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -checks all ./...

clean:
	rm -f sshidentifierlogger sshidentifierlogger.aarch64 sshidentifierlogger.amd64
