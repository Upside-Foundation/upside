<!--
order: 6
-->

# Events

The module emits Cosmos SDK [events](https://docs.cosmos.network/main/learn/advanced/events) when
its messages execute. As with every Cosmos SDK message in a transaction, the first event for each
message has type `message`, emitted by `BaseApp`. Its `action` attribute is the message type URL, and
its `sender` attribute is the signer. When the handler emits no `module` attribute, `BaseApp` adds one
taken from the type URL.

## MsgConvertVouchers

| Type             | Attribute Key | Attribute Value            |
| ---------------- | ------------- | -------------------------- |
| convert_vouchers | `"sender"`    | `{bech32_address}`         |
| convert_vouchers | `"amount"`    | `{amount}`                 |
| message          | action        | /upside.MsgConvertVouchers |
| message          | sender        | {bech32_address}           |
| message          | module        | upside                     |

## MsgTransferTokens

| Type            | Attribute Key | Attribute Value           |
| --------------- | ------------- | ------------------------- |
| transfer_tokens | `"sender"`    | `{bech32_address}`        |
| transfer_tokens | `"recipient"` | `{bech32_address}`        |
| transfer_tokens | `"amount"`    | `{amount}`                |
| message         | action        | /upside.MsgTransferTokens |
| message         | sender        | {bech32_address}          |
| message         | module        | upside                    |

## MsgUpdateTokenMapping

| Type    | Attribute Key | Attribute Value               |
| ------- | ------------- | ----------------------------- |
| message | action        | /upside.MsgUpdateTokenMapping |
| message | sender        | {bech32_address}              |
| message | module        | MsgUpdateTokenMapping         |

The handler emits no `module` attribute, so `BaseApp` takes it from the type URL, and the value is
`MsgUpdateTokenMapping` rather than `upside`.

## MsgSetMintPaused

| Type            | Attribute Key | Attribute Value   |
| --------------- | ------------- | ----------------- |
| set_mint_paused | `"paused"`    | `true` or `false` |

This message runs only through a governance proposal, so `BaseApp` emits no `message` event for it.
The `set_mint_paused` event is among the events of the block in which the proposal executes.
