.PHONY: tests coverage dev build

EXCLUDE_DIRS := ./internal/ui ./internal/testkit

tests:
	@echo "Running tests..."
	@go test -v $(shell go list ./... | grep -vE "$(subst $(eval) ,|,$(EXCLUDE_DIRS))") \
	-cover -coverprofile=coverage.out
	@echo "Tests completed."

dev:
	@go run ./cmd/web

build:
	@go build -o ./bin/web ./cmd/web/