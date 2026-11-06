//go:build rocksdb
// +build rocksdb

package cmd

import (
	"sort"

	"github.com/0xUp5ide/upside-foundation/app"
	"github.com/0xUp5ide/upside-foundation/cmd/upsided/opendb"
	versiondbclient "github.com/0xUp5ide/upside-foundation/versiondb/client"
	"github.com/linxGnu/grocksdb"
	"github.com/spf13/cobra"
)

func ChangeSetCmd() *cobra.Command {
	keys, _, _ := app.StoreKeys()
	storeNames := make([]string, 0, len(keys))
	for name := range keys {
		storeNames = append(storeNames, name)
	}
	sort.Strings(storeNames)

	return versiondbclient.ChangeSetGroupCmd(versiondbclient.Options{
		DefaultStores:  storeNames,
		OpenReadOnlyDB: opendb.OpenReadOnlyDB,
		AppRocksDBOptions: func(sstFileWriter bool) *grocksdb.Options {
			return opendb.NewRocksdbOptions(nil, sstFileWriter)
		},
	})
}
