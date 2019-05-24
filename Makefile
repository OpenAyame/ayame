VERSION=19.02.1

ayame: *.go
	GO111MODULE=on go build -ldflags '-X main.AyameVersion=${VERSION}' -o $@

.PHONY: all
all: ayame proto

linux-build:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags '-s -w -X main.AyameVersion=${VERSION}' -o ayame

check:
	GO111MODULE=on go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf ayame
