package embeddedShell

import (
	"io"

	"github.com/noffle/ipget/Godeps/_workspace/src/gopkg.in/errgo.v1"

	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/importer"
	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/importer/chunk"
	dag "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/merkledag"
	ft "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/unixfs"
)

func (s *Shell) Add(r io.Reader) (string, error) {
	dag, err := importer.BuildDagFromReader(
		s.node.DAG,
		chunk.DefaultSplitter(r),
		importer.BasicPinnerCB(s.node.Pinning.GetManual()), // TODO: make pinning configurable
	)
	if err != nil {
		return "", errgo.Notef(err, "add: importing DAG failed.")
	}
	k, err := dag.Key()
	if err != nil {
		return "", errgo.Notef(err, "add: getting key from DAG failed.")
	}
	return k.B58String(), nil
}

// AddLink creates a unixfs symlink and returns its hash
func (s *Shell) AddLink(target string) (string, error) {
	d, _ := ft.SymlinkData(target)
	nd := &dag.Node{Data: d}
	k, err := s.node.DAG.Add(nd)
	if err != nil {
		return "", err
	}

	return k.B58String(), nil
}
