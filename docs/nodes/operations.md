---
title: Node and validator operations
description: Run, monitor, restart, back up, upgrade, and troubleshoot an Upside L1 node.
slug: /nodes/operations
---

# Node and validator operations

## Run the node as a service

Use an unprivileged service account with access to the intended node home. Set an absolute binary path and an explicit `--home` in the service configuration. Use a separate service and home for each network.

The service must preserve the startup settings used during synchronization, including listening addresses, enabled APIs, and node-local gas prices. Configure graceful shutdown, bounded restart behavior, log collection, and alerts. Ensure an upgrade halt cannot become a rapid restart loop without an operator noticing.

For a validator, only one process or signing service may use the consensus key at a time. Disable any old service or automated restart path before a migration. A standby machine must not become a second active signer.

On Linux with systemd, a user service can run the mainnet profile from the dedicated node user's home. Save this as `~/.config/systemd/user/upsided-mainnet.service`:

```ini
[Unit]
Description=Upside mainnet node
StartLimitIntervalSec=300
StartLimitBurst=3

[Service]
Type=simple
ExecStart=%h/.local/bin/upsided-mainnet start --home %h/.upside-mainnet --minimum-gas-prices 0aup --rpc.pprof_laddr= --api.enable=false --grpc.enable=true --grpc.address 127.0.0.1:9090 --grpc-web.enable=false --json-rpc.enable=false
Restart=on-failure
RestartSec=10
KillSignal=SIGINT
TimeoutStopSec=300
LimitNOFILE=65535
NoNewPrivileges=true

[Install]
WantedBy=default.target
```

`%h` is the home of the user running the service. Create the unit directory first, and adjust the binary, node home, ports, and gas-price setting to match the configuration you synchronized. P2P settings remain in `config.toml`.

Stop the foreground node before enabling the service:

```bash
systemctl --user daemon-reload
systemctl --user enable --now upsided-mainnet.service
systemctl --user status upsided-mainnet.service
journalctl --user -u upsided-mainnet.service -f
```

For a service that must survive logout and start at boot, an administrator can enable lingering for the dedicated node user with `loginctl enable-linger USERNAME`. For testnet, use a separate unit, home, and ports.


## Monitor health


| Signal                      | Healthy behavior / reason to alert                                          |
| --------------------------- | --------------------------------------------------------------------------- |
| Chain identity              | Matches the intended network                                                |
| Latest height and timestamp | Advance and remain near trusted peers                                       |
| Sync status                 | `catching_up: false` after initial synchronization                          |
| P2P peers                   | Sufficient connectivity for this deployment topology                        |
| Validator state             | Bonded with `jailed: false` and expected consensus voting power             |
| Signing behavior            | Missed-block count stays comfortably below the network's downtime threshold |
| System resources            | CPU, RAM, I/O latency, free disk, and network traffic have headroom         |
| Time synchronization        | Host clock remains synchronized                                             |
| Release and governance      | Upcoming upgrades and parameter changes are noticed in time                 |


Inspect a local node using its configured RPC port:

```bash
curl --fail --silent --show-error http://127.0.0.1:26657/status | jq '.result'
curl --fail --silent --show-error http://127.0.0.1:26657/net_info \
  | jq '.result.n_peers'
```

For the tutorial localnet, use port `27657`. Check the validator's staking record through the [validator guide](../validators/create.md). Use the node's consensus address to query its signing record directly. The commands in this guide use `NODE_HOME` from the [mainnet](../networks/mainnet.md) or [testnet](../networks/testnet.md) node profile, `UPSIDE_BIN`, `CHAIN_ID`, `RPC_URL`, `WALLET_HOME`, `KEYRING_BACKEND`, and `GAS_PRICES` from [Networks and connection settings](../getting-started/networks.md), and `KEY_NAME` from [the wallet guide](../getting-started/wallets.md):

```bash
CONSENSUS_ADDRESS="$("$UPSIDE_BIN" comet show-address --home "$NODE_HOME")"
"$UPSIDE_BIN" query slashing signing-info "$CONSENSUS_ADDRESS" \
  --node "$RPC_URL" --home "$NODE_HOME" --output json
"$UPSIDE_BIN" query slashing params \
  --node "$RPC_URL" --home "$NODE_HOME" --output json
```

Compare `missed_blocks_counter` with the current `signed_blocks_window` and `min_signed_per_window`. Inspect `jailed_until` and `tombstoned` before an unjail attempt. To enumerate all records, use `query slashing signing-infos --page-limit 100`. Pass each response's `pagination.next_key` to `--page-key` until it is empty. Key encodings follow the installed command's help.

To check block votes, choose a height at least two blocks below the current height and query `/commit?height=HEIGHT`. Check that the response has `canonical: true`, then look for the validator's hexadecimal `validator_info.address` from `/status`. A vote for the block has a nonempty signature and `block_id_flag: 2`. A flag of `3` records a nil vote.

The CLI command `comet show-address` returns a Bech32 consensus address, which begins with `upvalcons1` on both public networks. Use the hexadecimal address from RPC status when comparing RPC commit records.

## Downtime and unjailing

Fix the cause of downtime before trying to unjail. Common causes include lost peers, slow or full disks, clock errors, a wrong binary version, or a stopped signer. Query the current slashing parameters and validator record. Wait until the jail period has elapsed, and check that the validator meets its minimum self-delegation.

