---
title: Rewards, the governance treasury, and governance
description: Understand reward allocation, Foundation validators, the governance treasury, and governance.
slug: /staking/governance
---

# Rewards, the governance treasury, and governance

## Staking

Rewards pass through three steps:

1. Mainnet currently allocates 93% to the governance treasury and 7% to staking. Governance can change this split, which is the `community_tax` parameter.
2. Active validators receive a share of the staking allocation based on consensus voting power.
3. Each validator takes its commission. Delegators, including the validator's self-delegation, share the remainder according to their stake.

For example, an allocation of 100 UP leaves 93 UP for the governance treasury and 7 UP for staking, before rounding. The active validator set shares the 7 UP according to consensus voting power.

If your validator receives 1 UP from that staking allocation and sets 10% commission, it earns 0.1 UP in commission. Its delegators share the other 0.9 UP.

The `community_tax` parameter applies to the collected distribution funds, including minted rewards and transaction fees sent to the fee collector. Applications can route their own fees to other accounts.

## Foundation validators

Foundation validators use the same staking, commission, and reward-distribution rules as other validators. Any additional contribution of their rewards to the community pool is a separate transfer under Foundation policy, made after protocol reward distribution.

Governance proposals can allocate community-pool funds to validator operators. Read the proposal's messages to identify recipients, amounts, and any conditions before voting.

## The governance treasury

The governance treasury holds funds for proposals approved through on-chain governance. Intended uses include community projects, protocol development, and bug bounties.

Use the [network profile](../getting-started/networks.md) to inspect the current allocation and treasury balance:

```bash
"$UPSIDE_BIN" query distribution params \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json

"$UPSIDE_BIN" query distribution community-pool \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

The CLI query name `community-pool` and the parameter name `community_tax` are protocol identifiers. This documentation calls the governed funds the governance treasury.

A `community_tax` value of `0.930000000000000000` represents 93%. The CLI can render this decimal as the scaled integer `930000000000000000`. Divide that representation by 10^18. Governance can change parameters, so use the current chain state when checking a historical or future reward estimate.

## Governance

Governance proposals can fund projects, change supported parameters, and schedule software upgrades. The governance parameters set the deposit and voting rules for all proposals on the chain.

Inspect the current rules and proposals:

```bash
"$UPSIDE_BIN" query gov params \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json

"$UPSIDE_BIN" query gov proposals \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

You can read proposal messages, voting deadline, status, and any scheduled activation height. A software-upgrade vote also requires validators to install the approved binary at the scheduled height.
