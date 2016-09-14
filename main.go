package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cli "github.com/codegangsta/cli"
	path "github.com/ipfs/go-ipfs/path"
	fallback "gx/ipfs/QmNpqXF6nNKesQgzCptL9GhMTw4QzWnid2evwT5TVU9zXC/fallback-ipfs-shell"
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

		shell, err := fallback.NewShell()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}

		if err := shell.Get(arg, outfile); err != nil {
			os.Remove(outfile)
			fmt.Fprintf(os.Stderr, "ipget failed: %s\n", err)
			os.Exit(2)
		}

		return nil
	}

	// Catch interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		os.Exit(1)
	}()

	app.Run(os.Args)
}
