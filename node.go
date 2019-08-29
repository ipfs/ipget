package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/interface-go-ipfs-core"
)

func spawn(ctx context.Context) (iface.CoreAPI, error) {
	defaultPath, err := config.PathRoot()
	if err != nil {
		// shouldn't be possible
		return nil, err
	}

	// Load plugins. This will skip the repo if not available.
	plugins, err := loader.NewPluginLoader(filepath.Join(defaultPath, "plugins"))
	if err != nil {
		return nil, fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return nil, fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return nil, fmt.Errorf("error initializing plugins: %s", err)
	}

	ipfs, err := open(ctx, defaultPath)
	if err == nil {
		return ipfs, nil
	}

	return tmpNode(ctx)
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

func tmpNode(ctx context.Context) (iface.CoreAPI, error) {
	dir, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return nil, fmt.Errorf("failed to get temp dir: %s", err)
	}

	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return nil, err
	}

	// Set some ipget specific defaults.
	cfg.Experimental.ShardingEnabled = true
	cfg.Experimental.QUIC = true
	cfg.Routing.Type = "dhtclient"
	cfg.Reprovider.Interval = "0"

	err = fsrepo.Init(dir, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init ephemeral node: %s", err)
	}
	return open(ctx, dir)
}
