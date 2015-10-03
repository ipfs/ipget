package mdutils

import (
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/blocks/blockstore"
	bsrv "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/blockservice"
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/exchange/offline"
	dag "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/merkledag"
	ds "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jbenet/go-datastore"
	dssync "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jbenet/go-datastore/sync"
)

func Mock() dag.DAGService {
	bstore := blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))
	bserv := bsrv.New(bstore, offline.Exchange(bstore))
	return dag.NewDAGService(bserv)
}
