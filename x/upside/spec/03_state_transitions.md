<!--
order: 3
-->

# State Transitions

<!-- define state transitions for accounts and storage -->

This document describes how the token mapping state changes.

## Token Mapping

### Create/Update

There are several scenarios when the token mapping is updated:

- When auto-deployment is enabled and a new IBC token arrives.
- When auto-deployment is enabled and `MsgConvertVouchers` is executed with a new token.
- When a `TokenMappingChangeProposal` gov proposal is executed.
- When an account with the token mapping permission executes `MsgUpdateTokenMapping`.

When a mapping is updated, moving balances from the old contract to the new one is up to the contract developer. For example, the new contract could let users swap old tokens for new ones.

#### Contract Migration

When a token mapping is updated, the old contract still exists.

### Delete

A `TokenMappingChangeProposal` with an empty `contract` deletes the external
mapping for its denom, if the denom is not a source token (`up0x...`).
`MsgUpdateTokenMapping` can't delete a mapping, because its basic validation
rejects an empty contract address. A mapping to a contract that the module
deployed automatically can't be deleted.
