[![Build Status](https://secure.travis-ci.org/noffle/ipget.png?branch=master)](http://travis-ci.org/noffle/ipget)

# Why ipget?

The IPFS CLI offers an `ipfs get` command that is very similar to `ipget`. Why bother with this utility? Well,

1. `ipget` has no dependencies: it's just a simple static binary.
2. `ipget` doesn't require an IPFS daemon to be running. (Though it will use it if available.)

This makes it ideal for users and projects that want a simple dependency-free utility that will Just Work(tm) whenever they want to retrieve files from IPFS with minimal fuss.

# Installation

```
$ go get github.com/noffle/ipget

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
