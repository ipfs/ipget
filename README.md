# ipget

[![Build Status](https://secure.travis-ci.org/ipfs/ipget.png?branch=master)](http://travis-ci.org/ipfs/ipget)


`ipget` is a standalone program analogous to [GNU `wget`](https://www.gnu.org/software/wget/). Unlike wget though, `ipget` specializes in downloading files and directory structures from the [IPFS network](https://ipfs.io).

`ipget` includes its own IPFS node, so you don't need IPFS installed on your
system. This makes it ideal for users and projects that want a simple utility
for whenever they want to retrieve files from IPFS.


## Install

Download a binary for your platform from [IPFS Distributions](https://dist.ipfs.io/#ipget).


## Install From Source

`ipget` doesn't use the vanilla Go package management system. It instead uses
the [gx](https://github.com/whyrusleeping/gx) (and
[gx-go](https://github.com/whyrusleeping/gx-go)) workflow. This means a slightly
different set of steps to install:

```
$ go get -d github.com/ipfs/ipget

$ cd ${GOPATH}/src/github.com/ipfs/ipget

$ make install
```


## Example

Find a fun IPFS address and `ipget` away!

```
$ ipget QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif -o nyan.gif
```
or with an `/ipfs` prefix:
```
$ ipget /ipfs/QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif -o nyan.gif
```
or even IPNS addresses:
```
$ ipget /ipns/QmQG1kwx91YQsGcsa9Z1p6BPJ3amdiSLLmsmAoEMwbX61b/files/cat.gif -o nyan.gif
```


## Usage
```
NAME:
   ipget - Retrieve and save IPFS objects.

USAGE:
   ipget [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
GLOBAL OPTIONS:
   --output value, -o value  specify output location
   --help, -h                show help
   --version, -v             print the version
```


## License

MIT
