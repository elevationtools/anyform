
(import 'jsonnet_lib/orchestrator.libsonnet')(std.thisFile) {
  stages: {
    shop: {},
    local cooking_tmpl = { depends_on: ['shop'] },
    cook_curry: cooking_tmpl,
    cook_rice: cooking_tmpl,
    eat: { depends_on: ['cook_curry', 'cook_rice'] },
  },

  config: {
    night: error 'required',

    nut_allergy: false,
    color: 'green',
  },
}
