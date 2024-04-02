
# Anyform

- [Table of Contents](/README.md)

## Development

### Run tests

```bash
make test
```

### Create a release

```bash
git tag v1.2.3
git push origin v1.2.3
```
(This triggers the GitHub Actions workflow defined in
`.github/workflows/release.yaml`)

Then, on the GitHub web UI, check that workflow succeeded and the release was
created.

> TODO: Consider changing this procedure.  For example, manually create a draft
> release via the GitHub web UI, then have it run tests and on success
> auto-populate with artifacts, then manually flip to non-draft.

