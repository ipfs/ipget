package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	pb "github.com/noffle/ipget/Godeps/_workspace/src/github.com/cheggaaa/pb"
	humanize "github.com/noffle/ipget/Godeps/_workspace/src/github.com/dustin/go-humanize"
	core "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/core"
	path "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/path"
	uarchive "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/unixfs/archive"
	cli "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jawher/mow.cli"
	tar "github.com/noffle/ipget/Godeps/_workspace/src/github.com/whyrusleeping/tar-utils"
	context "github.com/noffle/ipget/Godeps/_workspace/src/golang.org/x/net/context"
)

func main() {
	cmd := cli.App("ipget", "Retrieve and save IPFS objects.")
	cmd.Spec = "IPFS_PATH [-o]"

	hash := cmd.String(cli.StringArg{
		Name:  "IPFS_PATH",
		Value: "",
		Desc:  "the IPFS object path",
	})

	outFile := cmd.StringOpt("o output", "", "output file path")

	cmd.Action = func() {
		// Use the final segment of the object's path if no path was given.
		if *outFile == "" {
			ipfsPath, err := path.ParsePath(*hash)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ParsePath failure: %s", err)
				os.Exit(1)
			}
			segments := ipfsPath.Segments()
			*outFile = segments[len(segments)-1]
		}

		if err := get(*hash, *outFile); err != nil {
			os.Remove(*outFile)
			fmt.Fprintf(os.Stderr, "ipget failed: %s\n", err)
			os.Exit(2)
		}
	}
	cmd.Run(os.Args)
}

func get(path, outFile string) error {
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the ipfs node context if the process gets interrupted or killed.
	go func() {
		interrupts := make(chan os.Signal, 1)
		signal.Notify(interrupts, os.Interrupt, os.Kill)
		<-interrupts
		cancel()
	}()

	// Create the ipfs node.
	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
	})
	if err != nil {
		return fmt.Errorf("ipfs NewNode() failed: %s", err)
	}

	// Bootstrap.
	err = node.Bootstrap(core.DefaultBootstrapConfig)
	if err != nil {
		return fmt.Errorf("node Bootstrap() failed: %s", err)
	}
	fmt.Fprintf(os.Stderr, "IPFS Node bootstrapping (took %v)\n", time.Since(start))

	// Get a reader for downloading the object content.
	reader, length, err := cat(node.Context(), node, path)
	if err != nil {
		return fmt.Errorf("cat failed: %s", err)
	}

	bar := pb.New(int(length)).SetUnits(pb.U_BYTES)
	bar.Output = os.Stderr
	bar.ShowSpeed = false
	bar.Start()

	// Stream the file content from the IPFS node to the output file.
	extractor := &tar.Extractor{outFile}
	if err := extractor.Extract(reader); err != nil {
		return fmt.Errorf("copy failed: %s", err)
	}

	// Done!
	bar.Finish()
	fmt.Fprintf(os.Stderr, "Wrote %q to %q (%s) (took %v)\n", path, outFile,
		humanize.Bytes(length), time.Since(start))

	return nil
}

func cat(ctx context.Context, node *core.IpfsNode, fpath string) (io.Reader, uint64, error) {
	dagnode, err := core.Resolve(ctx, node, path.Path(fpath))
	if err != nil {
		return nil, 0, err
	}

	reader, err := uarchive.DagArchive(node.Context(), dagnode, fpath, node.DAG, false, 0)
	if err != nil {
		return nil, 0, err
	}

	// TODO: modify uarchive.DagArchive to output pb size
	length, err := dagnode.Size()
	if err != nil {
		return nil, 0, err
	}
	return reader, length, nil
}
