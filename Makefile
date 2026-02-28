BINARY = repotools
INSTALL_DIR = $(HOME)/bin

.PHONY: build install test clean

build:
	go build -o $(BINARY) ./cmd/repotools

install: build
	@rm -f $(INSTALL_DIR)/$(BINARY)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY)
