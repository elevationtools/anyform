
include $(MAKO_ROOT)/dep.mk

define DEPS
	$(REPO_ROOT)/deps/golang
endef

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	curl -Lo genfiles/download \
		'https://github.com/hairyhenderson/gomplate/releases/download/v4.0.0-pre-2/gomplate_linux-amd64'
	chmod +x genfiles/download

smoketest:
	gomplate --version | grep 4.0.0-pre-2

endif

