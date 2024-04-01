
# Anyform

Infrastructure-as-code tool allowing defining turn-up/down scripts as a directed
acyclic graph of stages that are potentially heterogeneous in their
implementation tools (Terraform, direct kubectl calls, etc).


## Installation via binaries

Download a binary release and put it in your `PATH`:
- https://github.com/elevationtools/anyform/releases

Also install `gomplate` into your `PATH`:
- https://github.com/hairyhenderson/gomplate/releases

Optionally but fairly importantly, install `jsonnet` into your `PATH` as well.
While not required, if you're using `anyform` then you'll likely need to use
`jsonnet` directly too:
- https://github.com/google/go-jsonnet/releases


## Building it yourself

Requirements:
- git
- GNU make
- curl
- bash
- standard unix CLI tools like: cp, rm, find, etc.

Checkout the desired version:
```bash
git clone https://github.com/elevationtools/anyform
git checkout v1.2.3  # optional, selects a version.
```

Build locally without docker:
```bash
make build
./module/cli/genfiles/bin/anyform
```

Alternatively, build for all supported platforms using docker:
```bash
git clone https://github.com/elevationtools/anyform
make docker_build
ls -l build/genfiles/bin
```

These instructions have been verified on:
- Ubuntu 22.04

