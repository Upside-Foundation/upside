local ibc = import 'ibc.jsonnet';

ibc {
  'upside_61872-1'+: {
    key_name: 'signer3',
    accounts: super.accounts + [{
      name: 'signer3',
      coins: '0taup',
      mnemonic: '${SIGNER3_MNEMONIC}',
    }],
    genesis+: {
      app_state+: {
        upside+: {
          params+: {
            ibc_timeout: 0,
          },
        },
      },
    },
  },
  relayer+: {
    chains: [super.chains[0] {
      fee_granter: 'up16z0herz998946wr659lr84c8c556da55wnla3a', //signer1
    }] + super.chains[1:],
  },
}
