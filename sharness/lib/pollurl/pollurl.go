// pollurl is a helper utility that waits for a http endpoint to be reachable
// and return with http.StatusOK.
//
// Run pollurl as a cli script:
//
//	go run pollurl.go [args...]
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

func main() {
	var (
		host     string
		endpoint = flag.String("ep", "/version", "which http endpoint path to hit")
		tries    = flag.Int("tries", 10, "how many tries to make before failing")
		timeout  = flag.Duration("tout", time.Second, "how long to wait between attempts")
		verbose  = flag.Bool("v", false, "verbose output")
	)
	flag.StringVar(&host, "host", "/ip4/127.0.0.1/tcp/5001", "the multiaddr host to dial on")
	flag.Parse()

	addr, err := ma.NewMultiaddr(host)
	if err != nil {
		log.Fatalln("NewMultiaddr() failed:", err)
	}
	p := addr.Protocols()
	if len(p) < 2 {
		log.Fatalln("need two protocols in host flag (/ip/tcp):", addr)
	}
	_, host, err = manet.DialArgs(addr)
	if err != nil {
		log.Fatalln("manet.DialArgs() failed:", err)
	}

	// construct url to dial
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   *endpoint,
	}
	targetURL := u.String()

	start := time.Now()
	if *verbose {
		log.Printf("starting at %s, tries: %d, timeout: %s, url: %s", start, *tries, *timeout, targetURL)
	}

	for *tries > 0 {
		err = checkOK(targetURL)
		if err != nil {
			if *verbose {
				log.Println("get failed:", err)
			}
			time.Sleep(*timeout)
			*tries--
		} else {
			if *verbose {
				log.Printf("ok - endpoint reachable with %d tries remaining, took %s", *tries, time.Since(start))
			}
			os.Exit(0)
		}
	}

	log.Print("failed")
	os.Exit(1)
}

func checkOK(endpoint string) error {
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pollurl: error reading response body: %s", err)
		}
		return fmt.Errorf("response not OK. %d %s %q", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}
