
FROM ubuntu:22.04

ARG ANYFORM_VERSION=0.0.0-unsetdocker

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get -y update
RUN apt-get -y install make curl

ENV REPO_ROOT=/anyform
ENV PATH="$REPO_ROOT/module/cli/genfiles/bin:$REPO_ROOT/deps/bin:$PATH"
ENV GOROOT="$REPO_ROOT/deps/golang/genfiles/go"
ENV GOPATH="$REPO_ROOT/deps/gopath"
ENV MAKO_ROOT="$REPO_ROOT/deps/mako/lib"
ENV JSONNET_PATH="$REPO_ROOT"

WORKDIR "$REPO_ROOT"

# Long running deps first to utilize the layer cache.
COPY /deps/ ./deps/
COPY /Makefile /activate.sh ./
RUN mako -C deps/golang && \
    mako -C deps/gomplate && \
    mako -C deps/jsonnet && \
    mako -C deps/jsonnet-bundler

COPY /module ./module
COPY /tests ./tests
RUN make test
RUN mako -C ./module/cli all_platforms

