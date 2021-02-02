GO ?= go

.PHONY: test
test:
	$(GO) test -short -v -coverprofile=cover.out ./...

.PHONY: cover
cover:
	$(GO) tool cover -func=cover.out -o cover_total.out
	$(GO) tool cover -html=cover.out -o cover.html