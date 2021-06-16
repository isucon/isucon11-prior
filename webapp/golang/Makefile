
SRC:=$(shell find . -type f -name '*.go' -not -path 'vendor/*')


.PHONY: all
all: build ;

.PHONY: build
build: bin/webapp

bin/webapp: ${SRC}
	@go build -o bin/webapp .

.PHONY: run
run: build
	@./bin/webapp
