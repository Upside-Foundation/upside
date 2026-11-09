<!-- order: 1 -->

# Concepts

## Upside

Upside is an EVM chain built on the Cosmos SDK with Ethermint. Users deploy and
call smart contracts with standard Ethereum tools, and the chain connects to
other Cosmos chains over IBC.

The `x/upside` module connects IBC and the EVM. It keeps a mapping between
native Cosmos denoms and EVM token contracts, converts tokens between the two
forms, and runs hooks when EVM contracts emit certain events. It also holds the
two parameters that control the native token's supply curve,
`mint_half_life_seconds` and `mint_paused` (see [Supply Curve](08_supply.md)).

On mainnet, governance proposal 6 turned off IBC transfers: the `transfer`
module's `send_enabled` and `receive_enabled` parameters are `false`, and
mainnet has no IBC channels. While that holds, `MsgTransferTokens` fails, and
no IBC vouchers arrive for `MsgConvertVouchers` to convert.

## Native Token

UP is the native token. Its base denom is `aup`, with 18 decimals
(1 UP = 10^18 `aup`), so amounts match Ethereum's units. `aup` pays gas, is
staked by validators and delegators, and is used for governance deposits.
`aup` is minted by `x/mint` according to the capped supply curve.

### IBC Voucher Conversion

Inherited from Cronos: when an IBC voucher whose denom equals the
`ibc_upside_denom` parameter arrives, it is converted to `aup` automatically.
That voucher is assumed to have 8 decimals:

- On arrival, the voucher is held in the module account and the same amount
  is minted as `aup`, multiplied by 10^10.
- On transfer out, the `aup` is burned and the voucher is released. Amounts
  are truncated to 8 decimals; the remainder stays on Upside.

This path only runs if an IBC channel delivers a voucher with exactly that
denom. The default `ibc_upside_denom` is inherited from Cronos, where it
is the denom of the CRO voucher. No Upside channel delivers it.

## Native Denoms Handled by the Module

- `aup`: the native token.
- IBC vouchers, `ibc/<hash>`: tokens received over IBC channels.
- Source tokens, `up0x<contract address>`: native denoms for tokens that
  originate from CRC21 contracts on Upside. A CRC21 contract follows the
  interface in [`ModuleCRC21.sol`](../../../contracts/src/ModuleCRC21.sol),
  which adds events for sending tokens over IBC and an `is_source` flag that
  marks a token originating on Upside.
- `gravity0x<contract address>`: a denom format inherited from Cronos's
  gravity bridge. Token mapping still accepts it, but this chain has no
  gravity bridge module, so no tokens arrive this way.

## CRC20 Token

A CRC20 token is Upside's version of an ERC20 token, with extensions that let
the module mint and burn it. A CRC20 contract can be mapped to a native denom,
which allows tokens to move between their native and EVM forms, and onward
over IBC.

## Auto-deployed Contract

A contract whose bytecode is embedded in the module and deployed by it. Some of
its operations can only be called by the module, so the module can trust it
directly. The module can auto-deploy a minimal CRC20 contract
([`ModuleCRC20.sol`](../../../contracts/src/ModuleCRC20.sol)) to wrap a native
token in the EVM automatically.

## Token Mapping

The module keeps two mappings from native denom to contract address: one for
auto-deployed contracts and one for external contracts.

When auto-deployment is enabled, incoming IBC tokens are wrapped in an
auto-deployed CRC20 contract automatically.

An external contract mapping can be registered for a denom through a
governance proposal or by an account with the token mapping permission.

## Admin and Permissions

The `upside_admin` parameter names an admin account, which has every
permission. The admin can grant permissions to other accounts with
`MsgUpdatePermissions`. Permissions are bit flags:

| Value | Permission              | Allows                                        |
| ----- | ----------------------- | --------------------------------------------- |
| 1     | `CanChangeTokenMapping` | `MsgUpdateTokenMapping`                       |
| 2     | `CanTurnBridge`         | `MsgTurnBridge`, which currently does nothing |
| 3     | `All`                   | Both of the above                             |

## Block List

The admin can store an address block list with `MsgStoreBlockList`. The list
is stored encrypted with [age](https://age-encryption.org). Validators that
hold the matching key decrypt it in the block proposal handler
([`app/proposal.go`](../../../app/proposal.go)): they leave transactions
involving blocked addresses, including addresses in EIP-7702 authorization
lists, out of the blocks they propose, and reject proposals that contain them.

This is separate from the `blocked-addresses` option in a node's `app.toml`,
which the node checks in its ante handler.
