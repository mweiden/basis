package main

import (
	"flag"
	"fmt"
	"github.com/mweiden/basis/problems/blockchain"
	"github.com/oklog/oklog/pkg/group"
	"io"
	"log"
	"net/http"
	"net"
	"context"
	"os/signal"
	"os"
	"syscall"
	"errors"
)

var (
	host = flag.String("host", "0.0.0.0", "Bind address for HTTP server")
	port = flag.Int("port", 5000, "Bind port for HTTP server")
)

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		return fmt.Errorf("received signal %s", sig)
	case <-cancel:
		return errors.New("canceled")
	}
}

// main
func main() {
	listen := fmt.Sprintf("%s:%d", *host, *port)

	api := blockchain.NewAPI(
		blockchain.NewBlockchain(blockchain.UnixTime),
		listen,
		func(url string) io.ReadCloser {
			response, _ := http.Get(url)
			return response.Body
		},
	)

	var g group.Group
	ctx, cancel := context.WithCancel(context.Background())

	// api
	g.Add(func () error {
		return api.Run(ctx)
	}, func (error) {
		cancel()
	})

	// http server
	mux := http.NewServeMux()
	mux.Handle("/", api)
	log.Printf("blockchain:web listening on %s", listen)
	ln, _ := net.Listen("tcp", listen)

	g.Add(func () error {
		return http.Serve(ln, api)
	}, func (error) {
		ln.Close()
	})

	// signal catcher
	cancelChan := make(chan struct{})

	g.Add(func () error {
		return interrupt(cancelChan)
	}, func (error) {
		close(cancelChan)
	})

	// run
	g.Run()
}
