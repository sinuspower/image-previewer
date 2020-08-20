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
			log.Printf("[ERROR] server shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("[INFO] listening port %d; cache size: %d images", s.port, s.cacheSize)
	fmt.Fprintln(s.logOutput)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", ErrListenAndServe, err)
	}

	<-idleConnsClosed
	fmt.Fprintln(s.logOutput)
	log.Println("[INFO] server stopped")
	err := cache.Clear()
	if err != nil {
		log.Println("[WARN] can not clear cache")
	} else {
		log.Println("[INFO] cache cleared")
	}

	return nil
}

func fillHandler(w http.ResponseWriter, r *http.Request) {
	fromHost := r.RemoteAddr
	path := r.URL.Path
	rqHeader := r.Header.Clone() // copy original request headers
	var e error

	log.Printf("[INFO] get request from %s; path: %s", fromHost, path)

	cutter, err := NewCutter(path)
	if err != nil {
		e = fmt.Errorf("%s: %w", ErrCreateCutter, err)
		sendResponse(w, 400, rqHeader, fromHost, nil, e)

		return
	}

	// make response from cache if requested image is in cache
	image, ok, err := cache.GetFile(path)
	if err != nil {
		log.Println("[WARN] can not get preview from cache:", err)
	}
	if ok {
		log.Println("[INFO] get preview from cache")
		sendResponse(w, 200, rqHeader, fromHost, image, nil)

		return
	}

	image, rsHeader, err := cutter.LoadImage(rqHeader)
	if err != nil {
		e = fmt.Errorf("%s: %w", ErrCanNotLoadImage, err)
		sendResponse(w, 500, rsHeader, fromHost, nil, e)

		return
	}

	image, err = cutter.Cut(image)
	if err != nil {
		e = fmt.Errorf("%s: %w", ErrCanNotCutImage, err)
		sendResponse(w, 500, rsHeader, fromHost, nil, e)

		return
	}

	// put resized image into cache
	err = cache.PutFile(path, image)
	if err != nil {
		log.Println("[WARN] can not put preview into cache:", err)
	} else {
		log.Println("[INFO] put preview into cache")
	}

	sendResponse(w, 200, rsHeader, fromHost, image, nil)
}

func sendResponse(w http.ResponseWriter, status int, header http.Header, toHost string, data []byte, err error) {
	// copy headers
	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(status)
	if err != nil {
		log.Println("[ERROR]", err)
		data = []byte(err.Error() + "\n")
	}

	written, err := w.Write(data)
	if err != nil {
		log.Println("[ERROR]", fmt.Errorf("%s: %w", ErrWritingResponse, err))

		return
	}

	log.Printf("[INFO] send response to %s, %d bytes, status %d", toHost, written, status)
}
