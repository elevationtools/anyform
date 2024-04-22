
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

### Interface / Breaking Changes

- Allow configuring a different command than just `./ctl up|down` for stages.

- Consider having stage specific configs, and then check if the config has
  actually changed for a stage, instead of just checking the timestamps like GNU
  make.

- Consider removing the noise `_DIR` suffix from environment variables.

- Better name than "config" for the application specific configuration passed to
  stages within the cell's `anyform.jsonnet` file.

### Additions / Non-breaking Changes

#### Bug/Quirk Fixes

- signal handling for subprocesses

#### Behavior Improvements

- Exponential backoff retries for stages, configurable in anyform.jsonnet.

- If a stage hasn't even been brought up at all, then assume it's down.  Perhaps
  guard with a flag like "--assume-down".

- Support "one shot" stages that shouldn't be run again if they've already
  succeeded.
  - Support manually forcing them to run again.
  - Perhaps add a new state, or new field in the state file, which if it's set
    then don't run the stage.

- Restamp on `down` in some scenarios: e.g. impl has changed.

- After stamping, remove any files that were there from a previous stamp but not
  this one.  Alternatively, just fully delete the stamp dir and recreate it.

- Option to select and filter stages.
  - Idea: `anyform mark (up|down) STAGE_NAME` allows forceably skipping a step.
    Creates the state file.

- Feature: `anyform down --skip-failures` to just optimistically try to bring
  down all stages, even if there was a failure while trying to bring down a
  child stage.

- Think about how to handle a stage calling a different orchestrator to support
  nested/modular anyforms.

- Enable controlling parallelism, at a minimum full on and full off.

- Add configo as a stamper.

- "make -C" like feature to chdir without having to chdir.

#### Printing/Logging Improvements

- Tee Stage.stdout/stderr to stage log file.

- Consider outputting error messages outside the context of a stage to somewhere
  other than stdout/stderr.

- `anyform env $stage_name` prints environment used for running stages.

#### Internal Code Clean-up

- Reconsider split of Config, Spec, and InnerConfig.
  - Especially that OrchestratorSpecFile is in AnyformConfig.

- Better error when anyform.jsonnet file not found.

#### Documentation

- Provide example of breaking a terraform lock:
  https://developer.hashicorp.com/terraform/cli/commands/force-unlock

- Add https://atmos.tools/ to alternatives comparison.

