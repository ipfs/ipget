package mockrouting

import (
	ds "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	sync "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-datastore/sync"
	mocknet "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/p2p/net/mock"
	dht "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/routing/dht"
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/util/testutil"
	context "github.com/noffle/ipget/Godeps/_workspace/src/golang.org/x/net/context"
)

type mocknetserver struct {
	mn mocknet.Mocknet
}

func NewDHTNetwork(mn mocknet.Mocknet) Server {
	return &mocknetserver{
		mn: mn,
	}
}

func (rs *mocknetserver) Client(p testutil.Identity) Client {
	return rs.ClientWithDatastore(context.TODO(), p, ds.NewMapDatastore())
}

func (rs *mocknetserver) ClientWithDatastore(ctx context.Context, p testutil.Identity, ds ds.Datastore) Client {

	// FIXME AddPeer doesn't appear to be idempotent

	host, err := rs.mn.AddPeer(p.PrivateKey(), p.Address())
	if err != nil {
		panic("FIXME")
	}
	return dht.NewDHT(ctx, host, sync.MutexWrap(ds))
}

var _ Server = &mocknetserver{}
