REVISION:=$(shell git rev-parse --short HEAD 2> /dev/null || cat REVISION)

GO_FILES:=$(shell find . -type f -name '*.go' -print)

GO_LDFLAGS_VERSION:=-X 'main.REVISION=${REVISION}'
GO_LDFLAGS_SYMBOL:=
ifdef RELEASE
	GO_LDFLAGS_SYMBOL:=-w -s
endif
GO_LDFLAGS_STATIC:=
ifdef RELEASE
	GO_LDFLAGS_STATIC:=-extldflags '-static'
endif
GO_LDFLAGS:=$(GO_LDFLAGS_VERSION) $(GO_LDFLAGS_SYMBOL) $(GO_LDFLAGS_STATIC)
GO_BUILD_TAGS:=debug
ifdef RELEASE
	GO_BUILD_TAGS:=release
endif
GO_BUILD_RACE:=-race
ifdef RELEASE
	GO_BUILD_RACE:=
endif
GO_BUILD_STATIC:=
ifdef RELEASE
	GO_BUILD_STATIC:=-a -installsuffix netgo
	GO_BUILD_TAGS:=$(GO_BUILD_TAGS),netgo
endif
GO_BUILD:=-tags=$(GO_BUILD_TAGS) $(GO_BUILD_RACE) $(GO_BUILD_STATIC) -ldflags "$(GO_LDFLAGS)"

.PHONY: all
all: build

.PHONY: build
build: bin/benchmarker

.PHONY: clean
clean:
	@$(RM) -r bin

bin/benchmarker: ${GO_FILES}
	@go build -o $@ ${GO_BUILD} .
