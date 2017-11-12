package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/mweiden/basis/problems/blockchain"
	"github.com/oklog/oklog/pkg/group"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	host = flag.String("host", "localhost", "Bind address for HTTP server")
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
	log.SetOutput(os.Stdout)
	flag.Parse()
	listen := fmt.Sprintf("%s:%d", *host, *port)
	log.SetPrefix(fmt.Sprintf("[%s] ", listen))
	rand.Seed(int64(*port))

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

	cronBeat := make(chan bool)
	go func() {
		for {
			// TODO: This is kind of a hack to fake different mining completions
			// between instances of mattcoin
			sleepTime := (10 + time.Duration(rand.Int()%5)) * time.Second
			time.Sleep(sleepTime)
			cronBeat <- true
		}
	}()

	// api
	g.Add(
		func() error { return api.Run(ctx) },
		func(error) { cancel() },
	)

	// mining cron
	g.Add(
		func() error { return api.Cron(ctx, cronBeat) },
		func(error) { cancel() },
	)

	// http server
	mux := http.NewServeMux()
	mux.Handle("/", api)
	log.Printf("blockchain:web listening on %s", listen)
	ln, _ := net.Listen("tcp", listen)

	g.Add(
		func() error { return http.Serve(ln, api) },
		func(error) { ln.Close() },
	)

	// signal catcher
	cancelChan := make(chan struct{})

	g.Add(
		func() error { return interrupt(cancelChan) },
		func(error) { close(cancelChan) },
	)

	// run
	g.Run()
}
