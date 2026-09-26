package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(address string, handler http.Handler) *HTTPServer {
	return &HTTPServer{server: &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}}
}

func (server *HTTPServer) Serve(listener net.Listener) error {
	err := server.server.Serve(listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (server *HTTPServer) Shutdown(ctx context.Context) error {
	if err := server.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown onairmux server: %w", err)
	}
	return nil
}
