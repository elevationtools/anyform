
local dirname = import 'elevation/dirname.libsonnet';

function(orchestratorFilePath) {
  impl_dir: dirname(orchestratorFilePath),
}

