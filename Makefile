VERSION = 19.08.0

LDFLAGS = -ldflags '-X main.AyameVersion=$(VERSION)'
GO111MODULE = on

.PHONY: ayame

ayame:
	GO111MODULE=$(GO111MODULE) go build $(LDFLAGS) -o $@


.PHONY: darwin linux

GOOS = $@
GOARCH = amd64

linux darwin:
	GO111MODULE=$(GO111MODULE) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o ayame-$(GOOS)

check:
	GO111MODULE=$(GO111MODULE) go test ./...

clean:
	rm -rf ayame

.PHONY: lint
lint:
	golangci-lint run ./...

fmt:
	golangci-lint run ./... --fix
