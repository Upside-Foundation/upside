local config = import 'default.jsonnet';

config {
  'upside_61872-1'+: {
    'start-flags': '--trace --inv-check-period 5',
    'app-config'+: {
      'minimum-gas-prices':: super['minimum-gas-prices'],
      'json-rpc'+: {
        api:: super['api'],
      },
    },
    accounts: [{
      name: 'community',
      coins: '10000000000000000000000taup',
      mnemonic: '${COMMUNITY_MNEMONIC}',
    }],
    genesis+: {
      app_state+: {
        upside: {
          params: {
            upside_admin: 'up12luku6uxehhak02py4rcz65zu0swh7wjngq88t',
            enable_auto_deployment: false,
          },
        },
        transfer:: super['transfer'],
      },
      consensus+: {
        params+: {
          block+: {
            time_iota_ms: '2000',
          },
        },
      },
    },
  },
}
