
local dirname = import 'elevation/dirname.libsonnet';

{
  // See //tests/curry_diamond/lib/dinner/anyform.libsonnet for an example showing
  // how this function is intended to be used.
  //
  // `pathToFileInImplDir` is the path to any file within the "dag template dir"
  // (the directory which contains the "stage template directories").  This is
  // intended to be called by a jsonnet file within the "dag template dir"` so
  // that it can just do:
  //  local anyform = import 'anyform/anyform.libsonnet';
  //  anyform.dag(std.thisFile) { ... }
  dag(pathToFileInImplDir): {
    impl_dir: dirname(pathToFileInImplDir),

    // A map of stages, each preferably created by $.stage(...).
    stages: {},

    // An object of any shape.  Implementations should specify the shape that is
    // required. This is converted into genfiles/config.json for use in stages.
    config: {},
  },

  // `depends_on` is a list of stage names.
  stage(depends_on=[]): {
    depends_on: depends_on,
  },
}
