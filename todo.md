
## Current State

just finished implementing skipping of already done operations.

## Todos

- Automatic cross platform release creation using github actions.

- Look into using gomplate as library rather than the CLI binary.

- allow passing options to gomplate, per directory or even per file, to change
  --left-delim --right-delim
  - Or just move away from gomplate to configo and support that with configo
    directly

- Enable controlling parallelism, at a minimum full on and full off.

- Reconsider split of Config, Spec, and InnerConfig.
  - Especially that OrchestratorSpecFile is in AnyformConfig.

- Better logging format.
  - Maybe a more natural way to call commands and have them echoed with output,
    similar to GNU make.

- Consider creating the output dir and stage output dirs automatically so `ctl`
  scripts don't have to.

- Think about how to handle a stage calling a different orchestrator to support
  nested/modular anyforms.

- Allow configuring a different command than just `./ctl up|down` for stages.

- signal handling for subprocesses

- Handle error output from StageStampers better (stdout/stderr for CLI based
  ones)
