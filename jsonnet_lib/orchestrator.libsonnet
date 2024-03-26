
local util = import 'jsonnet_lib/util.libsonnet';

function(orchestratorFilePath) {
  impl_dir: util.dirname(orchestratorFilePath),
}

