package main

import (
	"fmt"
	"os"

	path "github.com/ipfs/go-ipfs/path"
	cli "github.com/jawher/mow.cli"
	fallback "github.com/noffle/fallback-ipfs-shell"
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

		shell, err := fallback.NewShell()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}

		if err = shell.Get(*hash, *outFile); err != nil {
			os.Remove(*outFile)
			fmt.Fprintf(os.Stderr, "ipget failed: %s\n", err)
			os.Exit(2)
		}
	}
	cmd.Run(os.Args)
}
