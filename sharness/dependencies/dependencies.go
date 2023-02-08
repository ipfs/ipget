//go:build tools

package tools

import (
	_ "github.com/chriscool/go-sleep"
	_ "github.com/ipfs/go-ipfs/cmd/ipfs"
	_ "github.com/whyrusleeping/pollEndpoint"

	// We depend on ipget, then use a ../../ replace directive to ensure we end up using the
	// _same_ version of go-ipfs. If we update any dependencies in the main module, `go mod
	// tidy` should fail to produce clean results in this module which should fail CI.
	_ "github.com/ipfs/ipget"
)
