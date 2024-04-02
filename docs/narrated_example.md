
# Anyform

- [Table of Contents](/README.md)
- [Learning](/docs/learning.md)

## Narrated Example

This gives an example of the recommended usage pattern of Anyform.

> You may prefer the working example instead of this narrated example. It's
> also somewhat narrated, but also contains runnable code (it's even used as
> an integration test). See
> [`//tests/curry_diamond/README.md`](/tests/curry_diamond/README.md)

Imagine a user needs to bring up infrastructure in a project cleverly called "My
Infra".  They need to bring this up in 2 separate regions, `eu-north-1` and
`us-east-1` in production (for simplicity there's no staging/dev/etc in this
example because different environments work the same as different locations).

Assume the user has a repo that looks like the following to model this:
```
$MY_REPO
├── cells
│   └── prod
│       ├── eu-north-1
│       │      └── my_infra
│       └── us-east-1
│              └── my_infra
└── infra_lib
    ├── my_infra
    └── other_unreleated_project
 
``` 

- `infra_lib` contains libraries that can be called to bring up infrastructure.
- `cells` contains directories for each tangible instantiation of
infrastructure.  (e.g. `prod` deployed to `eu-north-1`).  Each of the leaf cell
directories calls into `infra_lib` and passes it the specific configuration
needed for that cell (other equivalent terms: cluster, deployment, etc).

### Operational usage after setup

First, let's start with how the operational usage would look like if Anyform
has already been setup for "My Infra".  The operator wouldn't care much about
the `infra_lib` directory.  Instead, they'd be focused on the `cells` directory,
which would then look something like this:

```
$MY_REPO/cells/prod
├── eu-north-1
│   ├── cell_config.jsonnet
│   └── my_infra
│       └── anyform.jsonnet
└── us-east-1
    ├── cell_config.jsonnet
    └── my_infra
        └── anyform.jsonnet
```

To bring up and down the "My Infra" project in `eu-north-1`, the user would then:
```
cd $MY_REPO/cells/prod/eu-north-1/my_infra

anyform up
# Now the infra is up.

anyform down
# Now the infra is down again.
```

And that's it! Upon successfully running `anyform up` the directory would look
like the following:
```
$MY_REPO/cells/prod/eu-north-1/my_infra
├── anyform.jsonnet
├── output/...
└── genfiles/...
```

`output/...` will contain output, if any, usually as one or more JSON files.
These SHOULD be checked into version control.

`genfiles/...` contains other output that SHOULD NOT be checked into version
control (consider putting `**/genfiles` in `.gitignore`).  This generally
contains temporary files, caching, etc.

## Creating this Anyform setup

To make the above work, the following must first be created.  This only needs to
be done once and then be used across multiple cells:
```
$MY_REPO/infra_lib/my_infra (a.k.a impl_dir)
├── anyform.libsonnet
├── stage_one
│   ├── ctl
│   └── some_terraform_code.tf
└── another_stage
    ├── ctl
    ├── some_kube_manifest.yaml
    └── another_kube_manifest.jsonnet
```

`stage_one` and `stage_two` are stage template directories containing the code
which actually brings up infrastructure. It can be Terraform, bash, a golang
binary, anything really.

`$STAGE_NAME/ctl` is an executable file which anyform calls to bring the stage's
infrastructure up or down.
- > See [`/common_stage_utilities.md`](/common_stage_utilities.md) to avoid
  boilerplate for common stage types.

`anyform.libsonnet` would look something like:
```jsonnet
(import 'anyform/jsonnet_lib/orchestrator.libsonnet')(std.thisFile) {
  stages: {
    stage_one: {},
    another_stage: {
      depends_on: ['stage_one'],
    },
  },

  // This config is made available to stages via $ANYFORM_CONFIG_JSON_FILE
  config: {
    // Cells must override this with the cell's `cell_config.jsonnet`.
    cell: error 'required',

    foo: {
      bar: 'a_default_value',
      baz: error 'a required value that must be specified per cell',
    },
  },
}
```

`$MY_REPO/cells/prod/eu-north-1/my_infra/anyform.jsonnet` would then
look something like...
```jsonnet
(import 'infra_lib/my_infra/anyform.libsonnet') {
  config: {
    cell: import '../cell_config.jsonnet',
    foo: { baz: 'cheese' },
  },
}
```

## genfiles details

The `./genfiles/...` directory looks something like:

```
$MY_REPO/cells/prod/eu-north-1/my_infra/genfiles/
├── stage_one
│   ├── stamp
│   │   ├── ctl
│   │   └── ...etc...
│   ├── state
│   └── logs
│       └── stdout_stderr
│
└── another_stage
    └── ...same as stage_one above...
```

Each stage gets its own directory in genfiles.  The "stamp" directory is where
the template is stamped into and where `ctl` is run from.  For Terraform stages,
an operator user can then navigate to this directory and run `terraform` CLI
commands directly to troubleshoot, repair, do advanced thing like migrate state,
shoot yourself in the foot, etc.

The `state` tracks whether it's most recently been brought up or down, and the
modification time tracks when this happened, which is used to avoid repeating
the operation if nothing has changed since the last time run.

`logs/stdout_stderr` contains the stdout and stderr (merge together) of the
latest run.
