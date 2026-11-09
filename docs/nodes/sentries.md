---
title: Configure validator sentries
description: Separate a validator from public P2P connections with sentry full nodes.
slug: /nodes/sentries
---

# Configure validator sentries

A sentry runs `upsided` as a full node and relays blocks, transactions, and consensus messages. Its consensus key is not registered as a validator. Run each sentry with its own home, P2P key, and consensus key.

A validator connects persistently to its sentries. The sentries maintain connections to the wider network. Use multiple sentries so one unavailable host does not disconnect the validator.

## Public network connections

New mainnet full nodes can use the seed and public AWS sentries listed in [Connect to mainnet](../networks/mainnet.md). Those endpoints provide access to the network. To hide a validator behind sentries, configure the validator and the sentries you operate together as described below.

## Private validator

Set the following entries in the validator's `config/config.toml`. Replace the example endpoints with your own sentry node IDs and reachable addresses:

```toml
[p2p]
pex = false
seed_mode = false
seeds = ""
external_address = ""
persistent_peers = "SENTRY_A_ID@SENTRY_A_HOST:26656,SENTRY_B_ID@SENTRY_B_HOST:26656"
unconditional_peer_ids = "SENTRY_A_ID,SENTRY_B_ID"
```

Restrict inbound validator P2P traffic to its sentries using the host or network firewall. Keep RPC, API, gRPC, and profiling services private or disabled. Restrict outbound connectivity according to your infrastructure requirements.

`unconditional_peer_ids` reserves connectivity despite peer-count limits. It does not authenticate additional peers or replace firewall rules.

## Sentry nodes

Set these entries on each of your sentries:

```toml
[p2p]
pex = true
seed_mode = false
persistent_peers = "VALIDATOR_ID@PRIVATE_VALIDATOR_HOST:26656,OTHER_SENTRY_ID@OTHER_SENTRY_HOST:26656"
private_peer_ids = "VALIDATOR_ID"
unconditional_peer_ids = "VALIDATOR_ID,OTHER_SENTRY_ID"
```

List every protected validator's P2P node ID in `private_peer_ids`. This prevents those addresses from being advertised through peer exchange. If the sentry connects to peers using private, non-routable addresses, set `addr_book_strict = false` for that topology.

Configure the sentry's listening and advertised addresses for its public P2P endpoint. Enable public access only to services you intend to provide. Keep management interfaces private.

Never copy the validator's consensus private key onto a sentry. A sentry's own `config/priv_validator_key.json`, generated during initialization, belongs to that full node and must not be reused as the validator's signing key.

## Seed nodes

A dedicated seed uses `pex = true` and `seed_mode = true`. It returns peer addresses and disconnects after discovery. Configure seeds in `seeds`. Use ordinary sentries in `persistent_peers` for ongoing block and transaction exchange.

## Verify the topology

Before registering or migrating a validator:

1. Check that every node reports the intended chain ID and reaches the current height.
2. Query local `/net_info` to confirm the intended persistent connections.
3. Check the validator's listening interfaces and ingress rules.
4. Verify that each sentry lists the validator in `private_peer_ids`.
5. Verify that only the intended validator signer has the active consensus key.

After activation, monitor block participation using [Node and validator operations](operations.md).
