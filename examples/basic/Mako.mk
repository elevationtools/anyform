
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/jsonnet
	$(REPO_ROOT)/cli
endef

ifeq "$(MAKO_STAGE)" "main"

#DEFAULT_TARGETS := smoketest
#DEFAULT_PREREQS :=
.PHONY: up
up:
	cd prod/tuesday && anyform up

.PHONY: down
down:
	cd prod/tuesday && anyform down
	

endif

include $(MAKO_ROOT)/base.mk

