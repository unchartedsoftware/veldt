version=0.1.0

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  build         - build the dist binary"
	@echo "  clean         - clean the dist build"
	@echo ""
	@echo "  tools         - go gets a bunch of tools for dev"
	@echo "  deps          - pull and setup dependencies"
	@echo "  update_deps   - update deps lock file"

clean:
	@rm -rf ./bin

build: clean
	@go vet ./...
	@golint ./...
	@go build -o ./bin/prism.bin server/main.go

deps:
	@glock sync -n github.com/unchartedsoftware/prism < Glockfile

update_deps:
	@glock save -n github.com/unchartedsoftware/prism > Glockfile

tools:
	go get github.com/robfig/glock
	go get github.com/golang/lint/golint
