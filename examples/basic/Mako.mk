
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/jsonnet
	$(REPO_ROOT)/module/cli
endef

ifeq "$(MAKO_STAGE)" "main"

DEFAULT_TARGETS := help
DEFAULT_PREREQS :=
.PHONY: help
help:
	@echo "Targets: up down spec"

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

clean:
	-rm -rf prod/*/genfiles

include $(MAKO_ROOT)/component.mk

