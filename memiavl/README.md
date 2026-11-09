# Memory-mapped IAVL

`memiavl` stores an IAVL tree in memory-mapped snapshots and records updates in a write-ahead log. It supports state replay and Merkle proofs. The separate [VersionDB](../versiondb/README.md) can serve historical key-value queries.

## Snapshot format

A tree snapshot contains four files. Integer fields use little-endian encoding.

| File | Contents |
| --- | --- |
| `metadata` | Three 32-bit fields: magic, format, and version. Total size: 12 bytes. |
| `nodes` | Branch records, each 48 bytes. |
| `leaves` | Leaf records, each 48 bytes. |
| `kvs` | Key-value records with 32-bit byte lengths. |

The current format number is `0`. See [snapshot.go](snapshot.go) for the reader, writer, file names, and validation rules.

A branch record contains height (1 byte), `preTrees` (1 byte), padding (2 bytes), version (4 bytes), subtree size (4 bytes), key-leaf index (4 bytes), and hash (32 bytes).

A leaf record contains version (4 bytes), key length (4 bytes), key offset into `kvs` (8 bytes), and hash (32 bytes). [persisted_node.go](persisted_node.go) defines the layouts and traversal rules.

Each `kvs` record contains a 32-bit key length, key bytes, a 32-bit value length, and value bytes. Snapshot readers use the operating system's page cache for mapped files.

## Change set file

The VersionDB migration tooling exports change sets for replay. Each version has a 16-byte header containing the version and payload size, each encoded in 8 bytes. The payload contains deletion flags, variable-length key lengths and keys, and value lengths and values for entries that are not deletions.

The encoding and compression readers are defined in [versiondb/client/changeset.go](../versiondb/client/changeset.go). Exported change set files and the live `memiavl` write-ahead log have distinct encodings; use their corresponding readers.

## Operation

Snapshot and log retention affect both recovery and historical proof availability. Preserve the files required by the selected recovery procedure, and verify replayed application hashes against trusted chain data before using restored state.

Run the module tests from the repository root:

```bash
make test-memiavl
```
