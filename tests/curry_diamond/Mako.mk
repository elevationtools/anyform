
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/jsonnet
	$(REPO_ROOT)/module/cli
endef

ifeq "$(MAKO_STAGE)" "main"

DEFAULT_TARGETS := test
DEFAULT_PREREQS :=
.PHONY: test
test:
	./run_tests

endif

clean:
	-rm -rf prod/*/genfiles prod/*/output

include $(MAKO_ROOT)/component.mk

