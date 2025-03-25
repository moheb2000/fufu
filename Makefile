VERSION=`git describe 2>/dev/null | cut -c 2- || echo "undefined"`
LDFLAGS=-ldflags="-X 'github.com/moheb2000/fufu.engineVersion=$(VERSION)'"
GORUN=go run $(LDFLAGS)

.PHONY: run
run:
	$(GORUN) ./cmd/engine
