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
  pb "github.com/cheggaaa/pb"
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

  reader, length, err := cat(node.Context(), node, path)
  if err != nil {
    fmt.Printf("%v\n", err)
    return
  }

  file, err := os.Create(outFile)
  if err != nil {
    fmt.Printf("%v", err)
    return
  }

  bar := pb.New(int(length)).SetUnits(pb.U_BYTES)
  bar.ShowSpeed = false
  bar.Start()
  writer := io.MultiWriter(file, bar)

  io.Copy(writer, reader)

  bar.Finish()

  fmt.Printf("wrote %v to %v (%v bytes)\n", path, outFile, length)
}

func cat(ctx context.Context, node *core.IpfsNode, fpath string) (io.Reader, uint64, error) {
  dagnode, err := core.Resolve(ctx, node, path.Path(fpath))
  if err != nil {
    return nil, 0, err
  }

  reader, err := uio.NewDagReader(ctx, dagnode, node.DAG)
  if err != nil {
    return nil, 0, err
  }
  length := uint64(reader.Size())

	return reader, length, nil
}
