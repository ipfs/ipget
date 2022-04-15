# ipget

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](https://ipfs.io/)

> wget for IPFS: retrieve files over IPFS and save them locally.

`ipget` is a standalone program analogous to [GNU `wget`](https://www.gnu.org/software/wget/). Unlike wget though, `ipget` specializes in downloading files and directory structures from the [IPFS network](https://ipfs.io).

`ipget` includes its own IPFS node, so you don't need IPFS installed on your
system. This makes it ideal for users and projects that want a simple utility
for whenever they want to retrieve files from IPFS.


## Install

Download a binary for your platform from [IPFS Distributions](https://dist.ipfs.io/#ipget).

### Install From Source

```
$ go get -d github.com/ipfs/ipget

$ cd ${GOPATH}/src/github.com/ipfs/ipget

$ make install
```

### Example

Find a fun IPFS address and `ipget` away!

```
$ ipget QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif -o nyan.gif
```
or with an `/ipfs` prefix:
```
$ ipget -o nyan.gif /ipfs/QmQ2r6iMNpky5f1m4cnm3Yqw8VSvjuKpTcK1X7dBR1LkJF/cat.gif
```
or even IPNS addresses (note this is an IPNS address example and does not work):
```
$ ipget /ipns/QmQG1kwx91YQsGcsa9Z1p6BPJ3amdiSLLmsmAoEMwbX61b/files/cat.gif
```

## Usage

```
NAME:
   ipget - Retrieve and save IPFS objects.

USAGE:
   ipget [global options] command [command options] [arguments...]

VERSION:
   0.8.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output value, -o value  specify output location
   --node value, -n value    specify ipfs node strategy ("local", "spawn", "temp" or "fallback") (default: "fallback")
   --peers value, -p value   specify a set of IPFS peers to connect to
   --progress                show a progress bar (default: false)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/ipfs/ipget/issues)!

This repository falls under the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

[![](https://cdn.rawgit.com/jbenet/contribute-ipfs-gif/master/img/contribute.gif)](https://github.com/ipfs/community/blob/master/CONTRIBUTING.md)

## License

[MIT](LICENSE)
