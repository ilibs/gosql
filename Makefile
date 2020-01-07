GO ?= go

.PHONY: test
test:
	@$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./... > tmp.out;
	@cat tmp.out; \
	if grep -q "^--- FAIL" tmp.out; then \
		rm tmp.out; \
		exit 1; \
	elif grep -q "build failed" tmp.out; then \
		rm tmp.out; \
		exit 1; \
	elif grep -q "setup failed" tmp.out; then \
		rm tmp.out; \
		exit 1; \
	fi;