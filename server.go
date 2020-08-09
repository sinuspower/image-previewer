package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type ProxyServer interface {
	ListenAndServe() error
}

type Server struct {
	port      uint16
	cacheSize uint16
	logOutput io.Writer
	server    *http.Server
}

var ErrListenAndServe = errors.New("error starting or closing listener")

func (s *Server) ListenAndServe() error {
	idleConnsClosed := make(chan struct{})

	go func() { // handle signals for graceful shutdown
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-done
		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Listening on port %d; cache size: %d images", s.port, s.cacheSize)
	fmt.Fprintln(s.logOutput)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", ErrListenAndServe, err)
	}

	<-idleConnsClosed
	fmt.Fprintln(s.logOutput)
	log.Println("Server stopped")

	return nil
}

func NewServer(port uint16, cacheSize uint16, logOutput io.Writer) ProxyServer {
	http.HandleFunc("/fill/", fillHandler)
	log.SetOutput(logOutput)

	return &Server{
		port:      port,
		cacheSize: cacheSize,
		logOutput: logOutput,
		server: &http.Server{
			Addr: ":" + strconv.Itoa(int(port)),
		},
	}
}

func fillHandler(w http.ResponseWriter, r *http.Request) {
	fromHost := r.RemoteAddr
	path := r.URL.Path

	log.Printf("Get request from %s; path: %s", fromHost, path)

	rsBody := []byte("Response body\n")

	written, err := w.Write(rsBody)
	if err != nil {
		log.Printf("Error sending response to %s: %s", fromHost, err)
	}
	log.Printf("Send response to %s, %d bytes", fromHost, written)
}