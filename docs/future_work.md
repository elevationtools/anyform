
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

- Restamp on `down` in some scenarios: e.g. impl has changed.

- `anyform env` to print environment used for running stages.

- Better name than "config" for the application specific configuration passed to
  stages within the cell's `anyform.jsonnet` file.

- After stamping, remove any files that were there from a previous stamp but not
  this one.  Alternatively, just fully delete the stamp dir and recreate it.

- Better error when anyform.jsonnet file not found.

- "make -C" like feature to chdir without having to chdir.

- Feature: `anyform markdown STAGE_NAME` allows forceably skipping a step.
  Creates the state file.

- Provide example of breaking a terraform lock:
  https://developer.hashicorp.com/terraform/cli/commands/force-unlock

- Feature: `anyform down --skip-failures` to just optimistically try to bring
  down all stages, even if there was a failure while trying to bring down a
  child stage.

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
  - Also delete the standard output dir during "down".

- Option to select and filter stages.

