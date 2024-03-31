
include $(MAKO_ROOT)/util.mk

define CHECK_ONLY_DEPS
	$(REPO_ROOT)/module/cli
endef

ifeq "$(MAKO_STAGE)" "main"

DEFAULT_TARGETS := genfiles/build_done
DEFAULT_PREREQS := Dockerfile* \
		$(shell find ../module -name genfiles -prune -o -print) | genfiles
genfiles/build_done:
	./build_all
	touch $@

genfiles:
	mkdir -p $@

endif

include $(MAKO_ROOT)/component.mk

