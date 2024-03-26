local util = import './util.libsonnet';

local caseTmpl = function(input, expected) {
  input: input,
  expected: expected,
};

local cases = [
  caseTmpl('/', '/'),
  caseTmpl('.', '.'),
  caseTmpl('./', '.'),
  caseTmpl('..', '.'),
  caseTmpl('../', '.'),
  caseTmpl('/foo', '/'),
  caseTmpl('/foo/bar', '/foo'),
  caseTmpl('foo', '.'),
  caseTmpl('foo/bar', 'foo'),
];

[
  case {
    actual: util.dirname(self.input),
    result: if self.actual == self.expected then 'success' else 'FAILURE'
  }
  for case in cases
]

