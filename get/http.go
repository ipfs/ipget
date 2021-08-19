package get

import (
	"context"

	ipfshttp "github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

func Http(ctx context.Context) (iface.CoreAPI, error) {
	httpApi, err := ipfshttp.NewLocalApi()
	if err != nil {
		return nil, err
	}
	err = httpApi.Request("version").Exec(ctx, nil)
	if err != nil {
		return nil, err
	}
	return httpApi, nil
}
