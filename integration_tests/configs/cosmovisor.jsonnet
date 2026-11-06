local config = import 'default.jsonnet';

config {
  'upside_61872-1'+: {
    'app-config'+: {
      'app-db-backend': 'rocksdb',
      'minimum-gas-prices': '100000000000taup',
      'iavl-lazy-loading':: super['iavl-lazy-loading'],
    },
    validators: [super.validators[0] {
      'app-config'+: {
        store: {
          streamers: ['versiondb'],
        },
      },
    }] + super.validators[1:],
    genesis+: {
      app_state+: {
        bank+: {
          params: {
            send_enabled: [
              {
                denom: 'aup',
                enabled: true,
              },
              {
                denom: 'taup',
                enabled: false,
              },
            ],
          },
        },
        feemarket+: {
          params+: {
            no_base_fee: false,
            base_fee:: super.base_fee,
          },
        },
        gov: {
          voting_params: {
            voting_period: '10s',
          },
          deposit_params: {
            max_deposit_period: '10s',
            min_deposit: [
              {
                denom: 'taup',
                amount: '1',
              },
            ],
          },
        },
      },
    },
  },
}
