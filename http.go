package main

import (
	"context"

	iface "github.com/ipfs/boxo/coreiface"
	ipfshttp "github.com/ipfs/kubo/client/rpc"
)

func http(ctx context.Context) (iface.CoreAPI, error) {
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
