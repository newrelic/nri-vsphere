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
BINS_DIR           = $(TARGET_DIR)/bin/linux_amd64
INTEGRATION       := vsphere
SHORT_INTEGRATION := vsphere
BINARY_NAME        = nri-$(INTEGRATION)
CONTAINER_IMAGE    = $(PROJECT_NAME)-builder
CONTAINER          = $(PROJECT_NAME)
CONTAINER_PATH     = /go/src/$(PROJECT_NAME)

LINTER         = golangci-lint
LINTER_VERSION = 1.27.0
SNYK_BIN       = snyk-linux
SNYK_VERSION   = v1.361.3

all: build
build-local: clean compile test tidy
build: build-container-image delete-container test-container delete-container

build-container-image:
	@echo "=== $(PROJECT_NAME) === [ $@ ]: Building the container image"
	@docker build --no-cache -t $(CONTAINER_IMAGE) -f Dockerfile.test .

test-container:
	@echo "=== $(PROJECT_NAME) === [ $@ ]: Testing the integration"
	@echo "make test" | docker run --name $(CONTAINER) -i $(CONTAINER_IMAGE)
	@docker cp $(CONTAINER):$(CONTAINER_PATH)/coverage.xml .

compile-container: bin
	@echo "=== $(PROJECT_NAME) === [ $@ ]: Creating the binaries"
	@echo "make create-bins" | docker run --name $(CONTAINER) -i -e VERSION $(CONTAINER_IMAGE)
	@docker cp $(CONTAINER):$(CONTAINER_PATH)/$(TARGET)/bin/linux_amd64/$(BINARY_NAME) $(BINS_DIR)

delete-container:
	-docker rm -f $(CONTAINER) 2>/dev/null

bin:
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(BINS_DIR)

clean:
	@echo "=== $(PROJECT_NAME) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv $(BIN_DIR) $(BINS_DIR) coverage.xml $(TARGET)

compile: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@$(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-linux: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-linux    ]: building commands:"
	@GOOS=linux $(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-darwin: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-darwin    ]: building commands:"
	@GOOS=darwin $(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...
compile-windows: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-windows    ]: building commands:"
	@GOOS=windows $(GO_CMD) build -o $(BIN_DIR)/$(BINARY_NAME).exe ./cmd/...


test: deps lint test-unit test-integration
test-unit: compile
	@echo "=== $(PROJECT_NAME) === [ unit-test        ]: running unit tests..."
	@gocov test $(GO_PKGS) | gocov-xml > coverage.xml

test-integration: compile
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	@$(GO_CMD) test -v -tags=integration ./integration-test/.

test-security: bin deps
	@echo "=== $(PROJECT_NAME) === [ security-test        ]: running security tests..."
	@wget https://github.com/snyk/snyk/releases/download/$(SNYK_VERSION)/$(SNYK_BIN) -O $(BIN_DIR)/$(SNYK_BIN)
	@chmod +x $(BIN_DIR)/snyk-linux
	@$(BIN_DIR)/$(SNYK_BIN) auth $(SNYK_TOKEN)
	@$(BIN_DIR)/$(SNYK_BIN) test

lint: lint-deps
	@echo "=== $(PROJECT_NAME) === [ lint             ]: Validating source code running $(LINTER)..."
	@$(LINTER) run ./...

tidy:
	@echo "=== $(PROJECT_NAME) === [ tidy ]: Tidying up go mod..."
	@$(GO_CMD) mod tidy

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
include $(CURDIR)/build/ci.mk
include $(CURDIR)/build/release.mk

.PHONY: all build