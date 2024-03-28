
export REPO_ROOT := $(CURDIR)

.PHONY: build
build: cli/genfiles/anyform

cli/genfiles/anyform: $(shell find cli lib -name genfiles -prune -o -print)
	. ./activate.sh && mako -C cli
	
.PHONY: examples_basic
examples_basic:
	. ./activate.sh && mako -C examples/basic

