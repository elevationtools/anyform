
# Anyform

- [Table of Contents](/README.md)

## Future Work

### Issues with `gomplate`

Solving the following issues with `gomplate` would be helpful for Anyform

- Calling templates in another file is awkward/boilerplate heavy.
  ```golang
  {{ $_ := file.Read "otherfile.gmp" | tpl }}
  {{ tmpl.Exec "otherfiletmpl" (merge (dict "foo" "bar") .) }}
  ```

- `--input-dir` and `--output-dir` has problematic semantics.  One of these two
  options would solve it:
    - Option 1) Remove files in `--output-dir` that aren't in `--input-dir` to clean up
      previous runs automatically without having to `rm -rf` the whole output dir.
    - Option 2) Make the timestamps of the output files match the timestamps of
      the input files.  This would allow `rm -rf` to work while not confusing
      `make`'s timestamp checking.


## Unorganized Todos

- Support `--version` or `anyform version`

- Reconsider split of Config, Spec, and InnerConfig.
  - Especially that OrchestratorSpecFile is in AnyformConfig.

- Add configo as a stamper.

- Allow configuring a different command than just `./ctl up|down` for stages.

- Better logging/error handling.
  - Maybe a more natural way to call commands and have them echoed with output,
    similar to GNU make.
  - Handle error output from StageStampers better (stdout/stderr for CLI based
    ones)

- Think about how to handle a stage calling a different orchestrator to support
  nested/modular anyforms.

- Enable controlling parallelism, at a minimum full on and full off.

- signal handling for subprocesses

- Consider creating the output dir and stage output dirs automatically so `ctl`
  scripts don't have to.

