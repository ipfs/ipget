package embeddedShell

import (
	"io"

	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/core"
	"github.com/noffle/ipget/Godeps/_workspace/src/golang.org/x/net/context"
)

// Interface ...
type Interface interface {
	// Add reads everything from the Reader and returns the IPFS hash
	Add(io.Reader) (string, error)

	// Cat returns a reader returning the data under the IPFS path
	Cat(path string) (io.ReadCloser, error)

	// AddDir(dir string) (string, error)
	Get(hash string, outdir string) error

	// List(path string) ([]*LsLink, error)
	//FileList(path string) (*UnixLsObject, error)
	//FindPeer(peer string) (*PeerInfo, error)

	//ID(peer ...string) (*IdOutput, error)
	// Version() (string, string, error)

	// Refs(hash string, recursive bool) (<-chan string, error)

	// NewObject(template string) (string, error)
	// AddLink(target string) (string, error)
	// Patch(root string, action string, args ...string) (string, error)
	// PatchLink(root string, path string, childhash string, create bool) (string, error)

	// Pin(path string) error
	// Unpin(path string) error

	// Publish(node string, value string) error
	// Resolve(id string) (string, error)
	// ResolvePath(path string) (string, error)
}

// Shell ...
type Shell struct {
	ctx  context.Context
	node *core.IpfsNode
}

// func NewReadOnlyShell() *Shell {}

func NewShell(node *core.IpfsNode) *Shell {
	return NewShellWithContext(node, context.Background())
}

func NewShellWithContext(node *core.IpfsNode, ctx context.Context) *Shell {
	return &Shell{ctx, node}
}
