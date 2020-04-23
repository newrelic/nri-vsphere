WORKDIR      := $(shell pwd)
NATIVEOS     := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH   := $(shell go version | awk -F '[ /]' '{print $$5}')
GO_PKGS      := $(shell go list ./... | grep -v -e "/vendor/" -e "/example")
GO_FILES     := $(shell find cmd -type f -name "*.go")

GO_CMD        = go
GOLINTER      = golangci-lint

BIN_DIR    = $(WORKDIR)/bin
TARGET = target
TARGET_DIR       = $(WORKDIR)/$(TARGET)
INTEGRATION  := vmware-vsphere
SHORT_INTEGRATION  := vsphere
BINARY_NAME   = nri-$(INTEGRATION)

GOTOOLS       = github.com/kardianos/govendor gopkg.in/alecthomas/gometalinter.v2 github.com/axw/gocov/gocov github.com/AlekSi/gocov-xml


all: build
build: clean test compile


clean:
	@echo "=== $(INTEGRATION) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml $(TARGET)

compile: compile-only
compile-only: deps
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-linux: deps
	@echo "=== $(PROJECT_NAME) === [ compile-linux    ]: building commands:"
	@GOOS=linux go  build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-darwin: deps
	@echo "=== $(PROJECT_NAME) === [ compile-linux    ]: building commands:"
	@GOOS=darwin go  build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-windows: deps
	@echo "=== $(PROJECT_NAME) === [ compile-linux    ]: building commands:"
	@GOOS=windows go  build -o $(BIN_DIR)/$(BINARY_NAME).exe ./cmd/...


test: deps test-unit test-integration
test-unit:
	@echo "=== $(PROJECT_NAME) === [ unit-test        ]: running unit tests..."
	@gocov test $(GO_PKGS) | gocov-xml > coverage.xml
	ls

test-integration:
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	@docker-compose -f ./integration-test/docker-compose.yml up -d --build
	@go test -v -tags=integration ./integration-test/. || (ret=$$?; docker-compose -f ./integration-test/docker-compose.yml  down && exit $$ret)
	@docker-compose -f ./integration-test/docker-compose.yml  down
lint: deps
	@echo "=== $(PROJECT_NAME) === [ lint             ]: Validating source code running $(GOLINTER)..."
	@$(GOLINTER) run ./...


deps: tools deps-only
tools: check-version
	@echo "=== $(INTEGRATION) === [ tools ]: Installing tools required by the project..."
	@go get $(GOTOOLS)
tools-update: check-version
	@echo "=== $(INTEGRATION) === [ tools-update ]: Updating tools required by the project..."
	@go get -u $(GOTOOLS)
deps-only:
	@echo "=== $(INTEGRATION) === [ deps ]: Installing package dependencies required by the project..."
	@govendor sync


check-version:
ifdef GOOS
ifneq "$(GOOS)" "$(NATIVEOS)"
	$(error GOOS is not $(NATIVEOS). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
endif
ifdef GOARCH
ifneq "$(GOARCH)" "$(NATIVEARCH)"
	$(error GOARCH variable is not $(NATIVEARCH). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
endif

# Import fragments
include package.mk

.PHONY: all build