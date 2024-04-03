
include $(MAKO_ROOT)/util.mk

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	curl -Lo genfiles/download \
		'https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64'
	chmod +x genfiles/download

smoketest:
	jq --version | grep 1.7.1

endif

include $(MAKO_ROOT)/dep.mk

