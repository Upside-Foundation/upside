//go:build !objstore
// +build !objstore

package rootmulti

import (
	"fmt"

	"cosmossdk.io/store/types"
	"github.com/0xUp5ide/upside-foundation/memiavl"
)

func (rs *Store) GetObjKVStore(key types.StoreKey) types.ObjKVStore {
	panic("placeholder does not support GetObjKVStore")
}

func (rs *Store) loadExtraStore(db *memiavl.DB, key types.StoreKey, params storeParams) (types.CommitStore, error) {
	panic(fmt.Sprintf("unrecognized store type %v", params.typ))
}
