package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	path "gx/ipfs/QmT3rzed1ppXefourpmoZ7tyVQfsGPQZ1pHDngLmCvXxd3/go-path"
	fallback "gx/ipfs/QmaWDhoQaV6cDyy6NSKFgPaUAGRtb4SMiLpaDYEsxP7X8P/fallback-ipfs-shell"
	cli "gx/ipfs/Qmc1AtgBdoUHP8oYSqU81NRYdzohmF45t5XNwVMvhCxsBA/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipget"
	app.Usage = "Retrieve and save IPFS objects."
	app.Version = "0.3.2"
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

	// Catch interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		os.Exit(1)
	}()

	// TODO(noffle): remove this once https://github.com/urfave/cli/issues/427 is
	// fixed.
	args := movePostfixOptions(os.Args)

	app.Run(args)
}

// movePostfixOptions finds the Qmfoobar hash argument and moves it to the end
// of the argument array.
func movePostfixOptions(args []string) []string {
	var idx = 1
	the_args := make([]string, 0)
	for {
		if idx >= len(args) {
			break
		}

		if args[idx][0] == '-' {
			if !strings.Contains(args[idx], "=") {
				idx++
			}
		} else {
			// add to args accumulator
			the_args = append(the_args, args[idx])

			// remove from real args list
			new_args := make([]string, 0)
			new_args = append(new_args, args[:idx]...)
			new_args = append(new_args, args[idx+1:]...)
			args = new_args
			idx--
		}

		idx++
	}

	// append extracted arguments to the real args
	return append(args, the_args...)
}
