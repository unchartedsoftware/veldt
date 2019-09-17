version=0.1.0

_allpackages = $(shell go list ./...)
# memoize allpackages, so that it's executed only once and only if used
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  build         - build the source code"
	@echo "  lint          - lint the source code"
	@echo "  test          - test the source code"
	@echo "  fmt           - format the code with gofmt"
	@echo "  install       - install dependencies"

lint:
	@go vet $(allpackages)
	@golint $(allpackages)

test:
	@ginkgo -r -cover -covermode=count
	@rm -rf "coverage" && mkdir "coverage"
	@find . -name "*.coverprofile" | xargs gocovmerge > "coverage/coverage.out"
	@find . -name "*.coverprofile" -type f -delete
	@go tool cover -html="coverage/coverage.out" -o "coverage/coverage.html"

fmt:
	@go fmt $(allpackages)

build: lint
	@go build $(allpackages)

install:
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/onsi/ginkgo/ginkgo
	@go get -u github.com/wadey/gocovmerge
	@go mod vendor
