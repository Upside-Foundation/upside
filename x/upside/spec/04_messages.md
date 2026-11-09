<!-- order: 4 -->

# Messages

All messages are defined in [`proto/upside/tx.proto`](../../../proto/upside/tx.proto).

## MsgConvertVouchers

> Users normally do this through a smart contract and do not send this message directly.

Converts native tokens to their mapped CRC20 tokens. If no mapping exists and
auto-deployment is enabled, the module deploys its embedded CRC20 contract for
the denom first. A voucher whose denom equals `ibc_upside_denom` is converted
to `aup` instead (see [Concepts](01_concepts.md#ibc-voucher-conversion)).

This message fails if:

- The coin denom is not an IBC voucher (`ibc/...`), a source token
  (`up0x...`), or a `gravity0x...` denom.
- No mapping exists and auto-deployment is disabled.

Fields:

- `address`: the signer, a bech32 address on Upside.
- `coins`: the coins to convert.

## MsgTransferTokens

> Users normally do this through a smart contract and do not send this message directly.

Transfers IBC tokens out of Upside. For the `ibc_upside_denom` voucher, the
decimal conversion is done automatically.

It calls the IBC transfer module with `timeoutHeight` set to zero and
`timeoutTimestamp` set from the `ibc_timeout` parameter.

This message fails if:

- The sender does not have enough balance.
- The IBC transfer fails.

Fields:

- `from`: the signer, a bech32 address on Upside.
- `to`: the destination address on the other chain.
- `coins`: the coins to transfer.

## MsgUpdateTokenMapping

Creates or updates an external token mapping. Only the admin or an account
with the `CanChangeTokenMapping` permission can send it.

This message fails if:

- The sender is not authorized.
- The contract address or denom is malformed. An empty contract address is
  rejected, so this message can't delete a mapping (see
  [State Transitions](03_state_transitions.md#delete)).
- The contract is already mapped to another denom.

Fields:

- `sender`: the signer, a bech32 address on Upside.
- `denom`: the denom to map.
- `contract`: the hex address of the contract.
- `symbol`: the token symbol. Used only for source tokens (`up0x...`).
- `decimal`: the number of decimals. Used only for source tokens.

## MsgTurnBridge

Inherited from Cronos, where it enabled or disabled the gravity bridge. This
chain has no gravity bridge, and the message handler does nothing.

## MsgUpdateParams

Replaces all module parameters. Only the governance module account can send
it, so it runs as part of a passed governance proposal. Every parameter must
be supplied; see [Parameters](07_params.md) for validation.

## MsgSetMintPaused

Sets `mint_paused` to pause or resume `x/mint`. Only the governance module
account can send it, through a governance proposal. Fields are `authority`
(the governance account) and `paused` (the requested boolean value).

This is the only message that changes `mint_paused`; `MsgUpdateParams` must
repeat the stored value. Resuming minting does not mint the amounts skipped
while paused. See [Supply](08_supply.md) and [Parameters](07_params.md).

## MsgUpdatePermissions

Sets the permission bits for an account (see
[Concepts](01_concepts.md#admin-and-permissions)). Only the admin can send it.

## MsgStoreBlockList

Stores the encrypted block list (see [Concepts](01_concepts.md#block-list)).
Only the admin can send it.
