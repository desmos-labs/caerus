#!/usr/bin/make -f

BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)

export GO111MODULE = on

###############################################################################
###                                   All                                   ###
###############################################################################

all: lint test-unit build

###############################################################################
###                                Protobuf                                 ###
###############################################################################
protoVer=0.11.6
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user $(shell id -u):$(shell id -g) $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh
	$(MAKE) update-swagger-docs

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=master

proto-update-deps:
	$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update

.PHONY: proto-all proto-gen proto-lint proto-check-breaking proto-update-deps

###############################################################################
###                               Build flags                               ###
###############################################################################

build_tags = netgo

# These lines here are essential to include the muslc library for static linking of libraries
# (which is needed for the wasmvm one) available during the build. Without them, the build will fail.
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

ldflags =
ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

ldflags := $(strip $(ldflags))
BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

###############################################################################
###                                 Build                                   ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*pb*.go' -not -name '*mock*.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*pb*.go' -not -name '*mock*.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*pb*.go' -not -name '*mock*.go' | xargs goimports -w -local github.com/desmos-labs/caerus
.PHONY: format

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

stop-test-db:
	@echo "Stopping test database..."
	@docker stop apis-test-db || true && docker rm apis-test-db || true
.PHONY: stop-test-db

start-test-db: stop-test-db
	@echo "Starting test database..."
	@docker run --name apis-test-db -e POSTGRES_USER=caerus -e POSTGRES_PASSWORD=password -e POSTGRES_DB=caerus -d -p 6433:5432 postgres
.PHONY: start-test-db

start-test-chain:
	@echo "Starting test chain"
	@scripts/get_desmos_bin.sh 5.1.0
	@scripts/spawn_test_chain.sh

test-unit: start-test-db
	@echo "Executing unit tests..."
	@go test -p=1 -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit
