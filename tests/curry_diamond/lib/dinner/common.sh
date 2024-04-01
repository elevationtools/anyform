
set -euo pipefail

function common_main() {
  case "$1" in
    up|down) "run_$1" ;;
    *) echo "Unknown command: $1" >&2; exit 1 ;;
  esac
}

