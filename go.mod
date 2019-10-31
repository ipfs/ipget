module github.com/ipfs/ipget

go 1.12

require (
	github.com/ipfs/go-ipfs v0.4.22-0.20191030074653-cf150c1bf096
	github.com/ipfs/go-ipfs-addr v0.0.1 // indirect
	github.com/ipfs/go-ipfs-config v0.0.11
	github.com/ipfs/go-ipfs-files v0.0.4
	github.com/ipfs/go-ipfs-http-client v0.0.5
	github.com/ipfs/interface-go-ipfs-core v0.2.3
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.1.3
	github.com/multiformats/go-multiaddr v0.1.1
	github.com/urfave/cli v1.21.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

// go 1.13 stuff
replace github.com/golangci/golangci-lint => github.com/golangci/golangci-lint v1.18.0

replace github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.3.5-0.20190526074819-1df300866540
