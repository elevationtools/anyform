
The jsonnet library located in this repo at
[`//jsonnet_lib/anyform`](/jsonnet_lib/anyform) was added via:
```bash
jsonnet-bundler install \
  https://github.com/elevationtools/anyform/jsonnet_lib/anyform@$VERSION
```
This may seem odd since it's local, but it was done so that tests in
`//tests` uses the same strategy that users of anyform must do.

