
include $(MAKO_ROOT)/util.mk

define CHECK_ONLY_DEPS
	$(REPO_ROOT)/module/cli
endef

.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Targets: local github"

ifeq "$(MAKO_STAGE)" "main"

$(call mako_define_target, local, genfiles/local_done, \
		build_all_platforms_via_docker Dockerfile* | genfiles)
genfiles/local_done:
	CONTAINER_REGISTRY=none ./build_all_platforms_via_docker
	touch $@

# Build using github ghcr.io for layer caching. The final image is pushed as
# well, but not used.
$(call mako_define_target, github, genfiles/github_done, \
		build_all_platforms_via_docker Dockerfile* | genfiles)
genfiles/github_done:
	CONTAINER_REGISTRY=ghcr.io/elevationtools ./build_all_platforms_via_docker
	touch $@

genfiles:
	mkdir -p $@

endif

include $(MAKO_ROOT)/base.mk

