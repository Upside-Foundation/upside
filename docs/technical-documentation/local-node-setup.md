# Local Node Setup (Development)

This repo contains two supported ways to spin up local nodes for development:

1. **Single-node local validator** via the provided script: `./start-local-node.sh`
2. **Multi-node local network** by running multiple `upsided` instances on different ports

Both run entirely locally and are intended for development/testing.

## Prerequisites

- Built `upsided` binary (see [Build & Install](./build-and-install.md))
- `jq` (used by the helper script): `sudo apt-get install -y jq`

## Option A: Single-node local validator (fastest)

Run the helper script:

```bash
./start-local-node.sh
```

What it does (at a high level):

- Builds and installs `upsided`
- Deletes any existing local state under `~/.upside*`
- Initializes a local chain with chain-id `upside_61872-1`
- Imports a few deterministic dev keys into the `test` keyring
- Creates a genesis validator (`gentx` + `collect-gentxs`)
- Starts the node with EVM JSON-RPC and REST API enabled

Important:

- The default node home for this chain is `~/.upside`.
- The script **wipes** local data (`rm -rf ~/.upside*`). Don’t use it on a machine where you care about existing local chain state.

## Option B: Multi-node local network (manual)

This section is intentionally command-focused and only includes commands that have been validated in this repo.

See [Join Local Network](./join-local-network.md) for adding a second node/validator to an already-running local network.
