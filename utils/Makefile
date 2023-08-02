###############################################################################
###                                   All                                   ###
###############################################################################

all: lint test-unit

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
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*_mock.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*_mock.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*_mock.go' | xargs goimports -w -local github.com/desmos-labs/caerus
.PHONY: format

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test-unit:
	@echo "Executing unit tests..."
	@go test -p=1 -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit
