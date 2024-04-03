
local anyform = import 'anyform/anyform.libsonnet';
local stage = anyform.stage;

anyform.dag(std.thisFile) {
  stages: {
    shop: stage([]),
    local cooking_tmpl = stage(['shop']),
    cook_curry: cooking_tmpl,
    cook_rice: cooking_tmpl,
    eat: stage(['cook_curry', 'cook_rice']),
  },

  config: {
    night: error 'required',

    nut_allergy: false,
    color: 'green',
  },
}
