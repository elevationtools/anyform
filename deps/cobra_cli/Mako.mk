

define DEPS
	$(REPO_ROOT)/deps/golang
endef

include $(MAKO_ROOT)/dep.mk

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	go install github.com/spf13/cobra-cli@latest

smoketest:
	cobra-cli --help | grep 'Cobra is a CLI' > /dev/null

endif

