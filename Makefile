VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_USER ?= $(shell whoami)@$(shell hostname)
PKG := github.com/user-cube/$PROJECT_NAME
LDFLAGS := -ldflags "-X $(PKG)/cmd.Version=$(VERSION) -X $(PKG)/cmd.BuildDate=$(BUILD_DATE) -X $(PKG)/cmd.GitCommit=$(COMMIT) -X $(PKG)/cmd.BuildUser=$(BUILD_USER)"

.PHONY: all
all: clean build

.PHONY: build
build:
	@echo "Building $PROJECT_NAME $(VERSION) ($(COMMIT))"
	@go build $(LDFLAGS) -o $PROJECT_NAME main.go

.PHONY: install
install:
	@echo "Installing $PROJECT_NAME $(VERSION) to GOPATH"
	@go install $(LDFLAGS)

.PHONY: clean
clean:
	@echo "Cleaning build artifacts"
	@rm -f $PROJECT_NAME
	@rm -rf dist

.PHONY: test
test:
	@echo "Running tests"
	@go test -v ./...

.PHONY: lint
lint:
	@echo "Running linters"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping lint"; \
	fi

.PHONY: release
release: clean
	@echo "Creating release with GoReleaser $(VERSION)"
	@if ! command -v goreleaser > /dev/null; then \
		echo "Error: goreleaser not found. Install with 'go install github.com/goreleaser/goreleaser@latest'"; \
		exit 1; \
	fi
	@VERSION=$(VERSION) GIT_COMMIT=$(COMMIT) BUILD_DATE=$(BUILD_DATE) goreleaser release --clean

.PHONY: release-snapshot
release-snapshot: clean
	@echo "Creating snapshot release with GoReleaser (no publish)"
	@if ! command -v goreleaser > /dev/null; then \
		echo "Error: goreleaser not found. Install with 'go install github.com/goreleaser/goreleaser@latest'"; \
		exit 1; \
	fi
	@VERSION=$(VERSION) GIT_COMMIT=$(COMMIT) BUILD_DATE=$(BUILD_DATE) goreleaser release --snapshot --clean

.PHONY: build-release
build-release: clean
	@echo "Building release version of $PROJECT_NAME $(VERSION) ($(COMMIT))"
	@go build $(LDFLAGS) -o $PROJECT_NAME main.go
	@echo "Built $PROJECT_NAME binary with release information"
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Run ./$PROJECT_NAME version to verify"

.PHONY: help
help:
	@echo "$PROJECT_NAME Makefile"
	@echo "---------------"
	@echo "Available targets:"
	@echo "  all              - Clean and build $PROJECT_NAME"
	@echo "  build            - Build the $PROJECT_NAME binary"
	@echo "  install          - Install $PROJECT_NAME to your GOPATH/bin"
	@echo "  clean            - Remove built binary and dist directory"
	@echo "  test             - Run tests"
	@echo "  lint             - Run linters (requires golangci-lint)"
	@echo "  release          - Create a full release using GoReleaser"
	@echo "  release-snapshot - Create a local release snapshot for testing (no publish)"
	@echo "  build-release    - Build $PROJECT_NAME binary with release information"
	@echo "  help             - Show this help message"