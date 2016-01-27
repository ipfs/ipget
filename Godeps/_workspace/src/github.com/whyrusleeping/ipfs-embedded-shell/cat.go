package embeddedShell

import (
	"io"

	"github.com/noffle/ipget/Godeps/_workspace/src/gopkg.in/errgo.v1"

	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/core"
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/path"
	unixfsio "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/unixfs/io"
)

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Cat(p string) (io.ReadCloser, error) {
	ipfsPath, err := path.ParsePath(p)
	if err != nil {
		return nil, errgo.Notef(err, "cat: could not parse %q", p)
	}
	nd, err := core.Resolve(s.ctx, s.node, ipfsPath)
	if err != nil {
		return nil, errgo.Notef(err, "cat: could not resolve %s", ipfsPath)
	}
	dr, err := unixfsio.NewDagReader(s.ctx, nd, s.node.DAG)
	if err != nil {
		return nil, errgo.Notef(err, "cat: failed to construct DAG reader")
	}
	return dr, nil
}
