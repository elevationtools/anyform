
include $(MAKO_ROOT)/dep.mk

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	curl -Lo genfiles/download \
		'https://github.com/jsonnet-bundler/jsonnet-bundler/releases/download/v0.5.1/jb-linux-amd64'
	chmod +x genfiles/download

smoketest:
	jsonnet-bundler --version 2>&1 | grep v0.5.1

endif

