# Installation

```
$ go get github.com/noffle/fallback-ipfs-shell
```

# Usage
```
import (
  shell "github.com/noffle/fallback-ipfs-shell/shell"
)
```

### shell.NewShell() (shell.Shell, error)

Returns a Shell interface, preferring a local HTTP API node if it can find one,
but falling back to producing a new ephemeral node that self-bootstraps.

