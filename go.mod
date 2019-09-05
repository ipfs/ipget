module github.com/ipfs/ipget

go 1.12

require (
	github.com/ipfs/go-ipfs v0.4.22
	github.com/ipfs/go-ipfs-config v0.0.3
	github.com/ipfs/go-ipfs-files v0.0.4
	github.com/ipfs/go-ipfs-http-client v0.0.2
	github.com/ipfs/interface-go-ipfs-core v0.0.9
	github.com/libp2p/go-libp2p-peer v0.1.1
	github.com/libp2p/go-libp2p-peerstore v0.0.6
	github.com/multiformats/go-multiaddr v0.0.4
	github.com/urfave/cli v1.21.0
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

replace github.com/golangci/golangci-lint => github.com/mhutchinson/golangci-lint v1.17.2-0.20190819125825-d18f2136e32b

replace github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.3.5-0.20190904082202-d79a9f0c64db
