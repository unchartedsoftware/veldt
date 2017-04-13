version=0.1.0

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
	@go vet $(shell glide novendor)
	@go list ./... | grep -v /vendor/ | xargs -L1 golint

test:
	@ginkgo -r -cover -covermode=count
	@rm -rf "coverage" && mkdir "coverage"
	@find . -name "*.coverprofile" | xargs gocovmerge > "coverage/coverage.out"
	@find . -name "*.coverprofile" -type f -delete
	@go tool cover -html="coverage/coverage.out" -o "coverage/coverage.html"

fmt:
	@go fmt $(shell glide novendor)

build: lint
	@go build $(shell glide novendor)

install:
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/Masterminds/glide
	@go get -u github.com/onsi/ginkgo/ginkgo
	@go get -u github.com/wadey/gocovmerge
	@glide install
