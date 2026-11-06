<!--
order: 6
-->

# Events

The module emits Cosmos SDK [events](https://docs.cosmos.network/main/learn/advanced/events) when
its messages execute. As with every Cosmos SDK message, the first event for each message has type
`message` and an `action` attribute, emitted by `BaseApp`.

## MsgConvertVouchers

| Type             | Attribute Key | Attribute Value    |
| ---------------- | ------------- | ------------------ |
| convert_vouchers | `"sender"`    | `{bech32_address}` |
| convert_vouchers | `"amount"`    | `{amount}`         |
| message          | module        | upside             |
| message          | action        | ConvertVouchers    |

## MsgTransferTokens

| Type            | Attribute Key | Attribute Value    |
| --------------- | ------------- | ------------------ |
| transfer_tokens | `"sender"`    | `{bech32_address}` |
| transfer_tokens | `"recipient"` | `{bech32_address}` |
| transfer_tokens | `"amount"`    | `{amount}`         |
| message         | module        | upside             |
| message         | action        | TransferTokens     |

## MsgUpdateTokenMapping

| Type    | Attribute Key | Attribute Value    |
| ------- | ------------- | ------------------ |
| message | action        | UpdateTokenMapping |
