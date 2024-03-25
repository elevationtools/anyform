
function activate_anyform_repo_bash() {
  export REPO_ROOT="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
  export PATH="$REPO_ROOT/deps/bin:$PATH"
  export GOPATH="$REPO_ROOT/deps/gopath"
  export MAKO_ROOT="$REPO_ROOT/deps/mako/lib"
}

activate_anyform_repo_bash

