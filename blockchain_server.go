package main

import (
	"flag"
	"fmt"
	"github.com/mweiden/basis/problems/blockchain"
	"log"
	"net/http"
	"io"
)

var (
	host = flag.String("host", "0.0.0.0", "Bind address for HTTP server")
	port = flag.Int("port", 5000, "Bind port for HTTP server")
)

// main
func main() {
	listen := fmt.Sprintf("%s:%d", *host, *port)

	api := blockchain.NewAPI(
		blockchain.NewBlockchain(blockchain.UnixTime),
		listen,
		func (url string) io.ReadCloser {
			response, _ := http.Get(url)
			return response.Body
		},
	)

	mux := http.NewServeMux()
	mux.Handle("/", api)

	log.Printf("blockchain:web listening on %s", listen)
	log.Print(http.ListenAndServe(listen, mux))
}
