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
	port      int
	cacheSize int
	logOutput io.Writer
	server    *http.Server
}

var (
	ErrListenAndServe  = errors.New("error starting or closing listener")
	ErrWritingResponse = errors.New("error writing response to client")
	ErrCreateCutter    = errors.New("error during cutter creation")
	ErrCanNotLoadImage = errors.New("can not load image from server")
	ErrCanNotCutImage  = errors.New("can not cut image")
)

func NewServer(port int, cacheSize int, logOutput io.Writer) ProxyServer {
	http.HandleFunc("/fill/", fillHandler)
	log.SetOutput(logOutput)

	return &Server{
		port:      port,
		cacheSize: cacheSize,
		logOutput: logOutput,
		server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
	}
}

func (s *Server) ListenAndServe() error {
	idleConnsClosed := make(chan struct{})

	go func() { // handle signals for graceful shutdown
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-done
		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Printf("Server shutdown error: %v", err)
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

func fillHandler(w http.ResponseWriter, r *http.Request) {
	fromHost := r.RemoteAddr
	path := r.URL.Path
	var e error

	log.Printf("Get request from %s; path: %s", fromHost, path)

	cutter, err := NewCutter(path)
	if err != nil {
		e = fmt.Errorf("%s: %w", ErrCreateCutter, err)
		sendResponse(w, 400, fromHost, nil, e)

		return
	}

	image, err := cutter.LoadImage()
	if err != nil {
		e = fmt.Errorf("%s: %w", ErrCanNotLoadImage, err)
		sendResponse(w, 500, fromHost, nil, e)

		return
	}

	image, err = cutter.Cut(image)
	if err != nil {
		e = fmt.Errorf("%s: %w", ErrCanNotCutImage, err)
		sendResponse(w, 500, fromHost, nil, e)

		return
	}

	sendResponse(w, 200, fromHost, image, nil)
}

func sendResponse(w http.ResponseWriter, status int, toHost string, data []byte, err error) {
	w.WriteHeader(status)
	if err != nil {
		log.Println(err)
		data = []byte(err.Error() + "\n")
	}

	written, err := w.Write(data)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", ErrWritingResponse, err))
		return //nolint:go-lint
	}

	log.Printf("Send response to %s, %d bytes, status %d", toHost, written, status)
}
