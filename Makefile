GO_CMD        = go
WORKDIR      := $(shell pwd)
NATIVEOS     := $(shell $(GO_CMD) version | awk -F '[ /]' '{print $$4}')
NATIVEARCH   := $(shell $(GO_CMD) version | awk -F '[ /]' '{print $$5}')
GO_PKGS      := $(shell $(GO_CMD) list ./... | grep -v -e "/vendor/" -e "/example")
GO_FILES     := $(shell find cmd -type f -name "*.go")
GO_TOOLS      = github.com/axw/gocov/gocov github.com/AlekSi/gocov-xml

BIN_DIR            = $(WORKDIR)/bin
TARGET             = target
TARGET_DIR         = $(WORKDIR)/$(TARGET)
INTEGRATION       := vsphere
SHORT_INTEGRATION := vsphere
BINARY_NAME        = nri-$(INTEGRATION)
CONTAINER_IMAGE    = $(PROJECT_NAME)-builder
CONTAINER          = $(PROJECT_NAME)

LINTER         = golangci-lint
LINTER_VERSION = 1.27.0
SNYK_BIN       = snyk-linux
SNYK_VERSION   = v1.361.3

all: build
build-local: clean test compile
build: bin
	@docker build --no-cache -t $(CONTAINER_IMAGE) .
	@docker run --privileged=true --name $(CONTAINER) $(CONTAINER_IMAGE) \
		/etc/init.d/docker start && make test compile 
	@docker cp $(CONTAINER):/go/src/$(PROJECT_NAME)/bin/$(BINARY_NAME) $(BIN_DIR); docker rm -f $(CONTAINER);

bin:
	@mkdir -p $@

clean:
	@echo "=== $(PROJECT_NAME) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml $(TARGET)

compile: compile-only
compile-only: deps
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@$(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-linux: deps
	@echo "=== $(PROJECT_NAME) === [ compile-linux    ]: building commands:"
	@GOOS=linux $(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-darwin: deps
	@echo "=== $(PROJECT_NAME) === [ compile-darwin    ]: building commands:"
	@GOOS=darwin $(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-windows: deps
	@echo "=== $(PROJECT_NAME) === [ compile-windows    ]: building commands:"
	@GOOS=windows $(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME).exe ./cmd/...


test: deps lint test-unit test-integration
test-unit:
	@echo "=== $(PROJECT_NAME) === [ unit-test        ]: running unit tests..."
	@gocov test $(GO_PKGS) | gocov-xml > coverage.xml

test-integration: compile
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	@$(GO_CMD) test -v -tags=integration ./integration-test/.

bin:
	@mkdir $(BIN_DIR)

test-security: bin deps
	@echo "=== $(PROJECT_NAME) === [ security-test        ]: running security tests..."
	@wget https://github.com/snyk/snyk/releases/download/$(SNYK_VERSION)/$(SNYK_BIN) -O $(BIN_DIR)/$(SNYK_BIN)
	@chmod +x $(BIN_DIR)/snyk-linux
	@$(BIN_DIR)/$(SNYK_BIN) auth $(SNYK_TOKEN)
	@$(BIN_DIR)/$(SNYK_BIN) test

lint: deps lint-deps
	@echo "=== $(PROJECT_NAME) === [ lint             ]: Validating source code running $(LINTER)..."
	@$(LINTER) run ./...


deps: tools deps-only
tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools ]: Installing tools required by the project..."
	@$(GO_CMD) get $(GO_TOOLS)
tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update ]: Updating tools required by the project..."
	@$(GO_CMD) get -u $(GO_TOOLS)
deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps ]: Installing package dependencies required by the project..."
	@$(GO_CMD) mod download
lint-deps:
	@echo "=== $(PROJECT_NAME) === [ lint-deps ]: Installing linting dependencies required by the project..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$($(GO_CMD) env GOPATH)/bin v$(LINTER_VERSION)

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