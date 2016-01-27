package testutil

import (
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	syncds "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-datastore/sync"
	ds2 "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/util/datastore2"
)

func ThreadSafeCloserMapDatastore() ds2.ThreadSafeDatastoreCloser {
	return ds2.CloserWrap(syncds.MutexWrap(datastore.NewMapDatastore()))
}
