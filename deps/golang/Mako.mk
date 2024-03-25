
include $(MAKO_ROOT)/dep.mk

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	curl -Lo genfiles/download.tar.gz \
		'https://go.dev/dl/go1.22.1.linux-amd64.tar.gz'
	cd genfiles && \
		tar xzvf download.tar.gz

smoketest:
	go version | grep 1.22

endif

