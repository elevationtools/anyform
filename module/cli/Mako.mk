
include $(MAKO_ROOT)/util.mk

define DEPS
	$(REPO_ROOT)/deps/golang
endef

ifeq "$(MAKO_STAGE)" "main"

DEFAULT_TARGETS := genfiles/bin/anyform
DEFAULT_PREREQS := $(shell find . ../lib ../common -name genfiles -prune -o -name "*.go") \
	| genfiles/bin
genfiles/bin/anyform:
	go build -o $@ .

# 1: GOOS value
# 2: GOARCh value
define os_arch_target_impl

genfiles/bin/anyform-$(1)-$(2): genfiles/bin/anyform | genfiles/bin
	GOOS=$(1) GOARCH=$(2) go build -o $$@ .

all_platforms: genfiles/bin/anyform-$(1)-$(2)
endef
os_arch_target = $(eval $(call os_arch_target_impl,$(1),$(2)))

$(call os_arch_target,linux,amd64)
$(call os_arch_target,linux,arm64)
$(call os_arch_target,darwin,amd64)
$(call os_arch_target,darwin,arm64)
$(call os_arch_target,windows,amd64)
$(call os_arch_target,windows,arm64)

genfiles/bin:
	mkdir -p $@

endif

include $(MAKO_ROOT)/component.mk

