package shell

import (
	"fmt"
	"os"
	"os/signal"

	api "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs-api"
	fsrepo "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/repo/fsrepo"
	embedded "github.com/noffle/ipget/Godeps/_workspace/src/github.com/whyrusleeping/ipfs-embedded-shell"
	context "github.com/noffle/ipget/Godeps/_workspace/src/golang.org/x/net/context"
)

func NewShell() (Shell, error) {
	myShell, err := getApiShell()
	if err == nil {
		// fmt.Println("got an api shell!")
		return myShell, nil
	}

	myShell, err = getEmbeddedShell()
	if err == nil {
		// fmt.Println("got an embedded shell!")
		return myShell, nil
	}

	return nil, err
}

func getApiShell() (Shell, error) {
	apiShell := api.NewShell("http://127.0.0.1:5001")
	_, _, err := apiShell.Version()
	if err != nil {
		return nil, err
	}

	return apiShell, nil
}

func getEmbeddedShell() (Shell, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the ipfs node context if the process gets interrupted or killed.
	// TODO(noffle): is this needed?
	go func() {
		interrupts := make(chan os.Signal, 1)
		signal.Notify(interrupts, os.Interrupt, os.Kill)
		<-interrupts
		cancel()
	}()

	shell, err := tryLocal(ctx)
	if err == nil {
		return shell, nil
	}

	node, err := embedded.NewTmpDirNode(ctx)
	if err != nil {
		return nil, err
	}

	return embedded.NewShell(node), nil
}

func tryLocal(ctx context.Context) (Shell, error) {
	repoPath, err := getRepoPath()
	if err != nil {
		return nil, err
	}

	node, err := embedded.NewDefaultNodeWithFSRepo(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't get embedded shell: %s", err)
	}

	return embedded.NewShell(node), nil
}

func getRepoPath() (string, error) {
	repoPath, err := fsrepo.BestKnownPath()
	if err != nil {
		return "", err
	}
	return repoPath, nil
}
