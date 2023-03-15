LDFLAGS = -ldflags '-w -s'

.PHONY: ayame

ayame:
	go build $(LDFLAGS) -o bin/$@ cmd/ayame/main.go

.PHONY: darwin linux

GOOS = $@
GOARCH = amd64

linux darwin:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/ayame-$(GOOS) cmd/ayame/main.go

check:
	go test ./...

clean:
	rm -rf ayame

.PHONY: lint

lint:
	golangci-lint run ./...

fmt:
	golangci-lint run ./... --fix

.PHONY: init

init:
	cp -n ayame.example.ini ayame.ini
