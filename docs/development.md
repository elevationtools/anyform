
# Anyform

- [Table of Contents](/README.md)

## Development

### Run tests

```bash
make test
```

### Create a release

- Create a tag (can also be done via the GitHub web UI):
  ```bash
  git tag v1.2.3
  git push origin v1.2.3-rc0
  ```
  This triggers the GitHub Actions workflow defined in
  `.github/workflows/release.yaml`

- On the GitHub web UI, check that workflow succeeded and the draft release was
  created.
- Publish it as a pre-release.
- When there's an issue with the candidate, start over using suffix `-rcN+1`.
- Once a candidate has proven itself stable, create another tag `v1.2.3`
  pointing to the same commit (without the `-rcN` suffix), publish it as a full
  release.
- Once a candidate is no longer in use, delete it to clean up.

