{
  dotenv: '../../scripts/.env',
  'upside_61872-1': {
    cmd: 'upsided',
    'start-flags': '--trace',
    config: {
      db_backend: 'rocksdb',
      mempool: {
        version: 'v1',
      },
    },
    'app-config': {
      chain_id: 'upside_61872-1',
      'app-db-backend': 'rocksdb',
      'minimum-gas-prices': '0aup',
      'index-events': ['ethereum_tx.ethereumTxHash'],
      'iavl-lazy-loading': true,
      'json-rpc': {
        address: '127.0.0.1:{EVMRPC_PORT}',
        'ws-address': '127.0.0.1:{EVMRPC_PORT_WS}',
        api: 'eth,net,web3,debug,upside',
        'feehistory-cap': 100,
        'block-range-cap': 10000,
        'logs-cap': 10000,
      },
      evm: {
        'block-executor': 'sequential',
      },
      mempool: {
        'max-txs': 1000,
        'feebump': 10,
      },
    },
    validators: [{
      coins: '10001000000000000000000taup',
      staked: '1000000000000000000taup',
      mnemonic: '${VALIDATOR1_MNEMONIC}',
      client_config: {
        'broadcast-mode': 'sync',
      },
      'app-config': {
        memiavl: {
          enable: true,
          'zero-copy': true,
          'snapshot-interval': 5,
          'cache-size': 0,
          'async-commit-buffer': 5,
        },
        versiondb: {
          enable: true,
        },
        evm: {
          'block-executor': 'block-stm',
          'block-stm-workers': 32,
        },
      },
    }, {
      coins: '10001000000000000000000taup',
      staked: '1000000000000000000taup',
      mnemonic: '${VALIDATOR2_MNEMONIC}',
      client_config: {
        'broadcast-mode': 'sync',
      },
      config: {
        db_backend: 'pebbledb',
      },
      'app-config': {
        'app-db-backend': 'pebbledb',
      },
    }, {
      coins: '10001000000000000000000taup',
      staked: '1000000000000000000taup',
      mnemonic: '${VALIDATOR3_MNEMONIC}',
      client_config: {
        'broadcast-mode': 'sync',
      },
      config: {
        db_backend: 'goleveldb',
      },
      'app-config': {
        'app-db-backend': 'goleveldb',
      },
    }],
    accounts: [{
      name: 'community',
      coins: '10000000000000000000000taup',
      mnemonic: '${COMMUNITY_MNEMONIC}',
    }, {
      name: 'signer1',
      coins: '20000000000000000000000taup',
      mnemonic: '${SIGNER1_MNEMONIC}',
    }, {
      name: 'signer2',
      coins: '30000000000000000000000taup',
      mnemonic: '${SIGNER2_MNEMONIC}',
    }],
    genesis: {
      consensus: {
        params: {
          block: {
            max_bytes: '1048576',
            max_gas: '81500000',
          },
        },
      },
      app_state: {
        evm: {
          params: {
            evm_denom: 'taup',
          },
        },
        upside: {
          params: {
            upside_admin: '${UPSIDE_ADMIN}',
            enable_auto_deployment: true,
            ibc_upside_denom: '${IBC_UP_DENOM}',
          },
        },
        e2ee: {
          keys: [{
            address: 'up16z0herz998946wr659lr84c8c556da55wnla3a',
            key: 'age1k3mpspxytgvx6e0jja0xgrtzz7vw2p00c2a3xmq5ygfzhwh4wg0s35z4c8',
          }],
        },
        gov: {
          params: {
            expedited_voting_period: '1s',
            voting_period: '10s',
            max_deposit_period: '10s',
            min_deposit: [
              {
                denom: 'taup',
                amount: '1',
              },
            ],
          },
        },
        transfer: {
          params: {
            receive_enabled: true,
            send_enabled: true,
          },
        },
        staking: {
          params: {
            bond_denom: 'taup',
          },
        },
        feemarket: {
          params: {
            no_base_fee: false,
            base_fee: '100000000000',
          },
        },
      },
    },
  },
}
