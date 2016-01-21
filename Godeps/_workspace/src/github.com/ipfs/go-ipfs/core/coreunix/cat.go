package coreunix

import (
	uio "github.com/ipfs/go-ipfs/unixfs/io"
	core "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/core"
	path "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/path"
	context "github.com/noffle/ipget/Godeps/_workspace/src/golang.org/x/net/context"
)

func Cat(ctx context.Context, n *core.IpfsNode, pstr string) (*uio.DagReader, error) {
	dagNode, err := core.Resolve(ctx, n, path.Path(pstr))
	if err != nil {
		return nil, err
	}
	return uio.NewDagReader(ctx, dagNode, n.DAG)
}
