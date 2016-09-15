package main

import (
	"fmt"
	"os"

	cli "github.com/codegangsta/cli"
	path "github.com/ipfs/go-ipfs/path"
	fallback "gx/ipfs/QmXEk5yvscBkRAMLpHiCVkUoLrXYBkFrbNaiWprocDKS7Z/fallback-ipfs-shell"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipget"
	app.Usage = "Retrieve and save IPFS objects."
	app.Version = "0.2.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output,o",
			Usage: "specify output location",
		},
		cli.StringFlag{
			Name:  "node,n",
			Usage: "specify ipfs node strategy ('local', 'spawn', or 'fallback')",
			Value: "fallback",
		},
	}

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			fmt.Fprintf(os.Stderr, "usage: ipget <ipfs ref>\n")
			os.Exit(1)
		}

		outfile := c.String("output")
		arg := c.Args().First()

		// Use the final segment of the object's path if no path was given.
		if outfile == "" {
			ipfsPath, err := path.ParsePath(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ParsePath failure: %s\n", err)
				os.Exit(1)
			}
			segments := ipfsPath.Segments()
			outfile = segments[len(segments)-1]
		}

		var shell fallback.Shell
		var err error

		if c.String("node") == "fallback" {
			shell, err = fallback.NewShell()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
		} else if c.String("node") == "spawn" {
			shell, err = fallback.NewEmbeddedShell()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
		} else if c.String("node") == "local" {
			shell, err = fallback.NewApiShell()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "error: no such 'node' strategy, '%s'\n", c.String("node"))
			os.Exit(1)
			return nil
		}

		if err := shell.Get(arg, outfile); err != nil {
			os.Remove(outfile)
			fmt.Fprintf(os.Stderr, "ipget failed: %s\n", err)
			os.Exit(2)
		}

		return nil
	}

	app.Run(os.Args)
}
