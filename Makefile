TESTPKGS = $(shell go list ./... | grep -v cmd)

.PHONY: test/coverage-unit
test/coverage-unit: ## run unit test and measure test coverage
	@go test -race -p=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage.txt $(TESTPKGS)

.PHONY: codecov
codecov: SHELL=/usr/bin/env bash
codecov: ## send coverage result
	bash <(curl -s https://codecov.io/bash) -Z -F ${CODECOV_FLAG}
