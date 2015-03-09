package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/fishworks/api"
	"github.com/julienschmidt/httprouter"
)

var (
	apps []*api.App
)

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

// ServeRequest processes a single HTTP request.
func (s *HTTPServer) ServeRequest(w http.ResponseWriter, req *http.Request) {
	s.srv.Handler.ServeHTTP(w, req)
}

// NewServer sets up the required Server and does protocol specific checking.
func NewServer(proto, addr string) (*HTTPServer, error) {
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
			"/apps":  getAppsJSON,
		},
		"POST": {
			"/apps": createApp,
		},
	}

	for method, routes := range m {
		for route, funct := range routes {
			r.Handle(method, route, func(h httprouter.Handle) httprouter.Handle {
				return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
					log.Debugf("%s %s", r.Method, r.RequestURI)
					// Delegate request to the given handle
					h(w, r, p)
					return
				}
			}(funct))
		}
	}

	return r
}

// WriteJSON writes the value v to the http response stream as json with standard
// json encoding.
func WriteJSON(w http.ResponseWriter, v interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func ping(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte{'P', 'O', 'N', 'G'})
}

func getAppsJSON(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if len(apps) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		if err := WriteJSON(w, apps, http.StatusOK); err != nil {
			log.Error(err)
		}
	}
}

func createApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	app := api.NewApp("")
	apps = append(apps, app)
	app.Log("created initial release")
}
