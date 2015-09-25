# Installation

```
$ go get github.com/jawher/mow.cli github.com/ipfs/go-ipfs/core github.com/ipfs/go-ipfs/path github.com/ipfs/go-ipfs/unixfs/io github.com/cheggaaa/pb

$ go build

$ go install

$ ipget QmTJHuzG3mjgmvcfvTU4ykWXwD4QjA5aCk6QsU4BPaD8Hh/cat2.gif -o nyan.gif
```

# Usage
```
Usage: ipget IPFS_PATH [-o]

Retrieve and save IPFS objects.

Arguments:
  IPFS_PATH=""   the IPFS object path

Options:
  -o, --output=""   output file path
```
