
function runtest() {
  test_name="${1:?}"
  echo -n "$test_name: "
  if "$test_name" &> /dev/null; then
    echo "OK"
  else
    echo "FAILURE ($?)"
    ( set -euo pipefail; "$test_name" ) || true
    return 1
  fi
}

function expect_eq_stdin() {
  local actual="$1"
  local expected
  expected="$(< /dev/stdin)"
  expect_eq "$actual" "$expected"
}

function expect_eq() {
  local actual="$1"
  local expected="$2"
  if [[ "$actual" != "$expected" ]]; then
    echo "----------"
    echo "Actual:"
    echo "$actual" | sed 's/^/  /g'
    echo "----------"
    echo "Expected:"
    echo "$expected" | sed 's/^/  /g'
    echo "----------"

    numNewLines="$(echo -n "$actual$expected" | wc -l)"
    if (( numNewLines > 0 )) &> /dev/null; then
      echo "Diff:"
      diff -u <(echo "$actual") <(echo "$expected") || true
      echo "----------"
    fi
    return 1
  fi
}

function expect_json_eq_stdin() {
  actual="$1"
  expected="$(cat)"

  # normalize
  actual="$(echo "$actual" | jq --sort-keys 2>&1)"
  expected="$(echo "$expected" | jq --sort-keys 2>&1)"

  if [[ "$actual" != "$expected" ]]; then
    echo "----------"
    echo "Actual:"
    echo "$actual"
    echo "----------"
    echo "Expected:"
    echo "$expected"
    echo "----------"
    echo "Diff:"
    diff -u <(echo "$actual") --label actual \
            <(echo "$expected") --label expected || true
    echo "----------"
    return 1
  fi
}

