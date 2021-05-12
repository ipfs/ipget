package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	gopath "path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	files "github.com/ipfs/go-ipfs-files"
	iface "github.com/ipfs/interface-go-ipfs-core"
	ipath "github.com/ipfs/interface-go-ipfs-core/path"
	cli "github.com/urfave/cli/v2"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var (
	cleanup      []func() error
	cleanupMutex sync.Mutex
)

func main() {
	// Do any cleanup on exit
	defer doCleanup()

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
			if err == context.Canceled {
				return <-sigExitCoder
			}
			return cli.Exit(err, 2)
		}
		err = WriteTo(out, outPath, c.Bool("progress"))
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
		doCleanup()
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

func addCleanup(f func() error) {
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()
	cleanup = append(cleanup, f)
}

func doCleanup() {
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()

	for _, f := range cleanup {
		if err := f(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
