
function activate_anyform_repo_bash() {
  export REPO_ROOT="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
  . "$REPO_ROOT/activate.sh"

  # Avoids adding a duplicate.
  export PS1="(anyform) ${PS1#(anyform) }"

  # Fix any duplicates from running activate.bash twice.
  export PATH="$(remove_dups_from_colon_sep_string "$PATH")"
}

function remove_dups_from_colon_sep_string() {
  local list="${1:?}"
  (
    export IFS=:
    declare -A seen
    declare -a out
    for x in $list; do
      if (( seen[$x] != 1 )); then
        out+=("$x")
      fi
      seen[$x]=1
    done
    echo "${out[*]}"
  )
}

activate_anyform_repo_bash

