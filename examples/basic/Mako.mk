
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/jsonnet
	$(REPO_ROOT)/module/cli
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

.PHONY: spec
spec:
	cd prod/tuesday && anyform spec
	

endif

include $(MAKO_ROOT)/base.mk

