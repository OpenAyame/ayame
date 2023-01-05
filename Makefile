LDFLAGS = -ldflags '-w -s'

.PHONY: ayame

ayame:
	go build $(LDFLAGS) -o $@

.PHONY: darwin linux

GOOS = $@
GOARCH = amd64

linux darwin:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o ayame-$(GOOS)

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
	cp -n ayame.example.yaml ayame.yaml
