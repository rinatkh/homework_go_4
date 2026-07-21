SHELL := /bin/bash
GO ?= go
GO_PACKAGES := ./...
UNIT_PACKAGES := ./internal/...
INTEGRATION_PACKAGES := ./test/integration/...
BIN_DIR := bin
COVERAGE_FILE ?= coverage.out
COVERAGE_THRESHOLD ?= 80.0
PACKAGE_FILE ?= $(BIN_DIR)/homework_go_4-linux-amd64.tar.gz
CMDS := 01_maps 02_methods 03_structs 04_common

.PHONY: help deps-check mod-check fmt fmt-check vet test test-unit test-integration test-race coverage coverage-check build package clean run-all ci compile alignment test-maps test-methods test-structs test-common $(addprefix run-,$(CMDS))

help:
	@echo "Available commands:"
	@echo "  make compile          - compile all packages without running tests"
	@echo "  make test-maps        - test map tasks"
	@echo "  make test-methods     - test method tasks"
	@echo "  make test-structs     - test struct tasks"
	@echo "  make test-common      - test common tasks"
	@echo "  make ci               - full local CI"

deps-check:
	$(GO) mod download
	$(GO) mod verify

mod-check:
	$(GO) mod tidy
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		git diff --exit-code -- go.mod; \
		if [ -f go.sum ]; then git diff --exit-code -- go.sum; fi; \
	else \
		echo "Skipping git diff because this directory is not a git repository"; \
	fi

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './$(BIN_DIR)/*')

fmt-check:
	@files="$$(gofmt -l $$(find . -name '*.go' -not -path './$(BIN_DIR)/*'))"; \
	if [ -n "$$files" ]; then echo "Go files are not formatted:"; echo "$$files"; exit 1; fi

vet:
	$(GO) vet $(GO_PACKAGES)

test: test-unit test-integration

test-unit:
	$(GO) test $(UNIT_PACKAGES)

test-integration:
	$(GO) test $(INTEGRATION_PACKAGES)

compile:
	$(GO) test -run '^$$' $(GO_PACKAGES)

test-maps:
	$(GO) test ./internal/maps/...

test-methods:
	$(GO) test ./internal/methods/...

test-structs:
	$(GO) test ./internal/structs/...

test-common:
	$(GO) test ./internal/common/...

test-race:
	$(GO) test -race $(UNIT_PACKAGES)

coverage:
	$(GO) test $(UNIT_PACKAGES) -covermode=atomic -coverprofile=$(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)

coverage-check: coverage
	@coverage="$$(go tool cover -func=$(COVERAGE_FILE) | awk '/^total:/ {gsub("%", "", $$3); print $$3}')"; \
	awk -v coverage="$$coverage" -v threshold="$(COVERAGE_THRESHOLD)" 'BEGIN { \
		if (coverage + 0 < threshold + 0) { printf "coverage %.1f%% is below threshold %.1f%%\n", coverage, threshold; exit 1 } \
		printf "coverage %.1f%% is enough; threshold %.1f%%\n", coverage, threshold; \
	}'

run-all:
	@for cmd in $(CMDS); do echo "== $$cmd =="; $(GO) run ./cmd/$$cmd; done

run-01_maps:
	$(GO) run ./cmd/01_maps
run-02_methods:
	$(GO) run ./cmd/02_methods
run-03_structs:
	$(GO) run ./cmd/03_structs
run-04_common:
	$(GO) run ./cmd/04_common

alignment:
	$(GO) run golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest ./internal/structs/...

build:
	@mkdir -p $(BIN_DIR)
	@for cmd in $(CMDS); do $(GO) build -o $(BIN_DIR)/$$cmd ./cmd/$$cmd; done

package: build
	@mkdir -p $(BIN_DIR)
	tar -czf $(PACKAGE_FILE) -C $(BIN_DIR) $(CMDS)

ci: deps-check mod-check fmt-check vet test-unit test-integration test-race coverage-check build package

clean:
	rm -rf $(BIN_DIR) $(COVERAGE_FILE)
