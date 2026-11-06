<!--
order: 7
-->

# Parameters

The module has the following parameters, defined in
[`proto/upside/upside.proto`](../../../proto/upside/upside.proto). Governance
changes them with `MsgUpdateParams`, which must include every parameter.

| Key                      | Type   | Default                                                                  |
| ------------------------ | ------ | ------------------------------------------------------------------------ |
| `ibc_upside_denom`       | string | `"ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865"` |
| `ibc_timeout`            | uint64 | `86400000000000` (1 day, in nanoseconds)                                 |
| `upside_admin`           | string | `""`                                                                     |
| `enable_auto_deployment` | bool   | `false`                                                                  |
| `max_callback_gas`       | uint64 | `50000`                                                                  |
| `mint_half_life_seconds` | uint64 | `63115200` (2 Julian years)                                              |
| `mint_paused`            | bool   | `false`                                                                  |

- `ibc_upside_denom`: the IBC voucher denom that is converted to `aup` on
  arrival (see [Concepts](01_concepts.md#ibc-voucher-conversion)). The default
  is inherited from Cronos. Changing it does not move tokens already
  converted; that would need an explicit migration.

- `ibc_timeout`: the timeout used when the module sends tokens over IBC.

- `upside_admin`: the admin account (see
  [Concepts](01_concepts.md#admin-and-permissions)). Empty means no admin. If
  set, it must be a valid bech32 address.

- `enable_auto_deployment`: whether the module auto-deploys a CRC20 contract
  for new tokens. When disabled and no external contract is mapped, incoming
  tokens stay as native tokens, and users can move them with Cosmos messages.
  Tokens converted while it was enabled can still be withdrawn after it is
  disabled.

- `max_callback_gas`: the gas limit for the EVM call made when an IBC packet
  sent by a contract is acknowledged or times out.

- `mint_half_life_seconds`: the half-life of the native token's supply curve
  (see [Supply Curve](08_supply.md)). Must be at least `31557600`
  (1 Julian year), so that one proposal cannot mint most of the remaining
  supply in a short time.

- `mint_paused`: while `true`, `x/mint` mints nothing (see
  [Supply Curve](08_supply.md#pausing-minting)). Only `MsgSetMintPaused`
  changes it. `MsgUpdateParams` fails unless its `mint_paused` equals the
  stored value, so a params proposal must repeat the current value.
