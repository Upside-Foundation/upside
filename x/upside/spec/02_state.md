<!--
order: 2
-->

# State

The `x/upside` module keeps the following objects in state:

|                         | Key                                    | Value                      |
| ----------------------- | -------------------------------------- | -------------------------- |
| DenomToExternalContract | `[]byte{1} + []byte(denom)`            | 20-byte contract address   |
| DenomToAutoContract     | `[]byte{2} + []byte(denom)`            | 20-byte contract address   |
| ContractToDenom         | `[]byte{3} + 20-byte contract address` | `[]byte(denom)`            |
| Params                  | `[]byte{4}`                            | `ProtocolBuffer(Params)`   |
| AdminToPermissions      | `[]byte{5} + address bytes`            | permission bit flags, big-endian `uint64` |
| BlockList               | `[]byte{6}`                            | block list blob            |

- `DenomToExternalContract` maps a denom to an external CRC20 contract.
- `DenomToAutoContract` maps a denom to an auto-deployed CRC20 contract.
- `ContractToDenom` is the reverse map for both kinds of contract.
- `Params` holds the module parameters (see [Parameters](07_params.md)).
- `AdminToPermissions` holds the permissions granted to each account.
- `BlockList` holds the latest block list stored by the admin.
