package embeddedShell

import (
	"github.com/noffle/ipget/Godeps/_workspace/src/gopkg.in/errgo.v1"

	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/core"
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/path"
	tar "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/thirdparty/tar"
	uarchive "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/unixfs/archive"
)

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Get(ref, outdir string) error {
	ipfsPath, err := path.ParsePath(ref)
	if err != nil {
		return errgo.Notef(err, "get: could not parse %q", ref)
	}

	nd, err := core.Resolve(s.ctx, s.node, ipfsPath)
	if err != nil {
		return errgo.Notef(err, "get: could not resolve %s", ipfsPath)
	}

	r, err := uarchive.DagArchive(s.ctx, nd, outdir, s.node.DAG, false, 0)
	if err != nil {
		return err
	}

	ext := tar.Extractor{outdir}

	return ext.Extract(r)
}
