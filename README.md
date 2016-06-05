[![Build Status](https://secure.travis-ci.org/ipfs/ipget.png?branch=master)](http://travis-ci.org/ipfs/ipget)

# Why ipget?

The IPFS CLI offers an `ipfs get` command that is very similar to `ipget`. Why bother with this utility? Well,

1. `ipget` has no dependencies: it's just a standalone binary. This makes it
   great for including in distributions.
2. `ipget` doesn't require an IPFS daemon to be running. (Though it will use it if available.)

This makes it ideal for users and projects that want a simple dependency-free utility that will Just Work(tm) whenever they want to retrieve files from IPFS with minimal fuss.


# Install

Grab a binary for your platform at https://dist.ipfs.io/#ipget


# Install From Source

```
$ go get -d github.com/ipfs/ipget

$ cd ${GOPATH}/src/github.com/ipfs/ipget

$ make install
```


# Example

Find a fun IPFS address and `ipget` away!

```
$ ipget QmTJHuzG3mjgmvcfvTU4ykWXwD4QjA5aCk6QsU4BPaD8Hh/cat2.gif -o nyan.gif
```
or with an `/ipfs` prefix:
```
$ ipget /ipfs/QmTJHuzG3mjgmvcfvTU4ykWXwD4QjA5aCk6QsU4BPaD8Hh/cat2.gif -o nyan.gif
```
or even IPNS addresses:
```
$ ipget /ipns/QmQG1kwx91YQsGcsa9Z1p6BPJ3amdiSLLmsmAoEMwbX61b/files/cat2.gif -o nyan.gif
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
