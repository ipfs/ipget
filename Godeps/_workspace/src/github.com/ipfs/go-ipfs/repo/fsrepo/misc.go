package fsrepo

import (
	"os"

	"github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/repo/config"
	homedir "github.com/noffle/ipget/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
)

// BestKnownPath returns the best known fsrepo path. If the ENV override is
// present, this function returns that value. Otherwise, it returns the default
// repo path.
func BestKnownPath() (string, error) {
	ipfsPath := config.DefaultPathRoot
	if os.Getenv(config.EnvDir) != "" {
		ipfsPath = os.Getenv(config.EnvDir)
	}
	ipfsPath, err := homedir.Expand(ipfsPath)
	if err != nil {
		return "", err
	}
	return ipfsPath, nil
}
