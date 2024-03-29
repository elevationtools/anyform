
## Current State

`mako -C examples/basic up` stamps all stages but the stages aren't actually
run.

## Todos

- Better logging format.
  - Maybe a more natural way to call commands and have them echoed with output,
    similar to GNU make.
- Enable log level setting via flag.
- Store command stdout and stderr somewhere better.
- Mirror file mode from template dir to stamp dir (currently sets 0750).
- Don't overwrite `CONFIG_JSON_FILE` if already up to date.
  - Need to think about how to avoid rerunning if everything's up to date.  A
    large amount of GNU make ends up being reimplemented...hence the beauty of
    the v0 prototype that actually used GNU make.
- consider out creating the output dir and stage output dirs

### Longer term

- allow passing options to gomplate, per directory or even per file, to change
  --left-delim --right-delim
  - Or just move away from gomplate to configo and support that with configo
    directly

- maybe: don't update stamped file timestamp if it didn't change, similarly for
  output.

