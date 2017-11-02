// +build !test

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	l "github.com/go-kit/kit/log"
	"github.com/spiffcp/fohnhab"
)

func main() {
	var (
		httpAddr = flag.String("listen", ":8080", "HTTP listen and serve address for service")
	)
	flag.Parse()
	ctx := context.Background()
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	logger := l.NewLogfmtLogger(os.Stderr)
	s := fohnhab.NewService(logger)
	endpoints := fohnhab.MakeEndpoints(s, logger)
	go func() {
		log.Println("http:", *httpAddr)
		handler := fohnhab.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()
	// Block from exiting unless error or interupt is recieved
	log.Fatalln(<-errChan)
}
