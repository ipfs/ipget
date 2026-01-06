package main

import (
	"context"
	"fmt"
	"os"

	ipfshttp "github.com/ipfs/kubo/client/rpc"
	iface "github.com/ipfs/kubo/core/coreiface"
)

func http(ctx context.Context) (iface.CoreAPI, error) {
	fmt.Fprint(os.Stderr, "Downloading from local daemon... ")

	httpAPI, err := ipfshttp.NewLocalApi()
	if err != nil {
		return nil, err
	}
	err = httpAPI.Request("version").Exec(ctx, nil)
	if err != nil {
		return nil, err
	}
	return httpAPI, nil
}
