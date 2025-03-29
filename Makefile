VERSION=`git describe 2>/dev/null || echo "vundefined" | cut -c 2-`
LDFLAGS=-ldflags="-X 'github.com/moheb2000/fufu.engineVersion=$(VERSION)'"
GORUN=go run $(LDFLAGS)
GOBUILD=go build $(LDFLAGS)
BINARY_NAME=fufu
BUILD_PATH=build/$(BINARY_NAME)
SOURCE_PATH=./cmd/engine

.PHONY: run
run:
	$(GORUN) ./cmd/engine

.PHONY: build/linux
build/linux:
	CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 $(GOBUILD) -tags static -o=$(BUILD_PATH)-$(VERSION) $(SOURCE_PATH)
	cp -r assets build
	cp -r config.json build
	cp -r main.lua build

.PHONY: clean
clean:
	rm -rf build
