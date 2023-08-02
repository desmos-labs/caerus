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

stop-test-db:
	@echo "Stopping test database..."
	@docker stop apis-test-db || true && docker rm apis-test-db || true
.PHONY: stop-test-db

start-test-db: stop-test-db
	@echo "Starting test database..."
	@docker run --name apis-test-db -e POSTGRES_USER=apis -e POSTGRES_PASSWORD=password -e POSTGRES_DB=apis -d -p 6433:5432 postgres
.PHONY: start-test-db

start-test-chain:
	@echo "Starting test chain"
	@scripts/get_desmos_bin.sh 5.1.0
	@scripts/spawn_test_chain.sh

test-unit: start-test-db
	@echo "Executing unit tests..."
	@go test -p=1 -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit
