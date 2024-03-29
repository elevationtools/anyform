
## Current State

just finished: stage "up" checks if `ANYFORM_CONFIG_JSON_FILE` is newer than the state
file and that the last command was actually "up".  state file is a JSON file
called `genfiles/stage_name/state`.

need to:
- not always update the `ANYFORM_CONFIG_JSON_FILE` by checking the jsonnet-deps
  (or byte fingerprint) of the `AnyConfig.OrcherstratorSpecFile.`.
- check the `ANYFORM_IMPL_DIR/stage_name/*` modified times for the stage and all
  its parents.

## Todos

- Avoid redoing work that doesn't need to be done, GNU Make style or via
  checksums.
  - Important: DAG children impl or output changing shouldn't cause parents to
    rerun.

- Enable controlling parallelism, at a minimum full on and full off.

- Reconsider split of Config, Spec, and InnerConfig.
  - Especially that OrchestratorSpecFile is in AnyformConfig.

- Better logging format.
  - Maybe a more natural way to call commands and have them echoed with output,
    similar to GNU make.

- Consider creating the output dir and stage output dirs automatically so `ctl`
  scripts don't have to.

- allow passing options to gomplate, per directory or even per file, to change
  --left-delim --right-delim
  - Or just move away from gomplate to configo and support that with configo
    directly

- Think about how to handle a stage calling a different orchestrator to support
  nested/modular anyforms.

- Allow configuring a different command than just `./ctl up|down` for stages.

