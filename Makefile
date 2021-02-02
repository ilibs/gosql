GO ?= go

.PHONY: test
test:
	@$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...