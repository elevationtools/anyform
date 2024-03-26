
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/golang
endef

ifeq "$(MAKO_STAGE)" "main"

DEFAULT_TARGETS := genfiles/anyform
DEFAULT_PREREQS := $(shell find . ../lib -name genfiles -prune -o -name "*.go")
genfiles/anyform:
	go build -o $@ .

endif

include $(MAKO_ROOT)/component.mk

