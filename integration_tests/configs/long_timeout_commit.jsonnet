local default = import 'default.jsonnet';

default {
  'upside_61872-1'+: {
    config+: {
      consensus+: {
        timeout_commit: '15s',
      },
    },
    'app-config'+: {
      'blocked-addresses': ['up16z0herz998946wr659lr84c8c556da55wnla3a'],
    },
  },
}
