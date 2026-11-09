#!/bin/bash

make install

CHAINID="upside_61872-1"
MONIKER="localtestnet"

# Test-only mnemonics for local networks and integration tests. They are public.
# Never send funds to these accounts and never use them on mainnet or testnet.

# localKey address 0x7cb61d4117ae31a12e393a1cfa3bac666481d02e
VAL_KEY="localkey"
VAL_MNEMONIC="gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat"

# user1 address 0xc6fe5d33615a1c52c08018c47e8bc53646a0e101
USER1_KEY="user1"
USER1_MNEMONIC="copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom"

# user2 address 0x963ebdf2e1f8db8707d05fc75bfeffba1b5bac17
USER2_KEY="user2"
USER2_MNEMONIC="maximum display century economy unlock van census kite error heart snow filter midnight usage egg venture cash kick motor survey drastic edge muffin visual"

# remove existing daemon and client
rm -rf ~/.upside*

# Import keys from mnemonics
echo $VAL_MNEMONIC | upsided keys add $VAL_KEY --recover --keyring-backend test --algo "eth_secp256k1"
echo $USER1_MNEMONIC | upsided keys add $USER1_KEY --recover --keyring-backend test --algo "eth_secp256k1"
echo $USER2_MNEMONIC | upsided keys add $USER2_KEY --recover --keyring-backend test  --algo "eth_secp256k1"

upsided init $MONIKER --chain-id $CHAINID --default-denom aup

# Set gas limit in genesis
#cat $HOME/.upside/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.upside/config/tmp_genesis.json && mv $HOME/.upside/config/tmp_genesis.json $HOME/.upside/config/genesis.json
# cat $HOME/.upside/config/genesis.json | jq '.app_state["upside"]["params"]["ibc_upside_denom"]="aup"' > $HOME/.upside/config/tmp_genesis.json && mv $HOME/.upside/config/tmp_genesis.json $HOME/.upside/config/genesis.json
# cat $HOME/.upside/config/genesis.json | jq '.app_state["gravity"]["params"]["ethereum_blacklist"]=["0x000000000000000000000000000000000000dEaD"]' > $HOME/.upside/config/tmp_genesis.json && mv $HOME/.upside/config/tmp_genesis.json $HOME/.upside/config/genesis.json
cat $HOME/.upside/config/genesis.json | jq '.consensus["params"]["block"]["max_gas"]="10000000" | .app_state.evm.params.evm_denom = "aup"' > $HOME/.upside/config/tmp_genesis.json && mv $HOME/.upside/config/tmp_genesis.json $HOME/.upside/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
upsided genesis add-genesis-account "$(upsided keys show $VAL_KEY -a --keyring-backend test)" 10000000000000000000000aup --keyring-backend test
upsided genesis add-genesis-account "$(upsided keys show $USER1_KEY -a --keyring-backend test)" 10000000000000000000000aup --keyring-backend test
upsided genesis add-genesis-account "$(upsided keys show $USER2_KEY -a --keyring-backend test)" 10000000000000000000000aup --keyring-backend test

# Sign genesis transaction
upsided genesis gentx $VAL_KEY 1000000000000000000aup --amount=1000000000000000000000aup --chain-id $CHAINID --keyring-backend test

# Collect genesis tx
upsided genesis collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
upsided genesis validate

echo "blocked-addresses = ['up16z0herz998946wr659lr84c8c556da55dc34hh']" >> $HOME/.upside/config/app.toml

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
upsided start --pruning=nothing --rpc.unsafe --log_level info --json-rpc.api eth,txpool,personal,net,debug,web3 --minimum-gas-prices 200000aup --api.enable
