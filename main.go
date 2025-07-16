package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	gopath "path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/cheggaaa/pb/v3"
	files "github.com/ipfs/boxo/files"
	ipath "github.com/ipfs/boxo/path"
	iface "github.com/ipfs/kubo/core/coreiface"
	cli "github.com/urfave/cli/v3"
)

var (
	cleanup      []func() error
	cleanupMutex sync.Mutex
)

func main() {
	// Do any cleanup on exit
	defer doCleanup()

	cmd := &cli.Command{
		Name:      "ipget",
		Usage:     "Retrieve and save IPFS objects.",
		Version:   version,
		UsageText: "ipget [options] ipfs_object [ipfs_object ...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "specify output location",
			},
			&cli.StringFlag{
				Name:    "node",
				Aliases: []string{"n"},
				Usage: "specify ipfs node strategy (\"local\", \"spawn\", \"temp\" or \"fallback\")" +
					"\nlocal    connect to a local IPFS daemon" +
					"\nspawn    run ipget as an IPFS node using an existing repo, use 'temp' strategy if no repo" +
					"\ntemp     run ipget as an IPFS node using a temporary repo that is removed on command completion" +
					"\nfallback tries 'local' strategy first and then 'spawn' if no local daemon is available",
				Value: "fallback",
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
		},
		Action: ipgetAction,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "Received interrupt signal, shutting down...")
		cancel()
	}()

	// cli library requires flags before arguments
	err := cmd.Run(ctx, movePostfixOptions(os.Args))
	doCleanup()

	if err != nil {
		if errors.Is(err, context.Canceled) {
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ipgetAction(ctx context.Context, cmd *cli.Command) error {
	if !cmd.Args().Present() {
		return fmt.Errorf("usage: ipget <ipfs ref>")
	}

	outPath := cmd.String("output")
	iPath, err := parsePath(cmd.Args().First())
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
	switch cmd.String("node") {
	case "fallback":
		ipfs, err = http(ctx)
		if err != nil {
			ipfs, err = spawn(ctx)
		}
	case "spawn":
		ipfs, err = spawn(ctx)
	case "local":
		ipfs, err = http(ctx)
	case "temp":
		ipfs, err = temp(ctx)
	default:
		return fmt.Errorf("no such 'node' strategy, %q", cmd.String("node"))
	}
	if err != nil {
		return err
	}

	go connect(ctx, ipfs, cmd.StringSlice("peers"))

	out, err := ipfs.Unixfs().Get(ctx, iPath)
	if err != nil {
		return err
	}
	if err = WriteTo(out, outPath, cmd.Bool("progress")); err != nil {
		return err
	}
	return nil
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
	ipfsPath, err := ipath.NewPath(path)
	if err == nil {
		return ipfsPath, nil
	}
	origErr := err

	ipfsPath, err = ipath.NewPath("/ipfs/" + path)
	if err == nil {
		return ipfsPath, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, origErr
	}
	switch u.Scheme {
	case "ipfs", "ipld", "ipns":
		return ipath.NewPath(gopath.Join("/", u.Scheme, u.Host, u.Path))
	case "http", "https":
		return ipath.NewPath(u.Path)
	}
	return nil, fmt.Errorf("%q is not recognized as an IPFS path", path)
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
		err := os.Symlink(nd.Target, fpath)
		if err != nil {
			return err
		}
		switch runtime.GOOS {
		case "linux", "freebsd", "netbsd", "openbsd", "dragonfly":
			return files.UpdateModTime(fpath, nd.ModTime())
		default:
			return nil
		}
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
		return files.UpdateMeta(fpath, nd.Mode(), nd.ModTime())
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

		if err = files.UpdateMeta(fpath, nd.Mode(), nd.ModTime()); err != nil {
			return err
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
