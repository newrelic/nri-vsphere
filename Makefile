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
INTEGRATION  := vsphere
SHORT_INTEGRATION  := vsphere
BINARY_NAME   = nri-$(INTEGRATION)

GOTOOLS       = github.com/kardianos/govendor  github.com/axw/gocov/gocov github.com/AlekSi/gocov-xml

SNYK_VERSION  = v1.361.3
SNYK_BIN = snyk-linux

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
	@echo "=== $(PROJECT_NAME) === [ compile-darwin    ]: building commands:"
	@GOOS=darwin go  build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-windows: deps
	@echo "=== $(PROJECT_NAME) === [ compile-windows    ]: building commands:"
	@GOOS=windows go  build -o $(BIN_DIR)/$(BINARY_NAME).exe ./cmd/...


test: deps lint test-unit test-integration
test-unit:
	@echo "=== $(PROJECT_NAME) === [ unit-test        ]: running unit tests..."
	@gocov test $(GO_PKGS) | gocov-xml > coverage.xml

test-integration: compile
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	@go test -v -tags=integration ./integration-test/.

bin:
	@mkdir $(BIN_DIR)

test-security: bin deps
	@echo "=== $(PROJECT_NAME) === [ security-test        ]: running security tests..."
	@wget https://github.com/snyk/snyk/releases/download/$(SNYK_VERSION)/$(SNYK_BIN) -O $(BIN_DIR)/$(SNYK_BIN)
	@chmod +x $(BIN_DIR)/snyk-linux
	@$(BIN_DIR)/$(SNYK_BIN) auth $(SNYK_TOKEN)
	@$(BIN_DIR)/$(SNYK_BIN) test

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
	@$(GOPATH)/bin/govendor sync


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