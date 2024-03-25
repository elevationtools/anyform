
export REPO_ROOT := $(CURDIR)

build:
	. ./activate.sh && mako -C cli
	
