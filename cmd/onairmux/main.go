package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/wyw14/solo-035/internal/api"
	"github.com/wyw14/solo-035/internal/service"
)

func main() {
	address := flag.String("addr", "127.0.0.1:21242", "listen address")
	data := flag.String("data", filepath.Join(os.TempDir(), "onairmux-data"), "data directory")
	flag.Parse()
	runtime, err := service.NewRuntime(*data, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := service.Bootstrap(runtime); err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewHTTPServer(*address, api.NewServer(runtime, os.DirFS("web")).Router())
	errors := make(chan error, 1)
	go func() { errors <- server.Serve(listener) }()
	fmt.Printf("OnAirMux listening on http://%s\n", *address)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errors:
		if err != nil {
			log.Fatal(err)
		}
	case <-signals:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
