
# Anyform

- [Table of Contents](/README.md)

## Overview

### Key design goals

#### Primary goals

- Infrastructure turn-up/turn-down scripts can be expressed as a Directed
  Acyclic Graph (DAG) of one or more "stages", with dependents able to use the
  output of dependencies.

- Bringing up a new environment (prod/dev/staging/etc) or a new locations
  (us-west-1/eu-north-1/etc) only requires creating a single config file with a
  few lines in it.
  
- DAG stages can be written in any langugae (Terraform, python, bash calling
  CLIs like kubectl, golang, etc). A single DAG will often use multiple
  languages across stages.

#### Secondary goals

- Configuration can be written in a general purpose config language that can
  produce JSON. Currently, Jsonnet is used.

- Stage definitions can use a general purpose templating engine to fill missing
  pieces of tools like Terraform, Helm, etc (the gaps are identified
  below), rather than requiring specific solutions for each tool. Currently,
  gomplate is used.

#### Tertiary goals

- Allow plugging in different config and templating engines instead of Jsonnet
  and gomplate.  This is already easy within the golang code, but cannot yet be
  done post-compilation of `anyform`.
  

## Common Tooling Problems Addressed

This section describes some of the problems with common tools that motivated
Anyform.

### Terraform Problems

- Inability to parameterize `backend` configuration.
- Lack of solution when multi-staging is required. (and [Terraform admits it's
  required in many common
  cases](https://github.com/hashicorp/terraform/issues/27785#issuecomment-780017326)
  and [AWS Terraform Blueprints
  mentions](https://github.com/aws-ia/terraform-aws-eks-blueprints?tab=readme-ov-file#terraform-caveats)
  this is problematic for them too).
- Extreme boilerplate when trying to use a module as a library for multiple
  different deployments/environments/etc.  In particular, there is much
  redundancy of the `terraform` and `provider` blocks.

These are the same problems that motivated the existence of products like
Terragrunt and Terraspace.

### Jsonnet Problems

- Lack of computed imports.
- No convenient access to environment variables. (`-V` + `std.extVar` is tedious
  and can be problematic when jsonnet is invoked for you and you don't have the
  ability to provide `-V` to it).

### Helm Problems

- Lack of ability to parameterize `Values.yaml` with go templating.
- Lack of transitive dependency fetching preventing code reuse through layering.
  (though this arguable isn't solved by Anyform, aside from allowing you to
  avoid Helm as much a possible and instead utilize tools like Jsonnet).


## Comparing to Alternatives

### Terragrunt and Terraspace

The most obvious alternatives are Terragrunt and Terraspace. They are motivated
by the same problems in Terraform discussed above.  Compared to these,
Anyform has the following advantages:

- Anyform is not specific to Terraform, instead allowing stages to be written
  in a mix of Terraform, bash calling CLIs, etc.  The Terra\* tools only support
  Terraform stages.

- Anyform requires learning much less single-purpose and framework-specific
  knowledge.  Instead, it makes use of generally useful tools, like gomplate
  and Jsonnet.

- The ability to use Jsonnet (or other config language) provides a few benefits:
  - Allows using a more powerful and productive configuration language than pure
    Terraform, allowing writing much more declarative configurations than pure
    Terraform.  Terraform instead just reads the highly declarative
    configuration and converts it into API calls to the cloud to create infra.
  - Allows easily and efficienctly using your configuration outside of Terraform
    when Terraform isn't a good tool for the job. (No need for the slow "init"
    and managing providers and state).

An additional advantage over Terraspace is that it has no dependency on Ruby.

### Tonka

Tonka and Anyform aren't really trying to address the same problem and so
could be reasonably used in complementary roles. Tonka is focused on Jsonnet for
kubernetes application deployments.  Anyform is more intended to bring up the
infrastructure (cloud networking, managed kubernetes service, etc) that the
application deployments would then use.
