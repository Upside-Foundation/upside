local config = import 'default.jsonnet';

config {
  'upside_61872-1'+: {
    validators: super.validators + [{
      name: 'rollback-test-memiavl',
      'app-config': {
        memiavl: {
          enable: true,
        },
      },
    }, {
      name: 'rollback-test-iavl',
      'app-config': {
        memiavl: {
          enable: false,
        },
      },
    }],
  },
}
