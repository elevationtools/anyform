
include $(MAKO_ROOT)/dep.mk

define DEPS
	$(REPO_ROOT)/deps/golang
endef

ifeq "$(MAKO_STAGE)" "main"

Linux-x86_64:
	go install github.com/google/go-jsonnet/cmd/jsonnet@v0.20.0
	go install github.com/google/go-jsonnet/cmd/jsonnet-lint@v0.20.0

smoketest:
	jsonnet -version | grep 0.20.0
	jsonnet-lint -version | grep 0.20.0

endif

