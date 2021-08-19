package get

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	iface "github.com/ipfs/interface-go-ipfs-core"
	peer "github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	cleanup      []func() error
	cleanupMutex sync.Mutex
)

// Connect Gets us connected to the IPFS network
func Connect(ctx context.Context, ipfs iface.CoreAPI, peers []string) error {
	var wg sync.WaitGroup
	pinfos := make(map[peer.ID]*peer.AddrInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := pinfos[pii.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pii.ID}
			pinfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(pinfos))
	for _, pi := range pinfos {
		go func(pi *peer.AddrInfo) {
			defer wg.Done()
			log.Printf("attempting to connect to peer: %q\n", pi)
			err := ipfs.Swarm().Connect(ctx, *pi)
			if err != nil {
				log.Printf("failed to connect to %s: %s", pi.ID, err)
			}
			log.Printf("successfully connected to %s\n", pi.ID)
		}(pi)
	}
	wg.Wait()
	return nil
}

func AddCleanup(f func() error) {
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()
	cleanup = append(cleanup, f)
}

func DoCleanup() {
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()

	for _, f := range cleanup {
		if err := f(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
