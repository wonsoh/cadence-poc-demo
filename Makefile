PROJECT_ROOT = code.uber.internal/wonsoh/hello-world

PROGS = hello-world
SERVICES = hello-world

THRIFTRW_PLUGINS = go.uber.org/yarpc/encoding/thrift/thriftrw-plugin-yarpc
THRIFTRW_SRCS = $(shell find idl -name \*.thrift)

# Modifying any non-vendor go files will rebuild PROGS
hello-world: $(shell find . -not -path "./vendor/*" -name "*.go") Gopkg.lock

-include go-build/rules.mk

lint:: errcheck staticcheck unused

go-build/rules.mk:
	git submodule update --init

.PHONY: run
run: bins
	@./hello-world
