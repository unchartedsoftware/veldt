version=0.1.0

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  build         - build the dist binary"
	@echo "  lint          - lint the source code"
	@echo "  test          - test the source code"
	@echo "  fmt           - format the code with gofmt"
	@echo "  clean         - clean the dist build"
	@echo ""
	@echo "  deps          - pull and install tool dependencies"

clean:
	@rm -rf ./build

lint:
	@govendor vet +local 
	@go list ./... | grep -v /vendor/ | xargs -L1 golint

test:
	@govendor test +local

fmt:
	@govendor fmt +local

build: clean lint
	@govendor build +local

deps:
	@go get github.com/kardianos/govendor
	@go get github.com/golang/lint/golint
	
