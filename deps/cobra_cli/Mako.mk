
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/golang
endef

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	go install github.com/spf13/cobra-cli@latest

smoketest:
	cobra-cli --help | grep 'Cobra is a CLI' > /dev/null

endif

include $(MAKO_ROOT)/dep.mk

