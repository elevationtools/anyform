
## Todos

- Use tests in CI/CD

- Improve error message for stage ctl script error and stamping error (e.g.
  `EATER_NAME` wasn't set, {{ Env.EATER_NAME }} barfed, and the error message is
  escaped and almost unreadable, similarly for a normal script error).

- Integration tests.

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

