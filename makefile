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
	@echo "  deps          - pull and setup dependencies"
	@echo "  update_deps   - update deps lock file"

clean:
	@rm -rf ./build

lint:
	@go vet ./...
	@golint ./...

test:
	@go test ./...

fmt:
	@gofmt -l -w .

build: clean lint
	@go build ./...

deps:
	@go get github.com/robfig/glock
	@go get github.com/golang/lint/golint
	@glock sync -n github.com/unchartedsoftware/prism < Glockfile

update_deps:
	@glock save -n github.com/unchartedsoftware/prism > Glockfile
