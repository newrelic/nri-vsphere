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

all: build-local
build-local: clean compile test tidy

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
tools-vcsim-run: clean compile-linux
	@echo "=== $(PROJECT_NAME) === Running vcsim with an agent:"
	@if [ "$(NRIA_LICENSE_KEY)" = "" ]; then \
	    echo "Error: missing required env-var: NRIA_LICENSE_KEY\n" ;\
        exit 1 ;\
	fi
	@docker-compose -f tools/docker-compose.yml up -d --build
tools-vcsim-stop:
	@echo "=== $(PROJECT_NAME) === Stopping vcsim with agent:"
	@docker-compose -f tools/docker-compose.yml down


test: deps test-unit test-integration
test-unit: deps compile
	@echo "=== $(PROJECT_NAME) === [ unit-test        ]: running unit tests..."
	@gocov test $(GO_PKGS) | gocov-xml > coverage.xml

test-integration: compile
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	@$(GO_CMD) test -v -tags=integration ./integration-test/.

tidy:
	@echo "=== $(PROJECT_NAME) === [ tidy ]: Tidying up go mod..."
	@$(GO_CMD) mod tidy

deps: tools deps-only
tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools ]: Installing tools required by the project..."
	@$(GO_CMD) install $(GO_TOOLS)
tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update ]: Updating tools required by the project..."
	@$(GO_CMD) get -u $(GO_TOOLS)
deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps ]: Installing package dependencies required by the project..."
	@$(GO_CMD) mod download

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
include $(CURDIR)/build/ci.mk
include $(CURDIR)/build/release.mk

.PHONY: all build