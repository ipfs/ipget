package bitswap

import (
	bsnet "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/exchange/bitswap/network"
	peer "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/p2p/peer"
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/util/testutil"
)

type Network interface {
	Adapter(testutil.Identity) bsnet.BitSwapNetwork

	HasPeer(peer.ID) bool
}
