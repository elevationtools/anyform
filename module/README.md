
# Anyform

## `module` directory

Contains the single go module for anyform. Subdirs:

- `lib`
  - The primary location of Anyform code.

- `cli`
  - Utilizes the cobra library to expose `lib` as a CLI.
  - Very little functionality should be implemented here so that `lib` is easily
    usable on its own.

- `common/util`
  - Utilities common to both `lib` and `cli`.

