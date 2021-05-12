package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

type CfgOpt func(*config.Config)

func spawn(ctx context.Context) (iface.CoreAPI, error) {
	defaultPath, err := config.PathRoot()
	if err != nil {
		// shouldn't be possible
		return nil, err
	}

	if err := setupPlugins(defaultPath); err != nil {
		return nil, err
	}

	ipfs, err := open(ctx, defaultPath)
	if err == nil {
		return ipfs, nil
	}

	return tmpNode(ctx)
}

func setupPlugins(path string) error {
	// Load plugins. This will skip the repo if not available.
	plugins, err := loader.NewPluginLoader(filepath.Join(path, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func open(ctx context.Context, repoPath string) (iface.CoreAPI, error) {
	// Open the repo
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node
	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTClientOption,
		Repo:    r,
	})
	if err != nil {
		return nil, err
	}

	return coreapi.NewCoreAPI(node)
}

func temp(ctx context.Context) (iface.CoreAPI, error) {
	defaultPath, err := config.PathRoot()
	if err != nil {
		// shouldn't be possible
		return nil, err
	}

	if err := setupPlugins(defaultPath); err != nil {
		return nil, err
	}

	return tmpNode(ctx)
}

func tmpNode(ctx context.Context) (iface.CoreAPI, error) {
	dir, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return nil, fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Cleanup temp dir on exit
	addCleanup(func() error {
		return os.RemoveAll(dir)
	})

	identity, err := config.CreateIdentity(ioutil.Discard, []options.KeyGenerateOption{
		options.Key.Type(options.Ed25519Key),
	})
	if err != nil {
		return nil, err
	}
	cfg, err := config.InitWithIdentity(identity)
	if err != nil {
		return nil, err
	}

	// configure the temporary node
	cfg.Routing.Type = "dhtclient"

	err = fsrepo.Init(dir, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init ephemeral node: %s", err)
	}
	return open(ctx, dir)
}
