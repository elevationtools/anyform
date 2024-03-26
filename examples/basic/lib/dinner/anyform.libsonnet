
(import 'jsonnet_lib/orchestrator.libsonnet')(std.thisFile) {
  stages: {
    shop: {},
    local cooking_tmpl = { dependsOn: ['shop'] },
    cook_curry: cooking_tmpl,
    cook_rice: cooking_tmpl,
    eat: { dependsOn: ['cook_curry', 'cook_rice'] },
  },

  config: {
    night: error 'required',

    nut_allergy: false,
    color: 'green',
  },
}
