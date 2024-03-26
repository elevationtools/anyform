
export REPO_ROOT := $(CURDIR)

.PHONY: build
build: cli/genfiles/anyform

cli/genfiles/anyform: $(shell find cli lib -name genfiles -prune -o -print)
	. ./activate.sh && mako -C cli
	
.PHONY: smoketest
smoketest: cli/genfiles/anyform
	cli/genfiles/anyform up

