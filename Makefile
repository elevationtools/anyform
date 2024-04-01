# All targets downloads and activate the required dependency versions, including
# golang.  The only requirements needed a priori are:
#	- GNU make
# - curl
# - bash
# - git

export REPO_ROOT := $(CURDIR)

.DEFAULT_GOAL := build

# Build on the local machine.
.PHONY: build
build: submodules
	. ./activate.sh && mako -C module/cli

# Build within a docker container.
.PHONY: docker_build
docker_build: submodules
	. ./activate.sh && mako -C build local

.PHONY: examples_basic
examples_basic: submodules
	. ./activate.sh && mako -C examples/basic

.PHONY: clean
clean:
	-chmod -R u+w deps/gopath  # Annoyingly needed to be able to "rm -rf"
	rm -rf $(shell find . -name genfiles -type d -prune -print)
	git submodule deinit --all

.PHONY: submodules
submodules: deps/mako/README.md
deps/mako/README.md:
	git submodule update --init --recursive

