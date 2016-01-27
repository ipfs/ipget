package embeddedShell

import (
	"fmt"

	core "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/core"
	dag "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/merkledag"
	dagutils "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/merkledag/utils"
	path "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/path"
	ft "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/unixfs"
)

func (s *Shell) NewObject(template string) (string, error) {
	node := new(dag.Node)
	switch template {
	case "":
		break
	case "unixfs-dir":
		node.Data = ft.FolderPBData()
	default:
		return "", fmt.Errorf("unknown template %s", template)
	}
	k, err := s.node.DAG.Add(node)
	if err != nil {
		return "", err
	}

	return k.B58String(), nil
}

// TODO: extract all this logic from the core/commands/object.go to avoid dupe code
func (s *Shell) Patch(root, action string, args ...string) (string, error) {
	p, err := path.ParsePath(root)
	if err != nil {
		return "", err
	}

	rootnd, err := core.Resolve(s.ctx, s.node, p)
	if err != nil {
		return "", err
	}

	insertpath := args[0]
	childhash := args[1]

	childpath, err := path.ParsePath(childhash)
	if err != nil {
		return "", err
	}

	nnode, err := core.Resolve(s.ctx, s.node, childpath)
	if err != nil {
		return "", err
	}

	e := dagutils.NewDagEditor(s.node.DAG, rootnd)

	switch action {
	case "add-link":
		err := e.InsertNodeAtPath(s.ctx, insertpath, nnode, nil)
		if err != nil {
			return "", err
		}

		err = e.WriteOutputTo(s.node.DAG)
		if err != nil {
			return "", err
		}

		final, err := e.GetNode().Key()
		if err != nil {
			return "", err
		}

		return final.B58String(), nil
	default:
		return "", fmt.Errorf("unsupported action (impl not complete)")
	}
}

//TODO: hrm, maybe this interface could be better
func (s *Shell) PatchLink(root, npath, childhash string, create bool) (string, error) {
	p, err := path.ParsePath(root)
	if err != nil {
		return "", err
	}

	rootnd, err := core.Resolve(s.ctx, s.node, p)
	if err != nil {
		return "", err
	}

	childpath, err := path.ParsePath(childhash)
	if err != nil {
		return "", err
	}

	nnode, err := core.Resolve(s.ctx, s.node, childpath)
	if err != nil {
		return "", err
	}

	e := dagutils.NewDagEditor(s.node.DAG, rootnd)
	err = e.InsertNodeAtPath(s.ctx, npath, nnode, func() *dag.Node {
		return &dag.Node{Data: ft.FolderPBData()}
	})
	if err != nil {
		return "", err
	}

	err = e.WriteOutputTo(s.node.DAG)
	if err != nil {
		return "", err
	}

	final, err := e.GetNode().Key()
	if err != nil {
		return "", err
	}

	return final.B58String(), nil
}
