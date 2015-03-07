package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

// Server defines a server which serves API requests.
type Server interface {
	Serve() error
	Close() error
}

// HTTPServer is an API Server which listens and responds to HTTP requests.
type HTTPServer struct {
	srv *http.Server
	l   net.Listener
}

// Serve starts the HTTP server, accepting all new connections.
func (s *HTTPServer) Serve() error {
	return s.srv.Serve(s.l)
}

// Close shuts down the HTTP server, dropping all current connections.
func (s *HTTPServer) Close() error {
	return s.l.Close()
}

// NewServer sets up the required Server and does protocol specific checking.
func NewServer(proto, addr string) (Server, error) {
	switch proto {
	case "tcp":
		return setupTCPHTTP(addr)
	case "unix":
		return setupUnixHTTP(addr)
	default:
		return nil, fmt.Errorf("Invalid protocol format.")
	}
}

func setupTCPHTTP(addr string) (*HTTPServer, error) {
	r := createRouter()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{&http.Server{Addr: addr, Handler: r}, l}, nil
}

func setupUnixHTTP(addr string) (*HTTPServer, error) {
	r := createRouter()

	if err := syscall.Unlink(addr); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	mask := syscall.Umask(0777)
	defer syscall.Umask(mask)

	l, err := net.Listen("unix", addr)
	if err != nil {
		return nil, err
	}

	if err := os.Chmod(addr, 0660); err != nil {
		return nil, err
	}

	return &HTTPServer{&http.Server{Addr: addr, Handler: r}, l}, nil
}

func createRouter() *httprouter.Router {
	r := httprouter.New()

	m := map[string]map[string]func(http.ResponseWriter, *http.Request, httprouter.Params){
		"GET": {
			"/_ping": ping,
		},
	}

	for method, routes := range m {
		for route, funct := range routes {
			r.Handle(method, route, funct)
		}
	}

	return r
}

func ping(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte{'P', 'O', 'N', 'G'})
}
