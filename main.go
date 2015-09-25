package main

import (
  "fmt"
  "io"
  "os"
  "github.com/jawher/mow.cli"
  core "github.com/ipfs/go-ipfs/core"
  context "golang.org/x/net/context"
	path "github.com/ipfs/go-ipfs/path"
	uio "github.com/ipfs/go-ipfs/unixfs/io"
)

const (
  usage = `Usage: ipget <ipfs-path>
      ipget <ipfs-path> -o <local-path>
      ipget <ipfs-path> (-a | --api) <addr>
      ipget <ipfs-path> [--mdns] [--solarnet] [--bootstrap <list>]
      ipget <ipfs-path> [ (-t | --transports) <protos> ]


Options:

-o, --output <local-path>   output to local path
-t, --transports <protos>   toggle transports (tcp, utp, udt, ble, aud)
-a, --api <addr>            delegate to node api at given address
--mdns                      enable mdns discovery (default: true)
--public                    use the public ipfs network for bootstrap (default: true)
--bootstrap <file>          use custom bootstrap list`
  defaultOutPath = "./aus"
)

func main() {
  cmd := cli.App("ipget", "Retrieve and save IPFS objects.")
  cmd.Spec = "IPFS_PATH [-o]"
  hash := cmd.String(cli.StringArg{
    Name: "IPFS_PATH",
    Value: "",
    Desc: "the IPFS object path",
  })
  outFile := cmd.StringOpt("o output", "", "output file path")
  cmd.Action = func() {
    if *hash == "" {
      fmt.Printf("you gotsta have a hash\n")
      return
    }

    if *outFile == "" {
      *outFile = "./" + *hash
    }

    get(*hash, *outFile)
  }
  cmd.Run(os.Args)
}

func get(path, outFile string) {
  node, err := core.NewNode(context.Background(), &core.BuildCfg{
    Online: true,
  })
  if err != nil {
    panic(err)
  }

  err = node.Bootstrap(core.DefaultBootstrapConfig)
  if err != nil {
    fmt.Printf("%v\n", err)
    return
  }

  readers, _, err := cat(node.Context(), node, []string{path})
  if err != nil {
    fmt.Printf("%v\n", err)
    return
  }

  file, err := os.Create(outFile)
  if err != nil {
    fmt.Printf("%v", err)
    return
  }

  reader := io.MultiReader(readers...)
  io.Copy(file, reader)

  fmt.Printf("wrote %v to %v\n", path, outFile)
}

func cat(ctx context.Context, node *core.IpfsNode, paths []string) ([]io.Reader, uint64, error) {
	readers := make([]io.Reader, 0, len(paths))
	length := uint64(0)
	for _, fpath := range paths {
		dagnode, err := core.Resolve(ctx, node, path.Path(fpath))
		if err != nil {
			return nil, 0, err
		}

		read, err := uio.NewDagReader(ctx, dagnode, node.DAG)
		if err != nil {
			return nil, 0, err
		}
		readers = append(readers, read)
		length += uint64(read.Size())
	}
	return readers, length, nil
}
