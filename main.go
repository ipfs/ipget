package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	gopath "path"
	"path/filepath"
	"strings"
	"syscall"

	files "github.com/ipfs/go-ipfs-files"
	iface "github.com/ipfs/interface-go-ipfs-core"
	ipath "github.com/ipfs/interface-go-ipfs-core/path"
	cli "github.com/urfave/cli"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipget"
	app.Usage = "Retrieve and save IPFS objects."
	app.Version = "0.4.0"
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
		cli.StringSliceFlag{
			Name:  "peers,p",
			Usage: "specify a set of IPFS peers to connect to",
		},
		cli.BoolFlag{
			Name:  "progress",
			Usage: "show a progress bar",
		},
	}

	app.Action = func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if !c.Args().Present() {
			return fmt.Errorf("usage: ipget <ipfs ref>\n")
		}

		outPath := c.String("output")
		iPath, err := parsePath(c.Args().First())
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
			ipfs, err = http(ctx)
			if err == nil {
				break
			}
			fallthrough
		case "spawn":
			ipfs, err = spawn(ctx)
		case "local":
			ipfs, err = http(ctx)
		case "temp":
			ipfs, err = temp(ctx)
		default:
			return fmt.Errorf("no such 'node' strategy, %q", c.String("node"))
		}
		if err != nil {
			return err
		}

		go connect(ctx, ipfs, c.StringSlice("peers"))

		out, err := ipfs.Unixfs().Get(ctx, iPath)
		if err != nil {
			return cli.NewExitError(err, 2)
		}
		err = WriteTo(out, outPath, c.Bool("progress"))
		if err != nil {
			return cli.NewExitError(err, 2)
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

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
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
	ipfsPath := ipath.New(path)
	if ipfsPath.IsValid() == nil {
		return ipfsPath, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("%q could not be parsed: %s", path, err)
	}

	switch proto := u.Scheme; proto {
	case "ipfs", "ipld", "ipns":
		ipfsPath = ipath.New(gopath.Join("/", proto, u.Host, u.Path))
	case "http", "https":
		ipfsPath = ipath.New(u.Path)
	default:
		return nil, fmt.Errorf("%q is not recognized as an IPFS path", path)
	}
	return ipfsPath, ipfsPath.IsValid()
}

// WriteTo writes the given node to the local filesystem at fpath.
func WriteTo(nd files.Node, fpath string, progress bool) error {
	s, err := nd.Size()
	if err != nil {
		return err
	}

	var bar *pb.ProgressBar
	if progress {
		bar = pb.New64(s).Start()
	}

	return writeToRec(nd, fpath, bar)
}

func writeToRec(nd files.Node, fpath string, bar *pb.ProgressBar) error {
	switch nd := nd.(type) {
	case *files.Symlink:
		return os.Symlink(nd.Target, fpath)
	case files.File:
		f, err := os.Create(fpath)
		defer f.Close()
		if err != nil {
			return err
		}

		var r io.Reader = nd
		if bar != nil {
			r = bar.NewProxyReader(r)
		}
		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}
		return nil
	case files.Directory:
		err := os.Mkdir(fpath, 0777)
		if err != nil {
			return err
		}

		entries := nd.Entries()
		for entries.Next() {
			child := filepath.Join(fpath, entries.Name())
			if err := writeToRec(entries.Node(), child, bar); err != nil {
				return err
			}
		}
		return entries.Err()
	default:
		return fmt.Errorf("file type %T at %q is not supported", nd, fpath)
	}
}
