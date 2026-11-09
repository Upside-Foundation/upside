# VersionDB

VersionDB stores historical state key-value pairs in a separate RocksDB database. IAVL supplies the committed root hashes and Merkle proofs. VersionDB can back historical gRPC queries, including a standalone service started with `--grpc-only`.

The implementation uses RocksDB's user-defined timestamps[^1]. Its native library and build requirements differ from the default daemon release. Choose retention policies for IAVL proofs and VersionDB queries separately.

## Build prerequisite

VersionDB requires a RocksDB-enabled binary and the native RocksDB development
libraries. The default `make build` and pinned Docker release recipe omit the
`rocksdb` tag: they cannot enable VersionDB or run `upsided changeset`.

The existing Makefile selects the RocksDB build with:

```bash
COSMOS_BUILD_OPTIONS=rocksdb make build
./build/upsided version --long
./build/upsided changeset --help
```

Install compatible RocksDB headers/libraries and a C/C++ toolchain before using
that target; the Go binding version is pinned in `go.mod`. The separate Nix
definitions also include RocksDB, but are not the default release recipe.
Native build compatibility is platform-dependent and is not exercised by the
default Go CI workflow. Confirm the `rocksdb` build tag and the commands above
before enabling VersionDB or attempting migration.

## Configuration

To enable versiondb, set the `versiondb.enable` to `true` in `app.toml`:

```toml
[versiondb]
enable = true
```

On startup, a streaming service writes state changes to `$NODE_HOME/data/versiondb`. Enabling VersionDB also selects it as the backing store for gRPC queries. Import historical state before relying on queries at older heights.

For a nonempty VersionDB, its latest version must match the application's committed version. A mismatch prevents startup. See [Catch up with IAVL](#catch-up-with-iavl-tree).

## Migration

Work on a consistent backup with the node stopped. Keep the original database until the migrated state and application hash have been verified.

Migration has two stages:

1. Extract change sets from an archived IAVL database.
2. Import those change sets into VersionDB.

Concurrency and sorter settings control memory use. Measure resource requirements on the state you intend to migrate.

### Extract Change Sets

```bash
upsided changeset dump data --home "$NODE_HOME"
```

`dump` exports change sets into one directory per store. It uses the stores registered by the application unless `--stores` selects a different list. Output is segmented by block range. With the default zlib compression and one-million-block chunks, the directory has this structure:

```
data/acc/block-0.zz
data/acc/block-1000000.zz
data/acc/block-2000000.zz
...
data/authz/block-0.zz
data/authz/block-1000000.zz
data/authz/block-2000000.zz
...
```

Verify imported change sets against the trusted application hash for their target height. Keep migration work separate from the active node's databases.

#### Verify Change Sets

```bash
upsided changeset verify data
```

`verify` replays the change sets, rebuilds the IAVL tree, and prints the application hash and commit information at the target version. Compare that result with the corresponding trusted block header. The target defaults to the latest version in the input.

Verification can run incrementally. Save an intermediate snapshot, then continue from it:

```bash
upsided changeset verify data --save-snapshot snapshot --target-version 3000000
upsided changeset verify data --load-snapshot snapshot
```

The format of change set files are documented [here](../memiavl/README.md#change-set-file).

### Build VersionDB

To maximize the speed of initial data ingestion speed into rocksdb, we take advantage of the sst file writer feature to write out sst files first, then ingest them into final db, the sst files for each store can be written out in parallel. We also developed an external sorting algorithm to sort the data before writing the sst files, so the sst files don't have overlaps and can be ingested into the bottom-most level in final db.

```bash
upsided changeset build-versiondb-sst ./data ./sst
upsided changeset ingest-versiondb-sst "$NODE_HOME/data/versiondb" sst/*.sst --move-files --maximum-version 189
```

Use `--concurrency` and `--sorter-chunk-size` to control memory use. The example height `189` must be replaced with the target version of your input.

### Restore IAVL Tree

When migrating an existing archive node to versiondb, it's recommended to rebuild the `application.db` from scratch to reclaim disk space faster, we provide a command to restore a single version of IAVL trees from memiavl snapshot.

```bash
# create memiavl snapshot
upsided changeset verify data --save-snapshot snapshot
# restore application.db
upsided changeset restore-app-db snapshot application.db
```

With the node stopped, follow the database migration plan to install the verified replacement. Preserve validator signing state independently of application database restoration.

This command creates a RocksDB `application.db`. Set `app-db-backend="rocksdb"` in `app.toml` when using that database, and verify the restored height and application hash before activation.

### Catch Up With IAVL Tree

If a nonempty VersionDB lags behind `application.db`, startup fails. Import the missing change-set range or restore a consistent pair of databases. Verify that both stores reach the intended version before starting the node. Preserve validator signing state when restoring application data.

[^1]: https://github.com/facebook/rocksdb/wiki/User-defined-Timestamp-%28Experimental%29
