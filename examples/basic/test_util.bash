
function runtest() {
  test_name="${1:?}"
  echo -n "$test_name: "
  if test_output="$("$test_name")"; then
    echo "OK"
  else
    echo "FAILURE"
    echo "$test_output"
    return 1
  fi
}

# The timestamp changes so we just strip it off for testing.
function strip_timestamp() {
  sed -E 's/^[^ ]+ //g'
}

function expect_stripped_eq_stdin() {
  local actual="${1:?}"
  expect_eq_stdin "$(echo "$actual" | strip_timestamp 2>&1)"
}

function expect_eq_stdin() {
  local actual="${1:?}"
  local expected
  expected="$(< /dev/stdin)"
  expect_eq "$actual" "$expected"
}

function expect_eq() {
  local actual="${1:?}"
  local expected="${2:?}"
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
  actual="${1:?}"
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

