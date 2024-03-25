
# Requires REPO_ROOT is already set.

activate_anyform_repo_sh() {
  if test -z "$REPO_ROOT"; then
    echo "REPO_ROOT must already be set"
    exit 1
  fi
  export PATH="$REPO_ROOT/deps/bin:$PATH"
  export GOPATH="$REPO_ROOT/deps/gopath"
  export MAKO_ROOT="$REPO_ROOT/deps/mako/lib"
}

activate_anyform_repo_sh

