BINARY = repotools
INSTALL_DIR = $(HOME)/bin

.PHONY: build install test ci clean

build:
	go build -o $(BINARY) ./cmd/repotools

install: build
	@rm -f $(INSTALL_DIR)/$(BINARY)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)

test:
	go test ./...

ci: fmt-check test

fmt-check:
	@test -z "$$(gofmt -l ./src/ ./cmd/)" || (echo "gofmt needed:"; gofmt -l ./src/ ./cmd/; exit 1)

clean:
	rm -f $(BINARY)
