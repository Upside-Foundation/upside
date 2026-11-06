local config = import 'default.jsonnet';

config {
  'upside_61872-1'+: {
    validators: [validator {
      client_config: {
        'broadcast-mode': 'sync',
      },
    } for validator in super.validators],
  },
}
