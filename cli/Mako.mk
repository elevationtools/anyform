
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/golang
endef

ifeq "$(MAKO_STAGE)" "main"

DEFAULT_TARGETS := genfiles/bin/anyform
DEFAULT_PREREQS := $(shell find . ../lib -name genfiles -prune -o -name "*.go") \
	| genfiles/bin
genfiles/bin/anyform:
	go build -o $@ .

genfiles/bin:
	mkdir -p $@

endif

include $(MAKO_ROOT)/component.mk

