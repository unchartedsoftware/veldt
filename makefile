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
	@go test $(shell glide novendor)

fmt:
	@go fmt $(shell glide novendor)

build: lint
	@go build $(shell glide novendor)

install:
	@go get github.com/golang/lint/golint
	@go get github.com/Masterminds/glide
	@glide install
