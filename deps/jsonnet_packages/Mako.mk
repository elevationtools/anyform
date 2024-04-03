
include $(MAKO_ROOT)/dep.mk

define DEPS
	$(REPO_ROOT)/deps/jsonnet
	$(REPO_ROOT)/deps/jsonnet-bundler
endef

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	jsonnet-bundler install

smoketest:
	# Run a function we know is in one of the packages that's been installed.
	[ "/" = "$$(jsonnet -Se '(import "elevation/dirname.libsonnet")("/foo")')" ]

endif

clean:
	rm -rf vendor

