package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	gopath "path"
	"path/filepath"
	"strings"
	"syscall"

	ipath "gx/ipfs/QmT3rzed1ppXefourpmoZ7tyVQfsGPQZ1pHDngLmCvXxd3/go-path"
	fallback "gx/ipfs/QmaWDhoQaV6cDyy6NSKFgPaUAGRtb4SMiLpaDYEsxP7X8P/fallback-ipfs-shell"
	cli "gx/ipfs/Qmc1AtgBdoUHP8oYSqU81NRYdzohmF45t5XNwVMvhCxsBA/cli"
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

		outPath := c.String("output")
		iPath, err := parsePath(c.Args().First())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}

		// Use the final segment of the object's path if no path was given.
		if outPath == "" {
			trimmed := strings.TrimRight(iPath.String(), "/")
			_, outPath = filepath.Split(trimmed)
			outPath = filepath.Clean(outPath)
		}

		var shell fallback.Shell

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

		if err := shell.Get(iPath.String(), outPath); err != nil {
			os.Remove(outPath)
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

func parsePath(path string) (ipath.Path, error) {
	ipfsPath, err := ipath.ParsePath(path)
	if err == nil { // valid canonical path
		return ipfsPath, nil
	}
	u, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("%q could not be parsed: %s", path, err)
	}

	switch proto := u.Scheme; proto {
	case "ipfs", "ipld", "ipns":
		return ipath.ParsePath(gopath.Join("/", proto, u.Host, u.Path))
	case "http", "https":
		return ipath.ParsePath(u.Path)
	default:
		return "", fmt.Errorf("%q is not recognized as an IPFS path")
	}
}