Once the node is synchronized and the validator is eligible, submit from its operator wallet:

```bash
"$UPSIDE_BIN" tx slashing unjail \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

Check the validator's bonded and jailed status, and verify that signing resumes. Any penalty already applied remains in effect. The validator must still meet the active-set requirements. A tombstoned validator is permanently excluded from ordinary unjailing.

## Backups and migration

Back up the operator wallet recovery material separately from node configuration. For a file-based consensus signer, protect the consensus private key and preserve the latest signing state during a controlled shutdown and migration. Database backups must follow the selected database's consistency requirements.

Preserve the latest signing state when restoring chain data from a snapshot. Confirm that the former signer is stopped and cannot restart before bringing its replacement online. Resetting that state can allow conflicting signatures.

## Upgrades

Prepare before the upgrade height:

1. Obtain the approved release from the published source and verify its checksum.
2. Read its migration and compatibility requirements.
3. Prepare a recovery plan and rehearse the upgrade on an appropriate test network.
4. At the scheduled halt, follow the release instructions to replace the binary and resume the same node home.

Do not bypass a consensus upgrade with an unsafe skip flag. Verify height progression and validator signatures after the upgrade. Use a rollback procedure that accounts for database migrations and consensus signing state.

### Replay mainnet across the v2.0.0 upgrade

Mainnet applied the upgrade plan `v2.0.0` at height 3,898,900. A node that replays mainnet from genesis needs two binaries:

| Heights | Binary |
| --- | --- |
| 1 to 3,898,899 | The pre-v2 mainnet binary |
| 3,898,900 onward | `v2.0.0`, or a later release that keeps the `v2.0.0` upgrade handler, such as `v2.0.1` |

The pre-v2 binary isn't published yet, so you can't follow this sequence from public sources today. The current release can't replay the blocks before the upgrade by itself, as described in [Connect to mainnet](../networks/mainnet.md#release-and-synchronization).

When the pre-v2 binary reaches height 3,898,900, it stops processing blocks and logs a message that begins with the following text, followed by the plan information:

```text
UPGRADE "v2.0.0" NEEDED at height: 3898900
```

To switch binaries by hand:

1. Stop the node after it logs the message, and make sure the service doesn't restart the pre-v2 binary.
2. Verify the checksum of the `v2.0.0` binary and install it in place of the pre-v2 binary.
3. Start the node with the same home and startup settings.

To switch binaries with [Cosmovisor](https://github.com/cosmos/cosmos-sdk/tree/main/tools/cosmovisor), place both binaries under the node home before you start the node. Cosmovisor expects each binary to be named `upsided`:

```text
$NODE_HOME/cosmovisor/genesis/bin/upsided          the pre-v2 binary
$NODE_HOME/cosmovisor/upgrades/v2.0.0/bin/upsided  v2.0.0
```

Set these variables in the environment that runs Cosmovisor, such as `Environment=` lines in the systemd unit:

```bash
export DAEMON_NAME="upsided"
export DAEMON_HOME="$NODE_HOME"
export DAEMON_ALLOW_DOWNLOAD_BINARIES="false"
export DAEMON_RESTART_AFTER_UPGRADE="true"
```

Then start the node with `cosmovisor run start` in place of `upsided start`, with the same flags. When the pre-v2 binary logs the upgrade message, Cosmovisor stops it, backs up the `data` directory, and starts the `v2.0.0` binary. The backup needs free disk space about equal to the size of `data`. Verify the checksum of each binary before you place it; with `DAEMON_ALLOW_DOWNLOAD_BINARIES` set to `false`, Cosmovisor doesn't download binaries.

## Storage management

Pruning removes older application state according to a retention policy. Block history, transaction indexes, snapshots, and historical query databases have separate storage requirements. Changing one setting might not reclaim space used by another. Use the network's supported database settings and measure disk growth before changing retention.

Do not delete database files or switch storage backends in place to recover disk space. Stop the node before restoring or migrating its database. Restore a verified snapshot compatible with the selected backend into a separate home, then check its height and application hash before using it. The optional RocksDB-backed historical store has separate [versiondb instructions](https://github.com/Upside-Foundation/upside/blob/main/versiondb/README.md).

## Troubleshooting


| Problem                    | First checks                                                                                    |
| -------------------------- | ----------------------------------------------------------------------------------------------- |
| No peers                   | Peer IDs, ports, DNS, firewall, advertised address, network identity                            |
| No new blocks              | Compare independent endpoints to determine whether the problem affects your node or the network |
| App hash mismatch          | Binary version, historical upgrades, genesis, snapshot provenance                               |
| Registered but not signing | Active-set rank, jailed status, public-key match, signer connectivity                           |
| Frequent missed blocks     | CPU contention, disk latency, connectivity, clock, process restarts                             |
| Address decoding fails     | Binary prefix, network profile, and address type                                                |
| Port already in use        | Another node or API service using the same port                                                 |
| Key not found              | Wallet home, keyring backend, key name, and service/user identity                               |


Share only public diagnostic information when requesting help. Exclude recovery phrases, private keys, keyring contents, and credentials from logs and support attachments.
