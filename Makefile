.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test -race -cover ./internal/...

.PHONY: docker-build
docker-build:
	@docker build -t mock-http-server:local .