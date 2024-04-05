
# Anyform

- [Table of Contents](/README.md)


## Learning

### Examples

#### Hello world
TODO

#### Working example: [`//tests/curry_diamond`](/tests/curry_diamond)
This is a working example showing what Anyform can do and is also an integration
test so you know it actually works!  Read the README, run the integration test,
explore the files - especially the generated files.

#### Narrated Example: [`/docs/narrated_example.md`](/docs/narrated_example.md)
This is a "words-only" example, with code snippets but no runnable code.

#### AWS EKS example
TODO


### Getting Started / Walkthrough
TODO


### Core Concepts / Terminology

#### Stage
A program/script that brings up/down some infrastructure.

#### DAG
A directed acyclic graph of one or more stages.

#### Deployment / Cell
A deployment is an instance of some bundle of infrastructure.  Usually multiple
deployments of the same infrastructure are created, for example
`production/us-west-1` and `staging/eu-north-1` are deployments.  In this case
the deployment has an environment (prod/staging/etc) and a location
(us-west-1/eu-north-1/etc).  This pattern is assumed throughout the
documentation, but no assumption about this is made in the Anyform codebase.
Your deployments can exist along any dimensions.

#### Deployment Config (`./anyform.jsonnet`)
The configuration file for a single deployment.  This must be named
`anyform.jsonnet`.  It describes:
- The stages that in the DAG.
- How the stages depend on each other.
- Where the stage implementations are located.
- The configuration to be passed to all stages (all stages share the same
  config, but the config can be organized into stage specific sections if
  desired).

To bring up a new deployment, one should only need to create a new one of these
files, perhaps copying an existing one then tweak the settings in it for the
given deployment, and then run `anyform up`.

#### Deployment Dir
The directory containing the deployment config (`anyform.jsonnet`).  Upon
running `anyform up`, this directory will also contain generated files
including:

##### - Output Dir (`./output/`)
Stage outputs that should be checked into version control.  Stage's are
recommended to output into `output/STAGE_NAME/` but this is not a requirement
and there are legitimate reasons to do differently.  Putting output outside the
output dir is highly discouraged.

##### - Genfiles Dir (`./genfiles`)
Outputs and temporary files that should NOT be checked into version control.

#### DAG Template Directory
The directory containing stage template directories. Stage template directories
must live directly under the DAG template directory.

Usually, the DAG template directory will also contain a Jsonnet library to be
imported by the deployment config.  By convention, this is called
`anyform.libsonnet`, but can be called anything.

#### Stage Template Directories
Where the code for a stage actually lives.  Every file is considered a gomplate
template, unless otherwise configured. `.gomplate.yaml` and `.gomplateignore`
can live in the directory.  This can be important for things like binaries or
Helm charts, since they also use golang text templating.

#### Stage Stamp Directories
Each stage is stamped via gomplate to its own stage stamp directory. These live
in the deployment directory under genfiles `./genfiles/STAGE_NAME/stamp`. In
some cases it's possible to use the stamp directory directly. For example, with
Terraform this is the case, the `terraform` CLI can be used in the stage stamp
directory.  Not all stages support this though, as some will expect to be run
with environment variables set by Anyform.

To make your stamped directory usable directly, like the terraform stage, be
sure to only use anyform specific environment variables during stamping, and not
while running the stage `ctl`.

Note that stages aren't stamped until AFTER all dependency stages have run
successfully.  This means dependency stage outputs (both in `./genfiles/` an
`./output/`) are available to stages during stamp time (and also during run
time).

#### Stage `ctl`
An executable that is the entrypoint to the stage, written by the Anyform user,
called by Anyform after the stage has been stamped.  Called with the stage
environment variables set (discussed below).

Requirements:

`ctl` files must support being called 2 commands:
```
./ctl up
./ctl down
```

If environment variable `INTERACTIVE=false` then the stage MUST NOT prompt the
user.  Any other value than exactly `false` should be considered true.

`ctl` can be written in any language.


#### Stage Environment Variables

The following environment variables are set during stage stamp time, and `ctl`
run time.
```
ANYFORM_STAGE_NAME
ANYFORM_STAGE_STAMP_DIR
ANYFORM_CONFIG_JSON_FILE
ANYFORM_GENFILES
ANYFORM_IMPL_DIR
ANYFORM_OUTPUT_DIR
```

> To access during stamp time, use gomplate standard approach: `{{ .Env.FOO }}`


### Detailed Specification/Reference
TODO

