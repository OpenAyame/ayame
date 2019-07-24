VERSION=19.07.0

ayame: *.go
	GO111MODULE=on go build -ldflags '-X main.AyameVersion=${VERSION}' -o $@

.PHONY: all
all: ayame

darwin-build: *.go
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -ldflags '-X main.AyameVersion=${VERSION}' -o ayame-darwin
linux-build:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags '-s -w -X main.AyameVersion=${VERSION}' -o ayame-linux

check:
	GO111MODULE=on go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf ayame
