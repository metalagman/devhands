lint:
	@echo "Running linter checks"
	golangci-lint run
	go-arch-lint check

test:
	@echo "Running UNIT tests"
	@go clean -testcache
	go test -cover -race -short ./... | { grep -v 'no test files'; true; }

cover:
	@echo "Running test coverage"
	@go clean -testcache
	go test -cover -coverprofile=coverage.out -race -short ./internal/... | grep -v 'no test files'
	go tool cover -html=coverage.out

generate:
	@echo "Generating mocks"
	go generate ./...

.PHONY: build
build:
	@echo "Building the app to the bin dir"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/devhands ./cmd/devhands/*.go
