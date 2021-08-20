package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"	
	"strings"
	"syscall"

	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/ipget/get"
	cli "github.com/urfave/cli/v2"
)

func main() {
	// Do any cleanup on exit
	defer get.DoCleanup()

	app := cli.NewApp()
	app.Name = "ipget"
	app.Usage = "Retrieve and save IPFS objects."
	app.Version = "0.7.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "specify output location",
		},
		&cli.StringFlag{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "specify ipfs node strategy (\"local\", \"spawn\", \"temp\" or \"fallback\")",
			Value:   "fallback",
		},
		&cli.StringSliceFlag{
			Name:    "peers",
			Aliases: []string{"p"},
			Usage:   "specify a set of IPFS peers to connect to",
		},
		&cli.BoolFlag{
			Name:  "progress",
			Usage: "show a progress bar",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigExitCoder := make(chan cli.ExitCoder, 1)

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			return fmt.Errorf("usage: ipget <ipfs ref>\n")
		}

		outPath := c.String("output")
		iPath, err := get.ParsePath(c.Args().First())
		if err != nil {
			return err
		}

		// Use the final segment of the object's path if no path was given.
		if outPath == "" {
			trimmed := strings.TrimRight(iPath.String(), "/")
			_, outPath = filepath.Split(trimmed)
			outPath = filepath.Clean(outPath)
		}

		var ipfs iface.CoreAPI
		switch c.String("node") {
		case "fallback":
			ipfs, err = get.Http(ctx)
			if err == nil {
				break
			}
			fallthrough
		case "spawn":
			ipfs, err = get.Spawn(ctx)
		case "local":
			ipfs, err = get.Http(ctx)
		case "temp":
			ipfs, err = get.Temp(ctx)
		default:
			return fmt.Errorf("no such 'node' strategy, %q", c.String("node"))
		}
		if err != nil {
			return err
		}

		go get.Connect(ctx, ipfs, c.StringSlice("peers"))

		out, err := ipfs.Unixfs().Get(ctx, iPath)
		if err != nil {
			if err == context.Canceled {
				return <-sigExitCoder
			}
			return cli.Exit(err, 2)
		}
		err = get.WriteTo(out, outPath, c.Bool("progress"))
		if err != nil {
			if err == context.Canceled {
				return <-sigExitCoder
			}
			return cli.Exit(err, 2)
		}
		return nil
	}

	// Catch interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		sigExitCoder <- cli.Exit("", 128+int(sig.(syscall.Signal)))
		cancel()
	}()

	// cli library requires flags before arguments
	args := movePostfixOptions(os.Args)

	err := app.Run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		get.DoCleanup()
		os.Exit(1)
	}
}

// movePostfixOptions moves non-flag arguments to end of argument list.
func movePostfixOptions(args []string) []string {
	var endArgs []string
	for idx := 1; idx < len(args); idx++ {
		if args[idx][0] == '-' {
			if !strings.Contains(args[idx], "=") {
				idx++
			}
			continue
		}
		if endArgs == nil {
			// on first write, make copy of args
			newArgs := make([]string, len(args))
			copy(newArgs, args)
			args = newArgs
		}
		// add to args accumulator
		endArgs = append(endArgs, args[idx])
		// remove from real args list
		args = args[:idx+copy(args[idx:], args[idx+1:])]
		idx--
	}

	// append extracted arguments to the real args
	return append(args, endArgs...)
}
